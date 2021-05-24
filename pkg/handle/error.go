package handle

import (
	"fmt"
	"log"
)

func HandleError(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}