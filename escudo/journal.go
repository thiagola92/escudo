package escudo

import "github.com/thiagola92/escudo/escudo/assert"

type Journal struct {
	path   string
	file   *File
	shield *Shield

	Entries []*JournalEntry
	Status  int
}

func (journal *Journal) Commit() {
	defer assert.Catch()

	lock := journal.shield.getLock()

	defer lock.Close()

	journal.Status = REPLACING
	journal.commit(true)

	// TODO: Make the actual replacement of files.
	// Right now I'm just updating the journal.

	journal.Status = WRITING
	journal.commit(true)
}

func (journal *Journal) Close(commit bool) {
	defer assert.Catch()

	lock := journal.shield.getLock()

	defer lock.Close()

	if commit {
		journal.Status = REPLACING
		journal.commit(true)
	}

	// TODO: Make the actual replacement of files.
	// Right now I'm just updating the journal.

	journal.Status = DELETING
	journal.commit(false)
	journal.remove()
}

func (journal *Journal) LockFiles(files ...*File) {
	defer assert.Catch()

	journal.Entries = newJournalEntries(files)
	lock := journal.shield.getLock()

	defer lock.Close()

	// Attempt to get all locks.
	for index, entry := range journal.Entries {
		entry.file.Lock()

		if assert.Err != nil {
			for index >= 0 {
				entry.file.Close()
				index -= 1
			}

			return
		}
	}

	journal.Status = WRITING
	journal.commit(true)
}
