package escudo

import (
	"path"
)

var lockext = ".escl"
var tempext = ".esct"
var jourext = ".escj"

func Init(dirpath string) (shield *Shield, err error) {
	shield, err = newShield(path.Join(dirpath, ".escudo"))

	if err != nil {
		return nil, err
	}

	journal, err := shield.anyJournal()

	if err != nil {
		return nil, err
	}

	if journal != nil {
		err = journal.resume()

		if err != nil {
			return nil, err
		}
	}

	return shield, nil
}
