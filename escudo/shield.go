package escudo

import (
	"os"
)

type Shield struct {
	path string
}

// Get journal of the current process.
func (shield *Shield) GetJournal() (journal *Journal, err error) {
	journalpath := shield.journalpath(os.Getpid())
	journal, err = shield.openJournal(journalpath)

	if err != nil {
		return nil, err
	}

	return journal, nil
}
