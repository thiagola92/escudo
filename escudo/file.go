package escudo

import (
	"io/fs"
	"os"
	"path"

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

// Replace original file by the temporary file and
// create another temporary file so you can keep editing.
func (file *File) Commit() error {
	var err error

	err = file.commit()

	if err != nil {
		return err
	}

	err = file.openOrig()

	if err != nil {
		return err
	}

	return file.openTemp()
}

// Close the current File.
// Passing true will replace original file by the temporary file
// before closing everything.
func (file *File) Close(commit bool) error {
	if commit {
		err := file.commit()

		if err != nil {
			return err
		}
	}

	return file.close()
}

// It attempts to obtain the lock (it can fail).
func (file *File) Lock() error {
	var err error

	shared, err := file.openLock()

	if err != nil {
		return err
	}

	err = file.openOrig()

	if err != nil {
		return err
	}

	file.dir, err = os.Open(path.Dir(file.path))

	if err != nil {
		return err
	}

	// If is shared lock, can return sooner because we don't use temp file.
	if shared {
		return nil
	}

	return file.openTemp()
}

// Unlike os.OpenFile, this doesn't actually open file.
// It only setup the File for when you (or journal) decide to call File.Lock.
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

// Wrapper.
func (file *File) Seek(offset int64, whence int) (ret int64, err error) {
	if file.temp == nil {
		return file.orig.Seek(offset, whence)
	}

	return file.temp.Seek(offset, whence)
}

// Wrapper.
func (file *File) Truncate(size int64) error {
	if file.temp == nil {
		return nil
	}

	return file.temp.Truncate(size)
}

// Wrapper.
func (file *File) Write(b []byte) (int, error) {
	if file.temp == nil {
		return 0, nil
	}

	return file.temp.Write(b)
}

// Wrapper.
func (file *File) WriteString(s string) (int, error) {
	if file.temp == nil {
		return 0, nil
	}

	return file.temp.WriteString(s)
}
