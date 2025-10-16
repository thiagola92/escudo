package escudo

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/thiagola92/go-lockedfile/lockedfile"
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

	journal.Entries, err = toJournalEntries(files)

	if err != nil {
		return err
	}

	err = journal.save()

	if err != nil {
		return err
	}

	for index, entry := range journal.Entries {
		err := entry.file.Lock()

		if err != nil {
			for index >= 0 {
				entry.file.Close()
				entry.Status = LOCKING
				index -= 1
			}

			return err
		}

		entry.Status = WRITING
	}

	err = journal.save()

	if err != nil {
		return err
	}

	return nil
}

// - TODO: Remove locks from files? Otherwise the directory may feel dirty (many lock files).
func (journal *Journal) Close() error {
	lockpath := journal.shield.lockpath()
	lock, err := lockedfile.OpenFile(lockpath, os.O_WRONLY, 0770)

	if err != nil {
		return err
	}

	defer lock.Close()
	journal.file.Close()

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

func (journal *Journal) save() error {
	j, err := json.Marshal(journal)

	if err != nil {
		return err
	}

	lockpath := journal.shield.lockpath()
	lock, err := lockedfile.OpenFile(lockpath, os.O_WRONLY, 0770)

	if err != nil {
		return err
	}

	defer lock.Close()

	journal.file.Truncate(0)
	journal.file.Seek(0, 0)
	journal.file.WriteString(string(j))
	journal.file.Close()
	journal.file = OpenFile(journal.path, os.O_RDWR|os.O_CREATE, 0770)
	err = journal.file.Lock()

	return err
}

func (journal *Journal) remove() error {
	var err error

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

	err = os.Remove(journal.path)

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return nil
}
