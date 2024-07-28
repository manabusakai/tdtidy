package main

import (
	"context"
	"log"

	"github.com/manabusakai/tdtidy"
)

func main() {
	app, err := tdtidy.New(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}
