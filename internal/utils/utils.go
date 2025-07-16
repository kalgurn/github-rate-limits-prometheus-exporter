package utils

import (
	"errors"
	"fmt"
	"os"
	"log"
)

func RespError(err error) error {
	log.SetOutput(os.Stderr)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error during the call execution: %s\n", err)
		log.Printf("ERRO %s", errMsg)
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
