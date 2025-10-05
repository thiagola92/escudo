package main

import (
	"os"

	"github.com/thiagola92/escudo/escudo"
)

func main() {
	println("Init")

	escudo.Init("../example")

	file := escudo.OpenFile("../example/myfile.txt", os.O_WRONLY, 0777)
	err := file.ShieldUp()

	if err != nil {
		println(err.Error())
		return
	}

	file.Close()
}
