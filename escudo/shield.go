package escudo

import (
	"os"
)

type Shield struct {
	path string
}

// Get the journal of the current process.
func (shield *Shield) GetJournal() (*Journal, error) {
	journalpath := shield.journalpath(os.Getpid())
	return shield.openJournal(journalpath)
}
