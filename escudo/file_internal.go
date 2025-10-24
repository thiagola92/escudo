package escudo

import (
	"bytes"
	"crypto/sha256"
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
	file.temp, err = os.OpenFile(temppath, os.O_RDWR|os.O_CREATE, 0770)

	defer assert.Catch()
	assert.NoErr(err)

	info, err := file.temp.Stat()

	assert.NoErr(err)

	// To avoid damaging if has content.
	if info.Size() == 0 {
		assert.NoErrOn2(io.Copy(file.temp, file.orig))
	}
}

func (file *File) modified() bool {
	orig_current, err1 := file.orig.Seek(0, io.SeekCurrent)
	temp_current, err2 := file.temp.Seek(0, io.SeekCurrent)

	// Order matters.
	defer assert.Catch()
	assert.NoErr(err1)
	defer file.orig.Seek(orig_current, io.SeekStart)
	assert.NoErr(err2)
	defer file.temp.Seek(temp_current, io.SeekStart)

	assert.NoErrOn2(file.orig.Seek(0, io.SeekStart))
	assert.NoErrOn2(file.temp.Seek(0, io.SeekStart))

	orig_hash := sha256.New()
	temp_hash := sha256.New()

	assert.NoErrOn2(io.Copy(orig_hash, file.orig))
	assert.NoErrOn2(io.Copy(temp_hash, file.temp))

	diff := bytes.Compare(orig_hash.Sum(nil), temp_hash.Sum(nil))

	return diff != 0

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
		// Note: Anyone can get the lock between closing and removing,
		// but doing through journal we have the global lock to prevent this.
		assert.Closed(file.lock.Close())
		assert.FileNotExist(os.Remove(file.lockpath()))
	}

	file.orig = nil
	file.temp = nil
	file.dir = nil
	file.lock = nil
}
