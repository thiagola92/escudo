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
	assert.Catch(&err)

	journal.Entries, err = newJournalEntries(files)

	if err != nil {
		return err
	}

	lock, err := journal.shield.waitLock()

	if err != nil {
		return err
	}

	defer lock.Close()

	err = journal.lockall()

	if err != nil {
		return err
	}

	journal.Status = WRITING
	err = journal.replace()

	return err
}

// Commit changes to files and journal.
func (journal *Journal) Commit() (err error) {
	assert.Catch(&err)

	lock, err := journal.shield.waitLock()

	if err != nil {
		return err
	}

	defer lock.Close()

	journal.Status = REPLACING
	err = journal.replace()

	if err != nil {
		return err
	}

	err = journal.commitall()

	if err != nil {
		return err
	}

	journal.Status = WRITING
	err = journal.replace()

	return err
}

// Delete temporary files and journal.
func (journal *Journal) Close() (err error) {
	lock, err := journal.shield.waitLock()

	if err != nil {
		return err
	}

	defer lock.Close()

	journal.Status = DELETING
	err = journal.replace()

	if err != nil {
		return err
	}

	err = journal.closeall()

	if err != nil {
		return err
	}

	err = journal.close()

	return err
}
