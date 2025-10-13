package escudo

import (
	"os"
	"path"

	"github.com/thiagola92/go-lockedfile/lockedfile"
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

	file    *lockedfile.File // TODO: replace by ours File.
	entries []JournalEntry   // TODO: Decide the structure.
}

type JournalEntry struct {
	file File

	Path   string
	Status int
}

func OpenJournal(journalpath string) (*Journal, error) {
	var err error

	journal := &Journal{path: journalpath}
	journal.file, err = lockedfile.TryOpenFile(journalpath, os.O_RDWR|os.O_CREATE, 0770)

	if err != nil {
		return nil, err
	}

	return journal, nil
}

func AnyJournal(journalspath string) (*Journal, error) {
	var journal *Journal
	var err error

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
		journal, err = OpenJournal(journalpath)

		if err == nil {
			return journal, nil
		}
	}

	return nil, nil
}

// TODO: Give a bunch of files and let the journal attempt to open all of them.
func (journal *Journal) OpenFiles(files []*File) error {
	return nil
}
