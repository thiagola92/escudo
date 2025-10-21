package escudo

import "path/filepath"

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
