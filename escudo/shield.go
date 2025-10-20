package escudo

import (
	"os"
	"path"
	"strings"

	"github.com/thiagola92/go-lockedfile/lockedfile"
)

type Shield struct {
	path string
}

func (shield *Shield) GetJournal() (*Journal, error) {
	journalpath := shield.journalpath(os.Getpid())
	return shield.openJournal(journalpath)
}

func (shield *Shield) anyJournal() (*Journal, error) {
	var journal *Journal
	var err error

	journalspath := shield.journalspath()
	dir, err := os.Open(journalspath)

	if err != nil {
		return nil, err
	}

	defer dir.Close()

	files, err := dir.ReadDir(0)

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), lockext) {
			continue
		}

		journalpath := path.Join(journalspath, file.Name())
		journal, err = shield.openJournal(journalpath)

		if err == nil {
			return journal, nil
		}
	}

	return nil, nil
}

func (shield *Shield) openJournal(journalpath string) (*Journal, error) {
	var err error

	journal := &Journal{path: journalpath, shield: shield}
	journal.file = OpenFile(journalpath, os.O_RDWR|os.O_CREATE, 0770)
	err = journal.file.Lock()

	if err != nil {
		return nil, err
	}

	return journal, nil
}

func (shield *Shield) getlock() (*lockedfile.File, error) {
	// Note that this will lock until get the lock.
	lockpath := shield.lockpath()
	return lockedfile.OpenFile(lockpath, os.O_WRONLY, 0770)
}
