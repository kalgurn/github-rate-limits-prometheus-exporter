package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
)

// RespError logs and returns a wrapped error for callers to handle.
// This function does not terminate the process.
func RespError(err error) error {
	if err != nil {
		errMsg := fmt.Sprintf("there was an error during the call execution: %s", err)
		log.Printf("%s", errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func GetOSVar(envVar string) string {
	value, present := os.LookupEnv(envVar)
	if !present {
		err := fmt.Sprintf("environment variable %s not set", envVar)
		// Log the error but do not terminate; caller should validate the value.
		RespError(errors.New(err))
		return ""
	}
	return value
}
