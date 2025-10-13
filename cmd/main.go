package main

import (
	"log"
	"os"

	"github.com/thiagola92/escudo/escudo"
)

func main() {
	println("Init")

	escudo.Init("../example")

	file := escudo.OpenFile("../example/myfile.txt", os.O_WRONLY, 0777)
	err := file.ShieldUp()

	if err != nil {
		log.Fatal(err.Error())
	}

	defer file.Close()

	log.Println("Success")
}
