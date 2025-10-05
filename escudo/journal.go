package escudo

import (
	"fmt"
	"os"
	"path"
)

const (
	// When setted means that we need create the lock file
	// (and maybe the temporary file too).
	//
	// Finding this state in the journal means that we can't recover
	// from the crash and we should just delete the lock file.
	LOCKING = 100

	// When setted means that the user can edit the temporary file
	// as they wish.
	//
	// Finding this state in the journal means that we can't recover
	// from the crash and we should just delete the temporary & lock files.
	WRITING = 300

	// When setted means that we need to replace the original file
	// with the temporary file.
	//
	// Finding this state in the journal means that we CAN recover
	// from the crash and we should continue from where we stopped.
	REPLACING = 400

	// When setted means that we need to delete the temporary file.
	//
	// Finding this state in the journal means that we should
	// continue from where we stopped.
	DELETING = 500
)

type Journal struct {
	path string

	file    *os.File       // TODO: replace by ours File.
	entries []JournalEntry // TODO: Decide the structure
}

type JournalEntry struct {
	file   File
	Path   string
	Status int
}

func OpenJournal(journalspath string) (*Journal, error) {
	var err error
	var journal *Journal

	pid := os.Getpid()
	journalname := fmt.Sprintf("%d.escj", pid)
	journal.path = path.Join(journalspath, journalname)

	journal.file, err = os.OpenFile(journal.path, os.O_RDWR|os.O_CREATE, 0770)

	if err != nil {
		return nil, err
	}

	return journal, nil
}

func (journal *Journal) OpenFiles(files []*File) error {
	return nil
}
