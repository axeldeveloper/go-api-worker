package common

import (
	"os"
)

func Init() {
	os.Setenv("TZ", "America/Sao_Paulo")
}
