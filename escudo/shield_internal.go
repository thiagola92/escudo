package escudo

import (
	"os"
	"path"
	"strings"

	"github.com/thiagola92/escudo/escudo/assert"
	"github.com/thiagola92/go-lockedfile/lockedfile"
)

func (shield *Shield) anyJournal() *Journal {
	var journal *Journal

	journalspath := shield.journalspath()
	dir, err := os.Open(journalspath)

	assert.NoError(err)

	defer dir.Close()

	files, err := dir.ReadDir(0)

	assert.NoError(err)

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), lockext) {
			continue
		}

		journalpath := path.Join(journalspath, file.Name())
		journal = shield.openJournal(journalpath)

		if journal != nil {
			return journal
		}
	}

	return nil
}

func (shield *Shield) getLock() *lockedfile.File {
	// Block until get the lock.
	lockpath := shield.lockpath()
	lockfile, err := lockedfile.OpenFile(lockpath, os.O_WRONLY, 0770)

	assert.NoError(err)

	return lockfile
}

func (shield *Shield) openJournal(journalpath string) *Journal {
	journal := &Journal{path: journalpath, shield: shield, Status: INITIALIZING}
	journal.file = OpenFile(journalpath, os.O_RDWR|os.O_CREATE, 0770)
	journal.file.Lock()

	if assert.Err != nil {
		return nil
	}

	return journal
}
