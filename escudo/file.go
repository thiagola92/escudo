package escudo

import (
	"cmp"
	"io/fs"
	"os"

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

func (file *File) Commit() (err error) {
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

func (file *File) Push() (err error) {
	err = file.commit()

	if err != nil {
		return err
	}

	return file.close()
}

func (file *File) Lock() (err error) {
	shared, err := file.openLock()

	err = cmp.Or(
		err,
		file.openOrig(),
		file.openDir(),
	)

	if err != nil {
		return err
	}

	// Only need temporary file if is an exclusive lock.
	if !shared {
		return file.openTemp()
	}

	return nil
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

///////////////////////////////////////////////////////////
//////////////////////// Wrappers
///////////////////////////////////////////////////////////

func (file *File) Close() error {
	file.close()

	return nil
}

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
