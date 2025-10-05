package escudo

import (
	"os"
	"path"
)

type Escudo struct {
	escudopath   string
	lockpath     string
	journalspath string
}

func Init(dirpath string) (*Escudo, error) {
	var err error
	escudo := &Escudo{}

	escudo.escudopath = path.Join(dirpath, ".escudo")
	escudo.lockpath = path.Join(escudo.escudopath, "global.escl")
	escudo.journalspath = path.Join(escudo.escudopath, "journals")

	err = os.Mkdir(escudo.escudopath, 0770)

	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(escudo.lockpath, os.O_RDONLY|os.O_CREATE, 0770)

	if err != nil {
		return nil, err
	}

	err = file.Close()

	if err != nil {
		return nil, err
	}

	err = os.Mkdir(escudo.journalspath, 0770)

	if err != nil {
		return nil, err
	}

	return escudo, nil
}
