package escudo

import (
	"encoding/json"
	"os"

	"github.com/thiagola92/escudo/escudo/assert"
)

func (journal *Journal) lockall() (err error) {
	defer assert.Catch(&err)

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

	return nil
}

func (journal *Journal) commitall() (err error) {
	defer assert.Catch(&err)

	for _, entry := range journal.Entries {
		assert.NoErr(entry.file.Commit())
	}

	return nil
}

func (journal *Journal) closeall() (err error) {
	defer assert.Catch(&err)

	for _, entry := range journal.Entries {
		assert.NoErr(entry.file.Close())
	}

	return nil
}

func (journal *Journal) replace() (err error) {
	defer assert.Catch(&err)

	j, err := json.Marshal(journal)

	assert.NoErr(err)
	assert.NoErr(journal.file.Truncate(0))
	assert.NoErrOn2(journal.file.Seek(0, 0))
	assert.NoErrOn2(journal.file.WriteString(string(j)))
	assert.NoErr(journal.file.Commit())

	return nil
}

func (journal *Journal) close() (err error) {
	defer assert.Catch(&err)

	assert.NoErr(journal.file.Close())
	assert.FileNotExist(os.Remove(journal.file.path))

	return nil
}

func (journal *Journal) resume() (err error) {
	defer assert.Catch(&err)

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

	return nil
}
