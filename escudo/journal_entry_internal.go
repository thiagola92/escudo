package escudo

import (
	"path/filepath"
)

func newJournalEntry(file *File) (entry *JournalEntry, err error) {
	// NOTE: We really want absolute path?
	// This would means that the user can't move the directory...
	entry.file = file
	abspath, err := filepath.Abs(file.path)

	if err != nil {
		return nil, err
	}

	entry.Path = abspath

	return entry, nil
}

func newJournalEntries(files []*File) (entries []*JournalEntry, err error) {
	entries = make([]*JournalEntry, len(files))

	for index, file := range files {
		entry, err := newJournalEntry(file)

		if err != nil {
			return nil, err
		}

		entries[index] = entry
	}

	return entries, nil
}
