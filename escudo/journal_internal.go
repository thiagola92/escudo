package escudo

import (
	"encoding/json"
	"errors"
	"os"
)

func (journal *Journal) save(reopen bool) error {
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

	err = journal.file.Close(true)

	if err != nil {
		return err
	}

	if !reopen {
		return nil
	}

	journal.file = OpenFile(journal.path, os.O_RDWR|os.O_CREATE, 0770)
	err = journal.file.Lock()

	return err
}

func (journal *Journal) remove() error {
	var err error

	// Remove journal files.
	err = os.Remove(journal.file.lockpath())

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	err = os.Remove(journal.file.temppath())

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	err = os.Remove(journal.file.path)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	// Remove lock files from original files.
	for _, entry := range journal.Entries {
		err = os.Remove(entry.file.lockpath())

		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	return nil
}
