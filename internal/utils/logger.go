package utils

import (
	"github.com/s3nt3/myproxy/internal/logger"
)

func FatalCheck(err error, message string) {
	if err != nil {
		logger.Fatal(err)
	} else {
		logger.Logger.Debug(message)
	}
}

func ErrorCheck(err error, message string) bool {
	if err != nil {
		logger.Logger.Error(err)
	} else {
		logger.Logger.Debug(message)
	}

	return err == nil
}
