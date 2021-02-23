package source

import (
	"fmt"
	"testing"

	"github.com/alex-ant/ports/port"
	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	// Define test cases.
	cases := []struct {
		sourceFile       string
		expectedContents map[string]port.Info
		expectedError    string
	}{
		{
			sourceFile:    "invalid_path.json",
			expectedError: "source invalid_path.json doesn't exist",
		},
		{
			sourceFile: "ports_test.json",
			expectedContents: map[string]port.Info{
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
			},
			expectedError: "",
		},
	}

	// Run test cases.
	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			// Initialize new source reader.
			r, rErr := NewReader(c.sourceFile)

			if c.expectedError != "" {
				require.Error(t, rErr)
				require.Equal(t, c.expectedError, rErr.Error())
				return
			}
			require.NoError(t, rErr)

			// Define a map of expected results.
			resMap := make(map[string]port.Info)

			// Populate the map.
			readErr := r.Read(func(id string, pi port.Info) error {
				resMap[id] = pi
				return nil
			})
			require.NoError(t, readErr)

			// Compare results.
			require.Equal(t, c.expectedContents, resMap)
		})
	}
}
