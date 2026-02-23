package escudo

import (
	"os"
	"path"
	"strings"

	"github.com/thiagola92/escudo/escudo/assert"
	"github.com/thiagola92/go-lockedfile/lockedfile"
)

func newShield(shieldpath string) (shield *Shield, err error) {
	defer assert.Catch(&err)

	shield = &Shield{path: shieldpath}

	assert.FileExist(os.Mkdir(shield.path, 0770))

	file, err := os.OpenFile(shield.lockpath(), os.O_RDONLY|os.O_CREATE, 0770)

	assert.NoErr(err)
	assert.Closed(file.Close())
	assert.FileExist(os.Mkdir(shield.journalspath(), 0770))

	return shield, nil
}

func (shield *Shield) openJournal(journalpath string) (journal *Journal, err error) {
	defer assert.Catch(&err)

	journal = &Journal{path: journalpath, shield: shield, Status: INITIALIZING}
	journal.file = NewFile(journalpath, os.O_RDWR|os.O_CREATE, 0770)
	err = journal.file.Lock()

	assert.NoErr(err)

	return journal, nil
}

func (shield *Shield) anyJournal() (journal *Journal, err error) {
	defer assert.Catch(&err)

	journalspath := shield.journalspath()
	dir, err := os.Open(journalspath)

	assert.NoErr(err)

	defer dir.Close()

	files, err := dir.ReadDir(0)

	assert.NoErr(err)

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), lockext) {
			continue
		}

		journalpath := path.Join(journalspath, file.Name())
		journal, err = shield.openJournal(journalpath)

		if journal != nil {
			return journal, nil
		}
	}

	return nil, nil
}

func (shield *Shield) waitLock() (lockfile *lockedfile.File, err error) {
	defer assert.Catch(&err)

	lockpath := shield.lockpath()
	lockfile, err = lockedfile.OpenFile(lockpath, os.O_WRONLY, 0770)

	assert.NoErr(err)

	return lockfile, nil
}
