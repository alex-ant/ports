package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/gobuffalo/packr"
	"github.com/jackc/pgx"
	migrate "github.com/rubenv/sql-migrate"
)

const (
	migrationsDir   string = "./migrations"
	migrationsTable string = "ports-migrations"
)

// Client contains connection data.
type Client struct {
	user   string
	pass   string
	host   string
	port   int
	dbName string
	pool   *pgx.ConnPool
}

// New returns new Client.
func New(user, pass, host string, port int, dbName string, timeout int) (c *Client, err error) {
	c = &Client{
		user:   user,
		pass:   pass,
		host:   host,
		port:   port,
		dbName: dbName,
	}

	// Verify the connection.
	connected := make(chan bool)

	go func() {
		for {
			// Connect to the database.
			c.pool, err = pgx.NewConnPool(pgx.ConnPoolConfig{
				ConnConfig: pgx.ConnConfig{
					User:     user,
					Password: pass,
					Host:     host,
					Port:     uint16(port),
					Database: dbName,
				},
				AcquireTimeout: time.Second * 5,
			})

			if err == nil {
				// Ping the database.
				err = c.ping()
				if err == nil {
					log.Println("Connected to PostgreSQL")
					connected <- true
					break
				}
			}

			// Sleep one second if connection failed.
			log.Printf("Failed to connect to PostgreSQL [host: %s:%d, err: %v], reconnecting in 5 seconds", host, port, err)
			time.Sleep(5 * time.Second)
		}
	}()

	select {
	case <-connected:
	case <-time.After(time.Duration(timeout) * time.Second):
		err = errors.New("timeout occured connecting to PostgreSQL")
		return
	}

	// Run migrations.
	err = runDatabaseMigrations(user, pass, host, port, dbName)
	if err != nil {
		return nil, fmt.Errorf("failed to run DB migrations: %v", err)
	}

	return
}

func runDatabaseMigrations(user, pass, host string, port int, dbName string) error {
	var db *sql.DB
	var dbErr error

	db, dbErr = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, user, pass, dbName))
	if dbErr == nil {
		dbErr = db.Ping()
	}
	if dbErr != nil {
		if !strings.Contains(strings.ToLower(dbErr.Error()), "ssl is not enabled") {
			return fmt.Errorf("failed to connect database for migrations with SSL required: %v", dbErr)
		}

		// Attempt to connect without SSL.
		fmt.Println("SSL is not enabled in DB, connecting without it")

		db, dbErr = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbName))
		if dbErr != nil {
			return fmt.Errorf("failed to connect database for migrations with SSL disabled: %v", dbErr)
		}
	}
	defer db.Close()

	migrations := &migrate.PackrMigrationSource{
		Box: packr.NewBox(migrationsDir),
	}

	log.Println("running migrations")

	migrate.SetTable(migrationsTable)

	n, nErr := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if nErr != nil {
		return fmt.Errorf("failed to run migrations: %v", nErr)
	}

	log.Printf("Applied %d migrations!\n", n)

	return nil
}

func (c *Client) ping() error {
	pinged := make(chan error)

	go func() {
		_, err := c.pool.Exec(";")
		if err != nil {
			pinged <- err
			return
		}

		pinged <- nil
	}()

	select {
	case err := <-pinged:
		if err != nil {
			return err
		}
	case <-time.After(time.Second):
		return errors.New("timeout occured connecting to PostgreSQL")
	}

	return nil
}

// Close closes active connection.
func (c *Client) Close() {
	c.pool.Close()
}
