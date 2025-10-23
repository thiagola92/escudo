package escudo

import (
	"encoding/json"
	"os"

	"github.com/thiagola92/escudo/escudo/assert"
)

func (journal *Journal) replace() {
	var err error

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
	assert.FileNotExist(os.Remove(journal.file.lockpath()))
	assert.FileNotExist(os.Remove(journal.file.temppath()))
	assert.FileNotExist(os.Remove(journal.file.path))

	for _, entry := range journal.Entries {
		assert.NoErr(entry.file.Close())
		assert.FileNotExist(os.Remove(entry.file.lockpath()))
	}
}
