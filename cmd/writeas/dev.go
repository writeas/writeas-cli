package main

import (
	"os"
)

func isDev() bool {
	return os.Getenv("WRITEAS_DEV") == "1"
}
