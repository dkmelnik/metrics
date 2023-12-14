package configs

import (
	"errors"
	"flag"
)

func CheckUnknownFlags() error {
	if len(flag.Args()) > 0 {
		return errors.New("unknown flags or parameters")
	}
	return nil
}
