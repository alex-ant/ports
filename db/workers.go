package db

import (
	"fmt"

	"github.com/alex-ant/ports/ports"
)

// FetchPortInfo returns all port info items from the DB.
func (c *Client) FetchPortInfo() (map[string]*ports.PortInfo, error) {
	res := make(map[string]*ports.PortInfo)

	stmt := `
	SELECT id,
			name,
			city,
			country,
			alias,
			regions,
			lat,
			lng,
			province,
			timezone,
			unlocs,
			code
		FROM ports
	`

	rows, rowsErr := c.pool.Query(stmt)
	if rowsErr != nil {
		return nil, fmt.Errorf("failed to execute fetch port info query: %v", rowsErr)
	}

	defer rows.Close()
	for rows.Next() {
		var pi ports.PortInfo
		pi.Coordinates = make([]float32, 2)

		if scanErr := rows.Scan(
			&pi.Id,
			&pi.Name,
			&pi.City,
			&pi.Country,
			&pi.Alias,
			&pi.Regions,
			&pi.Coordinates[0],
			&pi.Coordinates[1],
			&pi.Province,
			&pi.Timezone,
			&pi.Unlocs,
			&pi.Code,
		); scanErr != nil {
			return nil, fmt.Errorf("failed to execute scan row: %v", scanErr)
		}

		res[pi.Id] = &pi
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching rows: %v", err)
	}

	return res, nil
}

// InsertPortInfo stores provided port info into the DB replacing the respective
// row on duplicate port id.
func (c *Client) InsertPortInfo(pi *ports.PortInfo) error {
	stmt := `
	INSERT INTO ports
				(id,
				name,
				city,
				country,
				alias,
				regions,
				lat,
				lng,
				province,
				timezone,
				unlocs,
				code)
	VALUES		($1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8,
				$9,
				$10,
				$11,
				$12)
	ON CONFLICT (id) DO UPDATE 
	SET name = $2, 
		city = $3,
		country = $4,
		alias = $5,
		regions = $6,
		lat = $7,
		lng = $8,
		province = $9,
		timezone = $10,
		unlocs = $11,
		code = $12
	`

	var lat, lng float32
	if len(pi.Coordinates) == 2 {
		lat, lng = pi.Coordinates[0], pi.Coordinates[1]
	}

	_, err := c.pool.Exec(
		stmt,
		pi.Id,
		pi.Name,
		pi.City,
		pi.Country,
		pi.Alias,
		pi.Regions,
		lat,
		lng,
		pi.Province,
		pi.Timezone,
		pi.Unlocs,
		pi.Code,
	)

	return err
}
