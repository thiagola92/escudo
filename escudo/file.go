package escudo

import (
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/gofrs/flock"
)

type File struct {
	path string
	flag int
	perm fs.FileMode

	orig *os.File     // Original file
	dir  *os.File     // File directory
	lock *flock.Flock // Lock file
	temp *os.File     // Temporary file
}

func OpenFile(filepath string, flag int, perm fs.FileMode) *File {
	file := &File{path: filepath}

	// Replace WRONLY with RDWR, otherwise we can't create
	// the temporary file with the same content as the original.
	if (flag & os.O_WRONLY) == os.O_WRONLY {
		flag = (flag ^ os.O_WRONLY) | os.O_RDWR
	}

	file.path = filepath
	file.flag = flag
	file.perm = perm

	return file
}

func (file *File) ShieldUp() (bool, error) {
	var err error

	// Get original file.
	file.orig, err = os.OpenFile(file.path, file.flag, file.perm)

	if err != nil {
		return false, err
	}

	// Get lock file.
	lockpath := file.lockpath()
	file.lock = flock.New(lockpath, flock.SetFlag(os.O_RDWR|os.O_CREATE), flock.SetPermissions(0770))

	if (file.flag & (os.O_WRONLY | os.O_RDWR)) == os.O_RDONLY {
		locked, err := file.lock.TryRLock()

		if (err != nil) || !locked {
			return false, err
		}

		// Return sooner because doesn't need temporary file for shared locks.
		return true, nil
	} else {
		locked, err := file.lock.TryLock()

		if (err != nil) || !locked {
			return false, err
		}
	}

	// Get directory
	file.dir, err = os.Open(path.Dir(file.path))

	if err != nil {
		return false, err
	}

	// Get temporary file.
	temppath := file.temppath()
	file.temp, err = os.OpenFile(temppath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0770)

	if err != nil {
		return false, err
	}

	_, err = io.Copy(file.temp, file.orig)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (file *File) Close() error {
	var err error

	if file.orig != nil {
		err = file.orig.Close()

		if err != nil {
			return err
		}
	}

	if file.temp != nil {
		err = file.temp.Sync()

		if err != nil {
			return err
		}

		err = file.dir.Sync()

		if err != nil {
			return err
		}

		err = file.temp.Close()

		if err != nil {
			return err
		}

		err = file.dir.Close()

		if err != nil {
			return err
		}

		err = os.Rename(file.temppath(), file.path)

		if err != nil {
			return err
		}
	}

	if file.lock != nil {
		err = file.lock.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
