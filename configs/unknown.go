package configs

import (
	"errors"
	"flag"
)

// CheckUnknownFlags checks whether there are unused passed flags.
func CheckUnknownFlags() error {
	if len(flag.Args()) > 0 {
		return errors.New("unknown flags or parameters")
	}
	return nil
}
