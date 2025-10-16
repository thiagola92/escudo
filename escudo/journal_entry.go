package escudo

import "path/filepath"

const (
	// When setted means that we need to create the lock file
	// (and maybe the temporary file too).
	//
	// Finding this state in the journal means that we can't recover
	// from the crash and we should just delete the lock file.
	//
	// NOTE: This state is never stored in the journal file.
	INITIALIZING = 0

	// When setted means that the user can edit the temporary file
	// as they wish.
	//
	// Finding this state in the journal means that we can't recover
	// from the crash and we should just delete the temporary & lock files.
	WRITING = 100

	// When setted means that we need to replace the original file
	// with the temporary file.
	//
	// Finding this state in the journal means that we CAN recover
	// from the crash and we should continue from where we stopped.
	REPLACING = 200

	// When setted means that we need to delete the temporary file.
	//
	// Finding this state in the journal means that we should
	// continue from where we stopped.
	DELETING = 300
)

type JournalEntry struct {
	file *File

	Path   string
	Status int
}

func newJournalEntry(file *File) (*JournalEntry, error) {
	// NOTE: We really want absolute path?
	// This would means that the user can't move the directory...
	entry := &JournalEntry{file: file}
	abspath, err := filepath.Abs(file.path)

	if err != nil {
		return nil, err
	}

	entry.Path = abspath
	entry.Status = INITIALIZING

	return entry, nil
}

func toJournalEntries(files []*File) ([]*JournalEntry, error) {
	entries := make([]*JournalEntry, len(files))

	for index, file := range files {
		entry, err := newJournalEntry(file)

		if err != nil {
			return nil, err
		}

		entries[index] = entry
	}

	return entries, nil
}
