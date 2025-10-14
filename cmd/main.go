package main

import (
	"log"
	"os"

	"github.com/thiagola92/escudo/escudo"
)

func main() {
	println("Init")

	shield, err := escudo.Init("../example")

	if err != nil {
		log.Fatal(err.Error())
	}

	file := escudo.OpenFile("../example/myfile.txt", os.O_RDONLY, 0777)
	file2 := escudo.OpenFile("../example/other.txt", os.O_RDWR, 0777)
	journal, err := shield.MyJournal()

	if err != nil {
		log.Fatal(err.Error())
	}

	defer file2.Close()
	defer file.Close()

	log.Println("Success")
}
