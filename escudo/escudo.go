package escudo

import (
	"path"

	"github.com/thiagola92/escudo/escudo/assert"
)

var lockext = ".escl"
var tempext = ".esct"
var jourext = ".escj"

func Init(dirpath string) (*Shield, error) {
	shield := newShield(path.Join(dirpath, ".escudo"))

	if assert.Err != nil {
		return nil, assert.Err
	}

	journal := shield.anyJournal()

	if assert.Err != nil {
		return nil, assert.Err
	}

	// TODO: complete journal
	println(journal)

	return shield, nil
}
