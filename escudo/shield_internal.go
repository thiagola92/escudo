package escudo

import (
	"os"
	"path"
	"strings"

	"github.com/thiagola92/escudo/escudo/assert"
	"github.com/thiagola92/go-lockedfile/lockedfile"
)

func (shield *Shield) waitLock() *lockedfile.File {
	lockpath := shield.lockpath()
	lockfile, err := lockedfile.OpenFile(lockpath, os.O_WRONLY, 0770)

	defer assert.Catch()
	assert.NoErr(err)

	return lockfile
}

func (shield *Shield) openJournal(journalpath string) *Journal {
	journal := &Journal{path: journalpath, shield: shield, Status: INITIALIZING}
	journal.file = OpenFile(journalpath, os.O_RDWR|os.O_CREATE, 0770)
	journal.file.Lock()

	defer assert.Catch()
	assert.NoErr(assert.Err)

	return journal
}

func (shield *Shield) anyJournal() *Journal {
	var journal *Journal

	journalspath := shield.journalspath()
	dir, err := os.Open(journalspath)

	defer assert.Catch()
	assert.NoErr(err)

	defer dir.Close()

	files, err := dir.ReadDir(0)

	assert.NoErr(err)

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
