package escudo

import "github.com/thiagola92/escudo/escudo/assert"

type Journal struct {
	path   string
	file   *File
	shield *Shield

	Entries []*JournalEntry
	Status  int
}

func (journal *Journal) Lock(files ...*File) (err error) {
	journal.Entries = newJournalEntries(files)

	if assert.Err != nil {
		return assert.Err
	}

	lock := journal.shield.waitLock()

	if assert.Err != nil {
		return assert.Err
	}

	defer lock.Close()

	// Attempt to get all locks.
	for index, entry := range journal.Entries {
		entry.file.Lock()

		// Undo all locks in case one fail.
		if assert.Err != nil {
			for index >= 0 {
				entry.file.Close()
				index -= 1
			}

			return assert.Err
		}
	}

	journal.Status = WRITING
	journal.replace()

	return assert.Err
}

func (journal *Journal) Commit() error {
	lock := journal.shield.waitLock()

	if assert.Err != nil {
		return assert.Err
	}

	defer lock.Close()

	journal.Status = REPLACING
	journal.replace()

	if assert.Err != nil {
		return assert.Err
	}

	// TODO: Make the actual replacement of files.
	// Right now I'm just updating the journal.

	journal.Status = WRITING
	journal.replace()

	return assert.Err
}

func (journal *Journal) Push() error {
	lock := journal.shield.waitLock()

	if assert.Err != nil {
		return assert.Err
	}

	defer lock.Close()

	journal.Status = REPLACING
	journal.replace()

	if assert.Err != nil {
		return assert.Err
	}

	// TODO: Make the actual replacement of files.
	// Right now I'm just updating the journal.

	journal.Status = DELETING
	journal.replace()

	if assert.Err != nil {
		return assert.Err
	}

	journal.close()

	return assert.Err
}
