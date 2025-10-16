package escudo

import (
	"os"

	"github.com/thiagola92/go-lockedfile/lockedfile"
)

func (shield *Shield) getlock() (*lockedfile.File, error) {
	// Note that this will lock until get the lock.
	lockpath := shield.lockpath()
	return lockedfile.OpenFile(lockpath, os.O_WRONLY, 0770)
}
