package db

import (
	"testing"

	"github.com/alex-ant/ports/port"
	"github.com/stretchr/testify/require"
)

func TestInsertPortInfo(t *testing.T) {
	// Define test port data.
	testPorts := map[string]port.Info{
		"AEAJM": {
			Name:        "Ajman",
			City:        "Ajman",
			Country:     "United Arab Emirates",
			Alias:       []string{},
			Regions:     []string{},
			Coordinates: [2]float32{55.5136433, 25.4052165},
			Province:    "Ajman",
			Timezone:    "Asia/Dubai",
			Unlocs:      []string{"AEAJM"},
			Code:        "52000",
		},
		"AEAUH": {
			Name:        "Abu Dhabi",
			City:        "Abu Dhabi",
			Country:     "United Arab Emirates",
			Alias:       []string{},
			Regions:     []string{},
			Coordinates: [2]float32{54.37, 24.47},
			Province:    "Abu Z¸aby [Abu Dhabi]",
			Timezone:    "Asia/Dubai",
			Unlocs:      []string{"AEAUH"},
			Code:        "52001",
		},
	}

	// Populate database.
	for id, pi := range testPorts {
		insertErr := testClient.InsertPortInfo(id, &pi)
		require.NoError(t, insertErr)
	}

	// Validate DB records.
	validateSTMT := `
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
		ORDER BY code
	`

	fetchPorts := func() map[string]port.Info {
		resMap := make(map[string]port.Info)

		rows, rowsErr := testClient.pool.Query(validateSTMT)
		if rowsErr != nil {
			require.NoError(t, rowsErr)
		}

		defer rows.Close()
		for rows.Next() {
			var id string
			var pi port.Info

			if scanErr := rows.Scan(
				&id,
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
				require.NoError(t, scanErr)
			}

			resMap[id] = pi
		}
		if err := rows.Err(); err != nil {
			require.NoError(t, err)
		}

		return resMap
	}

	// Compare results.
	require.Equal(t, testPorts, fetchPorts())

	// Update records.
	newTestPI := port.Info{
		Name:        "NewName",
		City:        "NewCity",
		Country:     "NewCountry",
		Alias:       []string{"alias1"},
		Regions:     []string{"region1"},
		Coordinates: [2]float32{1, 2},
		Province:    "NewProvince",
		Timezone:    "NewTimezone",
		Unlocs:      []string{"AEAJM"},
		Code:        "52000",
	}

	insertErr := testClient.InsertPortInfo("AEAJM", &newTestPI)
	require.NoError(t, insertErr)

	testPorts["AEAJM"] = newTestPI

	// Compare updated results.
	require.Equal(t, testPorts, fetchPorts())
}
