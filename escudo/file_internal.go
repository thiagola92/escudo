package escudo

import (
	"cmp"
	"io"
	"os"
	"path"

	"github.com/thiagola92/go-lockedfile/lockedfile"
)

func (file *File) close() (err error) {
	if file.orig != nil {
		err = cmp.Or(err, file.orig.Close())
	}

	if file.temp != nil {
		err = cmp.Or(
			err,
			file.temp.Close(),
			os.Remove(file.temppath()),
		)
	}

	if file.dir != nil {
		err = cmp.Or(err, file.dir.Close())
	}

	if file.lock != nil {
		err = cmp.Or(err, file.lock.Close())
	}

	file.orig = nil
	file.temp = nil
	file.dir = nil
	file.lock = nil

	return err
}

func (file *File) commit() (err error) {
	if file.orig == nil || file.dir == nil || file.temp == nil {
		return nil
	}

	err = cmp.Or(
		file.temp.Sync(),
		file.dir.Sync(),
		os.Rename(file.temppath(), file.path),
		file.dir.Sync(),
	)

	file.orig.Close()
	file.temp.Close()

	file.orig = nil
	file.temp = nil

	return err
}

func (file *File) openDir() (err error) {
	file.dir, err = os.Open(path.Dir(file.path))
	return err
}

func (file *File) openLock() (shared bool, err error) {
	shared = (file.flag & (os.O_WRONLY | os.O_RDWR)) == os.O_RDONLY
	lockpath := file.lockpath()

	if shared {
		file.lock, err = lockedfile.TryOpenFile(lockpath, os.O_RDONLY|os.O_CREATE, 0770)
	} else {
		file.lock, err = lockedfile.TryOpenFile(lockpath, os.O_RDWR|os.O_CREATE, 0770)
	}

	return shared, err
}

func (file *File) openOrig() (err error) {
	file.orig, err = os.OpenFile(file.path, file.flag, file.perm)
	return err
}

func (file *File) openTemp() (err error) {
	temppath := file.temppath()
	file.temp, err = os.OpenFile(temppath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0770)

	if err != nil {
		return err
	}

	_, err = io.Copy(file.temp, file.orig)

	return err
}

func (file *File) sync() (err error) {
	return cmp.Or(
		file.temp.Sync(),
		file.dir.Sync(),
	)
}
