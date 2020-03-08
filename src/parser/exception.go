package parser

import (
	"fmt"
	"os"
)

func errorHandler(err error) {
	if err != nil {
		fmt.Println("Encountered error:")
		fmt.Println(err)
		os.Exit(0)
	}
}
