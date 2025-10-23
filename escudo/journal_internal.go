package escudo

import (
	"encoding/json"
	"os"

	"github.com/thiagola92/escudo/escudo/assert"
)

func (journal *Journal) commit() (err error) {
	j, err := json.Marshal(journal)

	if err != nil {
		return err
	}

	err = journal.file.Truncate(0)

	if err != nil {
		return err
	}

	_, err = journal.file.Seek(0, 0)

	if err != nil {
		return err
	}

	_, err = journal.file.WriteString(string(j))

	if err != nil {
		return err
	}

	err = journal.file.Close()

	if err != nil {
		return err
	}

	journal.file = OpenFile(journal.path, os.O_RDWR|os.O_CREATE, 0770)
	journal.file.Lock()

	return nil
}

func (journal *Journal) remove() {
	assert.FileNotExist(os.Remove(journal.file.lockpath()))
	assert.FileNotExist(os.Remove(journal.file.temppath()))
	assert.FileNotExist(os.Remove(journal.file.path))

	for _, entry := range journal.Entries {
		assert.FileNotExist(os.Remove(entry.file.lockpath()))
		assert.FileNotExist(os.Remove(entry.file.temppath()))
	}
}
