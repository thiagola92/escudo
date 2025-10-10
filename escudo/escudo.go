package escudo

import (
	"os"
	"path"
)

var lockext = ".escl"
var tempext = ".esct"
var jourext = ".escj"

type Escudo struct {
	path string
}

func Init(dirpath string) (*Escudo, error) {
	var err error

	escudo := &Escudo{
		path: path.Join(dirpath, ".escudo"),
	}

	err = os.Mkdir(escudo.path, 0770)

	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(escudo.lockpath(), os.O_RDONLY|os.O_CREATE, 0770)

	if err != nil {
		return nil, err
	}

	err = file.Close()

	if err != nil {
		return nil, err
	}

	err = os.Mkdir(escudo.journalspath(), 0770)

	if err != nil {
		return nil, err
	}

	return escudo, nil
}
