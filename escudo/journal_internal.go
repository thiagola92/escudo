package escudo

import (
	"os"
)

func openJournal(journalpath string) (*Journal, error) {
	var err error

	journal := &Journal{path: journalpath}
	journal.file = OpenFile(journalpath, os.O_RDWR|os.O_CREATE, 0770)
	err = journal.file.Lock()

	if err != nil {
		return nil, err
	}

	return journal, nil
}
