package utils

import (
	"fmt"
	"os"
)

func Fatal(error error) {
	fmt.Println(error)
	os.Exit(1)
}
