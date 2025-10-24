package escudo

import (
	"path/filepath"

	"github.com/thiagola92/escudo/escudo/assert"
)

func newJournalEntry(file *File) (entry *JournalEntry) {
	var err error

	// NOTE: We really want absolute path?
	// This would means that the user can't move the directory...
	entry = &JournalEntry{file: file}
	entry.Path, err = filepath.Abs(file.path)

	defer assert.Catch()
	assert.NoErr(err)

	return entry
}

func newJournalEntries(files []*File) (entries []*JournalEntry) {
	entries = make([]*JournalEntry, len(files))

	defer assert.Catch()

	for index, file := range files {
		entry := newJournalEntry(file)

		assert.NoErr(assert.Err)

		entries[index] = entry
	}

	return entries
}
