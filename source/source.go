package source

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/alex-ant/ports/ports"
)

// Reader contains source file reader info.
type Reader struct {
	filePath string
}

// NewReader returns new Reader.
func NewReader(filePath string) (*Reader, error) {
	// Check whether the source file exists.
	fInfo, fInfoErr := os.Stat(filePath)
	if os.IsNotExist(fInfoErr) {
		return nil, fmt.Errorf("source %s doesn't exist", filePath)
	}

	if fInfo.IsDir() {
		return nil, fmt.Errorf("source %s is a directory", filePath)
	}

	return &Reader{
		filePath: filePath,
	}, nil
}

func (r *Reader) Read(handler func(pi *ports.PortInfo) error) error {
	// Open source file.
	f, fErr := os.Open(r.filePath)
	if fErr != nil {
		return fmt.Errorf("failed to read file %s: %v", r.filePath, fErr)
	}

	// Initialize JSON decoder.
	dec := json.NewDecoder(f)

	_, err := dec.Token()
	if err != nil {
		return fmt.Errorf("failed to open source JSON bracket 1: %v", err)
	}

	// Range over root map keys.
	for dec.More() {
		tt, err := dec.Token()
		if err != nil {
			return fmt.Errorf("failed to open source JSON bracket 2: %v", err)
		}

		// Extract port info value map.
		for dec.More() {
			var m interface{}
			var pi ports.PortInfo
			if err := dec.Decode(&m); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("failed to decode JSON: %v", err)
			}

			mB, err := json.Marshal(m)
			if err != nil {
				return fmt.Errorf("failed to marshal port item: %v", err)
			}

			err = json.Unmarshal(mB, &pi)
			if err != nil {
				return fmt.Errorf("failed to unmarshal port item map: %v", err)
			}

			// Save port ID.
			pi.Id = tt.(string)

			// Pass the retrieved port info into the handler.
			err = handler(&pi)
			if err != nil {
				return fmt.Errorf("handler error: %v", err)
			}

			// Assign next port ID.
			tt, err = dec.Token()
			if err != nil {
				return fmt.Errorf("failed to close source port map bracket: %v", err)
			}
		}

		_, err = dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("failed to close source bracket 2: %v", err)
		}
	}

	_, err = dec.Token()
	if err != io.EOF {
		return fmt.Errorf("failed to close source bracket 1: %v", err)
	}

	return nil
}
