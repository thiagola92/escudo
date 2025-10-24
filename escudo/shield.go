package escudo

import (
	"os"

	"github.com/thiagola92/escudo/escudo/assert"
)

type Shield struct {
	path string
}

// Get the journal of the current process.
func (shield *Shield) GetJournal() (journal *Journal, err error) {
	journalpath := shield.journalpath(os.Getpid())
	journal = shield.openJournal(journalpath)

	if assert.Err != nil {
		return nil, assert.Err
	}

	return journal, nil
}
