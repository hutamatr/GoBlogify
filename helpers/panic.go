package helpers

import (
	"log"
)

func PanicError(err error) {
	if err != nil {
		log.Printf("An error occurred: %v", err)
		// err := fmt.Errorf("an error occurred: %v", err)
		panic(err)
	}
}
