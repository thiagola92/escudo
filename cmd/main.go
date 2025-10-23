package main

import "cmp"

// func main() {
// 	println("Init")

// 	shield, err := escudo.Init("../example")

// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	file := escudo.OpenFile("../example/myfile.txt", os.O_RDONLY, 0777)
// 	file2 := escudo.OpenFile("../example/other.txt", os.O_RDWR, 0777)
// 	journal, err := shield.GetJournal()

// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	err = journal.LockFiles(file, file2)

// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	log.Println("Success")
// }

func main() {
	println(cmp.Or(1, g(), h()))
}

func g() int {
	println("g()")
	return 3
}

func h() int {
	println("h()")
	return 2
}
