package db

import (
	"errors"

	"github.com/alex-ant/ports/port"
)

// InsertPortInfo stores provided port info into the DB replacing the respective
// row on duplicate port id.
func (c *Client) InsertPortInfo(id string, pi *port.Info) error {
	if id == "" {
		return errors.New("empty port id provided")
	}

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

	_, err := c.pool.Exec(
		stmt,
		id,
		pi.Name,
		pi.City,
		pi.Country,
		pi.Alias,
		pi.Regions,
		pi.Coordinates[0],
		pi.Coordinates[1],
		pi.Province,
		pi.Timezone,
		pi.Unlocs,
		pi.Code,
	)

	return err
}
