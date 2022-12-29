package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sheepla/strans/api"
)

func main() {
	param, err := api.NewParam("en", "ja", os.Args[1], api.EngineGoogle)
	if err != nil {
		log.Println(err)
	}

	result, err := api.Translate(param, "")
	if err != nil {
		log.Println(err)
	}

	fmt.Println(result.Text)
}
