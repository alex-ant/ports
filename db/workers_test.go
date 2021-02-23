package db

import (
	"testing"

	"github.com/alex-ant/ports/ports"
	"github.com/stretchr/testify/require"
)

func TestInsertAndFetchPortInfo(t *testing.T) {
	// Define test port data.
	testPorts := map[string]*ports.PortInfo{
		"AEAJM": {
			Id:          "AEAJM",
			Name:        "Ajman",
			City:        "Ajman",
			Country:     "United Arab Emirates",
			Alias:       []string{},
			Regions:     []string{},
			Coordinates: []float32{55.5136433, 25.4052165},
			Province:    "Ajman",
			Timezone:    "Asia/Dubai",
			Unlocs:      []string{"AEAJM"},
			Code:        "52000",
		},
		"AEAUH": {
			Id:          "AEAUH",
			Name:        "Abu Dhabi",
			City:        "Abu Dhabi",
			Country:     "United Arab Emirates",
			Alias:       []string{},
			Regions:     []string{},
			Coordinates: []float32{54.37, 24.47},
			Province:    "Abu Z¸aby [Abu Dhabi]",
			Timezone:    "Asia/Dubai",
			Unlocs:      []string{"AEAUH"},
			Code:        "52001",
		},
	}

	// Populate database.
	for _, pi := range testPorts {
		insertErr := testClient.InsertPortInfo(pi)
		require.NoError(t, insertErr)
	}

	// Compare results.
	res, resErr := testClient.FetchPortInfo()
	require.NoError(t, resErr)

	require.Equal(t, testPorts, res)

	// Update records.
	newTestPI := ports.PortInfo{
		Id:          "AEAJM",
		Name:        "NewName",
		City:        "NewCity",
		Country:     "NewCountry",
		Alias:       []string{"alias1"},
		Regions:     []string{"region1"},
		Coordinates: []float32{1, 2},
		Province:    "NewProvince",
		Timezone:    "NewTimezone",
		Unlocs:      []string{"AEAJM"},
		Code:        "52000",
	}

	insertErr := testClient.InsertPortInfo(&newTestPI)
	require.NoError(t, insertErr)

	testPorts["AEAJM"] = &newTestPI

	// Compare updated results.
	res, resErr = testClient.FetchPortInfo()
	require.NoError(t, resErr)

	require.Equal(t, testPorts, res)
}
