package escudo

import (
	"io"
	"os"

	"github.com/thiagola92/go-lockedfile/lockedfile"
)

func (file *File) close() error {
	var err error

	if file.orig != nil {
		err = file.orig.Close()

		if err != nil {
			return err
		}

		file.orig = nil
	}

	if file.temp != nil {
		err = file.temp.Close()

		if err != nil {
			return err
		}

		file.temp = nil
	}

	if file.dir != nil {
		err = file.dir.Close()

		if err != nil {
			return err
		}

		file.dir = nil
	}

	if file.lock != nil {
		err = file.lock.Close()

		if err != nil {
			return err
		}

		file.lock = nil
	}

	return nil
}

func (file *File) commit() error {
	var err error

	if file.temp == nil {
		return nil
	}

	err = file.temp.Sync()

	if err != nil {
		return err
	}

	err = file.dir.Sync()

	if err != nil {
		return err
	}

	err = os.Rename(file.temppath(), file.path)

	if err != nil {
		return err
	}

	err = file.dir.Sync()

	if err != nil {
		return err
	}

	return nil
}

func (file *File) openLock() (bool, error) {
	var err error

	shared := (file.flag & (os.O_WRONLY | os.O_RDWR)) == os.O_RDONLY
	lockpath := file.lockpath()

	if shared {
		file.lock, err = lockedfile.TryOpenFile(lockpath, os.O_RDONLY|os.O_CREATE, 0770)
	} else {
		file.lock, err = lockedfile.TryOpenFile(lockpath, os.O_RDWR|os.O_CREATE, 0770)
	}

	return shared, err
}

func (file *File) openOrig() error {
	var err error

	// In case we are reopening.
	if file.orig != nil {
		err = file.orig.Close()

		if err != nil {
			return err
		}
	}

	file.orig, err = os.OpenFile(file.path, file.flag, file.perm)

	return err
}

func (file *File) openTemp() error {
	var err error

	// In case we are reopening.
	if file.temp != nil {
		err = file.temp.Close()

		if err != nil {
			return err
		}
	}

	temppath := file.temppath()
	file.temp, err = os.OpenFile(temppath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0770)

	if err != nil {
		return err
	}

	_, err = io.Copy(file.temp, file.orig)

	if err != nil {
		return err
	}

	return nil
}
