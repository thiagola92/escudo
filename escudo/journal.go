package escudo

import (
	"github.com/thiagola92/escudo/escudo/assert"
)

const (
	// When setted means that we need to create the lock file
	// (and maybe the temporary file too).
	//
	// This state is never found in the journal because is the
	// initial state, where we didn't accomplish nothing yet.
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
