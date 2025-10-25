package escudo

import (
	"io/fs"
	"os"

	"github.com/thiagola92/escudo/escudo/assert"
	"github.com/thiagola92/go-lockedfile/lockedfile"
)

type File struct {
	path string
	flag int
	perm fs.FileMode

	orig *os.File         // Original file
	dir  *os.File         // File directory
	lock *lockedfile.File // Lock file
	temp *os.File         // Temporary file
}

// Setup file to be open in the future.
func NewFile(filepath string, flag int, perm fs.FileMode) *File {
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

// Obtain permission to interact with file.
func (file *File) Lock() (err error) {
	shared := file.openLock()

	if assert.Err != nil {
		return assert.Err
	}

	file.openOrig()

	if assert.Err != nil {
		return assert.Err
	}

	file.openDir()

	if assert.Err != nil {
		return assert.Err
	}

	// Only need temporary file if is an exclusive lock.
	if !shared {
		file.openTemp()
	}

	return assert.Err
}

// Commit changes to file.
func (file *File) Commit() error {
	file.replace()

	if assert.Err != nil {
		return assert.Err
	}

	file.openOrig()

	if assert.Err != nil {
		return assert.Err
	}

	file.openTemp()

	return assert.Err
}

// Delete temporary file and lock.
func (file *File) Close() error {
	file.close()

	return assert.Err
}

///////////////////////////////////////////////////////////
//////////////////////// Wrappers
///////////////////////////////////////////////////////////

func (file *File) Seek(offset int64, whence int) (ret int64, err error) {
	if file.temp == nil {
		return file.orig.Seek(offset, whence)
	}

	return file.temp.Seek(offset, whence)
}

func (file *File) Truncate(size int64) error {
	if file.temp == nil {
		return nil
	}

	return file.temp.Truncate(size)
}

func (file *File) Write(b []byte) (int, error) {
	if file.temp == nil {
		return 0, nil
	}

	return file.temp.Write(b)
}

func (file *File) WriteString(s string) (int, error) {
	if file.temp == nil {
		return 0, nil
	}

	return file.temp.WriteString(s)
}
