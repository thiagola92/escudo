package escudo

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"strings"
)

type Journal struct {
	path   string
	file   *File
	shield *Shield

	Entries []*JournalEntry
}

func (shield *Shield) MyJournal() (*Journal, error) {
	journalpath := shield.journalpath(os.Getpid())
	return shield.openJournal(journalpath)
}

func (journal *Journal) LockFiles(files ...*File) error {
	var err error

	// Setup journal entries.
	journal.Entries, err = toJournalEntries(files)

	if err != nil {
		return err
	}

	// Hold global lock until the end of this function.
	lock, err := journal.shield.getlock()

	if err != nil {
		return err
	}

	defer lock.Close()

	// Attempt to get all locks.
	for index, entry := range journal.Entries {
		err := entry.file.Lock()

		if err != nil {
			for index >= 0 {
				entry.file.Close()
				entry.Status = INITIALIZING
				index -= 1
			}

			return err
		}

		entry.Status = WRITING
	}

	// Save state of journal.
	err = journal.save(true)

	if err != nil {
		return err
	}

	return nil
}

func (journal *Journal) Save() error {
	// Hold global lock until the end of this function.
	lock, err := journal.shield.getlock()

	if err != nil {
		return err
	}

	defer lock.Close()

	// Update state of journal.
	for _, entry := range journal.Entries {
		entry.Status = REPLACING
	}

	err = journal.save(true)

	if err != nil {
		return err
	}

	// TODO: Make the actual replacement of files.
	// Right now I'm just updating the journal.

	// Update state of journal.
	for _, entry := range journal.Entries {
		entry.Status = WRITING
	}

	err = journal.save(true)

	if err != nil {
		return err
	}

	return nil
}

func (journal *Journal) Close() error {
	// Hold global lock until the end of this function.
	lock, err := journal.shield.getlock()

	if err != nil {
		return err
	}

	defer lock.Close()

	// Update state of journal.
	for _, entry := range journal.Entries {
		entry.Status = DELETING
	}

	err = journal.save(false)

	if err != nil {
		return err
	}

	// Remove journal files.
	return journal.remove()
}

func (shield *Shield) anyJournal() (*Journal, error) {
	var journal *Journal
	var err error

	journalspath := shield.journalspath()
	dir, err := os.Open(journalspath)

	if err != nil {
		return nil, err
	}

	defer dir.Close()

	files, err := dir.ReadDir(0)

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), lockext) {
			continue
		}

		journalpath := path.Join(journalspath, file.Name())
		journal, err = shield.openJournal(journalpath)

		if err == nil {
			return journal, nil
		}
	}

	return nil, nil
}

func (shield *Shield) openJournal(journalpath string) (*Journal, error) {
	var err error

	journal := &Journal{path: journalpath, shield: shield}
	journal.file = OpenFile(journalpath, os.O_RDWR|os.O_CREATE, 0770)
	err = journal.file.Lock()

	if err != nil {
		return nil, err
	}

	return journal, nil
}

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

	err = journal.file.Close()

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
