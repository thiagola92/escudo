package escudo

import (
	"os"
	"path"
)

var lockext = ".escl"
var tempext = ".esct"
var jourext = ".escj"

type Shield struct {
	path string
}

func Init(dirpath string) (*Shield, error) {
	var err error

	shield := &Shield{
		path: path.Join(dirpath, ".escudo"),
	}

	err = os.Mkdir(shield.path, 0770)

	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(shield.lockpath(), os.O_RDONLY|os.O_CREATE, 0770)

	if err != nil {
		return nil, err
	}

	err = file.Close()

	if err != nil {
		return nil, err
	}

	err = os.Mkdir(shield.journalspath(), 0770)

	if err != nil {
		return nil, err
	}

	return shield, nil
}
