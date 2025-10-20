package escudo

type Journal struct {
	path   string
	file   *File
	shield *Shield

	Entries []*JournalEntry
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
				entry.file.Close(true)
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
