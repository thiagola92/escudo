package escudo

import (
	"encoding/json"
	"os"

	"github.com/thiagola92/escudo/escudo/assert"
)

func (journal *Journal) lockall() {
	defer assert.Catch()

	for index, entry := range journal.Entries {
		err := entry.file.Lock()

		// Attempts to undo all locks in case one fail.
		if err != nil {
			for index >= 0 {
				entry.file.Close()
				index -= 1
			}

			assert.NoErr(err)
		}
	}
}

func (journal *Journal) commitall() {
	defer assert.Catch()

	for _, entry := range journal.Entries {
		assert.NoErr(entry.file.Commit())
	}
}

func (journal *Journal) closeall() {
	defer assert.Catch()

	for _, entry := range journal.Entries {
		assert.NoErr(entry.file.Close())
	}
}

func (journal *Journal) replace() {
	j, err := json.Marshal(journal)

	defer assert.Catch()
	assert.NoErr(err)
	assert.NoErr(journal.file.Truncate(0))
	assert.NoErrOn2(journal.file.Seek(0, 0))
	assert.NoErrOn2(journal.file.WriteString(string(j)))
	assert.NoErr(journal.file.Commit())
}

func (journal *Journal) close() {
	defer assert.Catch()
	assert.NoErr(journal.file.Close())
	assert.FileNotExist(os.Remove(journal.file.path))
}

func (journal *Journal) resume() {
	defer assert.Catch()

	switch journal.Status {
	case WRITING:
		assert.NoErr(journal.Close())
	case REPLACING:
		assert.NoErr(journal.Commit())
	case DELETING:
		assert.NoErr(journal.Close())
	default:
		return
	}
}
