package escudo

import (
	"os"
	"path"
)

type Journal struct {
	path string

	file    *File
	entries []JournalEntry // TODO: Decide the structure.
}

func (shield *Shield) MyJournal() (*Journal, error) {
	return openJournal(shield.journalpath(os.Getpid()))
}

func (shield *Shield) AnyJournal() (*Journal, error) {
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
		journalpath := path.Join(journalspath, file.Name())
		journal, err = openJournal(journalpath)

		if err == nil {
			return journal, nil
		}
	}

	return nil, nil
}

func (journal *Journal) LockFiles(files []*File) error {
	for index, file := range files {
		err := file.Lock()

		if err != nil {
			for index >= 0 {
				file.Close()
				index -= 1
			}

			return err
		}
	}

	return nil
}
