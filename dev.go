package writeascli

import (
	"os"
)

func IsDev() bool {
	return os.Getenv("WRITEAS_DEV") == "1"
}
