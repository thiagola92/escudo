package escudo

import (
	"github.com/thiagola92/escudo/escudo/assert"
)

type Journal struct {
	path   string
	file   *File
	shield *Shield

	Entries []*JournalEntry
	Status  int
}

// Obtain permission to interact with multiple files.
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

	journal.lockall()

	if assert.Err != nil {
		return assert.Err
	}

	journal.Status = WRITING
	journal.replace()

	return assert.Err
}

// Commit changes to files and journal.
func (journal *Journal) Commit() (err error) {
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

	journal.commitall()

	if assert.Err != nil {
		return assert.Err
	}

	journal.Status = WRITING
	journal.replace()

	return assert.Err
}

// Delete temporary files and journal.
func (journal *Journal) Close() (err error) {
	lock := journal.shield.waitLock()

	if assert.Err != nil {
		return assert.Err
	}

	defer lock.Close()

	journal.Status = DELETING
	journal.replace()

	if assert.Err != nil {
		return assert.Err
	}

	journal.closeall()

	if assert.Err != nil {
		return assert.Err
	}

	journal.close()

	return assert.Err
}
