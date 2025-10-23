package escudo

import (
	"io"
	"os"
	"path"

	"github.com/thiagola92/escudo/escudo/assert"
	"github.com/thiagola92/go-lockedfile/lockedfile"
)

func (file *File) openLock() (shared bool) {
	var err error

	shared = (file.flag & (os.O_WRONLY | os.O_RDWR)) == os.O_RDONLY
	lockpath := file.lockpath()

	if shared {
		file.lock, err = lockedfile.TryOpenFile(lockpath, os.O_RDONLY|os.O_CREATE, 0770)
	} else {
		file.lock, err = lockedfile.TryOpenFile(lockpath, os.O_RDWR|os.O_CREATE, 0770)
	}

	defer assert.Catch()
	assert.NoErr(err)

	return shared
}

func (file *File) openOrig() {
	var err error

	file.orig, err = os.OpenFile(file.path, file.flag, file.perm)

	defer assert.Catch()
	assert.NoErr(err)
}

func (file *File) openDir() {
	var err error

	file.dir, err = os.Open(path.Dir(file.path))

	defer assert.Catch()
	assert.NoErr(err)
}

func (file *File) openTemp() {
	var err error

	temppath := file.temppath()
	file.temp, err = os.OpenFile(temppath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0770)

	defer assert.Catch()
	assert.NoErr(err)
	assert.NoErrOn2(io.Copy(file.temp, file.orig))
}

func (file *File) replace() {
	if file.orig == nil || file.dir == nil || file.temp == nil {
		return
	}

	defer assert.Catch()
	assert.NoErr(file.temp.Sync())
	assert.NoErr(file.dir.Sync())
	assert.Closed(file.orig.Close())
	assert.Closed(file.temp.Close())
	assert.NoErr(os.Rename(file.temppath(), file.path))
	assert.NoErr(file.dir.Sync())

	file.orig = nil
	file.temp = nil
}

func (file *File) close() {
	defer assert.Catch()

	if file.orig != nil {
		assert.Closed(file.orig.Close())
	}

	if file.temp != nil {
		assert.Closed(file.temp.Close())
		assert.FileNotExist(os.Remove(file.temppath()))
	}

	if file.dir != nil {
		assert.Closed(file.dir.Close())
	}

	if file.lock != nil {
		assert.Closed(file.lock.Close())
		// We don't delete lock HERE because we don't know
		// if someone may obtain it before we delete it.
	}

	file.orig = nil
	file.temp = nil
	file.dir = nil
	file.lock = nil
}
