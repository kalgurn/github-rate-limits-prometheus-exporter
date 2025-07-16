package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

var verbose bool

func init() {
	verboseEnv := GetOSVar("VERBOSE")
	verbose = strings.ToLower(verboseEnv) == "true" || verboseEnv == "1"
}

func RespError(err error) error {
	if err != nil {
		errMsg := fmt.Sprintf("there was an error during the call execution: %s", err)

		if verbose {
			log.SetOutput(os.Stderr)
			log.Printf("ERROR: %s", errMsg)
		}

		return errors.New(errMsg)
	}
	return nil
}

func GetOSVar(envVar string) string {
	value, present := os.LookupEnv(envVar)
	if !present {
		err := fmt.Sprintf("environment variable %s not set", envVar)
		RespError(errors.New(err))
		return ""
	}
	return value
}
