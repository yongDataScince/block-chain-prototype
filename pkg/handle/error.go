package handle

import "log"

func HandleError(err error) {
	log.Fatal(err)
}