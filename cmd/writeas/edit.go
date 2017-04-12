package main

import (
	"os"
)

var editors = []string{"WRITEAS_EDITOR", "EDITOR"}

func getConfiguredEditor() string {
	for _, v := range editors {
		if e := os.Getenv(v); e != "" {
			return e
		}
	}
	return ""
}
