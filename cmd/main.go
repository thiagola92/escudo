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

	file := escudo.NewFile("../example/myfile.txt", os.O_RDONLY, 0777)
	file2 := escudo.NewFile("../example/other.txt", os.O_RDWR, 0777)
	journal, err := shield.GetJournal()

	if err != nil {
		log.Fatal(err.Error())
	}

	err = journal.Lock(file, file2)

	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Success")
}
