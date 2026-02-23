package escudo

import (
	"path/filepath"

	"github.com/thiagola92/escudo/escudo/assert"
)

func newJournalEntry(file *File) (entry *JournalEntry, err error) {
	defer assert.Catch(&err)

	// NOTE: We really want absolute path?
	// This would means that the user can't move the directory...
	entry = &JournalEntry{file: file}
	entry.Path, err = filepath.Abs(file.path)

	assert.NoErr(err)

	return entry, nil
}

func newJournalEntries(files []*File) (entries []*JournalEntry, err error) {
	defer assert.Catch(&err)

	entries = make([]*JournalEntry, len(files))

	for index, file := range files {
		entry, err := newJournalEntry(file)

		assert.NoErr(err)

		entries[index] = entry
	}

	return entries, nil
}
