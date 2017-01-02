package main

import (
	"log"

	"github.com/sneakybueno/fli/fuego"
)

func main() {
	log.Println("Welcome to fli!")

	firebaseURL := "https://go-fli.firebaseio.com/"
	fStore := fuego.NewFStore(firebaseURL)

	json, err := fStore.ShallowGet("")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(json)
}
