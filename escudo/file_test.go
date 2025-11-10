package escudo

import (
	"os"
	"strings"
	"testing"

	"github.com/thiagola92/escudo/escudo/assert_t"
)

const test_filename = "../test/file.txt"
const test_str0 = "STRING 0"
const test_str1 = "STRING 1"

func TestFileFlag(t *testing.T) {
	file := NewFile(test_filename, os.O_RDONLY, 0770)

	if (file.flag & os.O_RDONLY) != os.O_RDONLY {
		t.Errorf("Flag should be os.O_RDONLY, instead is: %d", file.flag)
	}

	file = NewFile(test_filename, os.O_WRONLY, 0770)

	if (file.flag & os.O_RDWR) != os.O_RDWR {
		t.Errorf("Flag should have changed to os.O_RDWR, but is: %d", file.flag)
	}

	file = NewFile(test_filename, os.O_RDWR, 0770)

	if (file.flag & os.O_RDWR) != os.O_RDWR {
		t.Errorf("Flag should be os.O_RDWR, but is: %d", file.flag)
	}
}

func TestFileExclusiveLock(t *testing.T) {
	file := NewFile(test_filename, os.O_RDWR|os.O_CREATE, 0770)
	err := file.Lock()

	defer file.Close()

	if err != nil {
		t.Errorf("Failed to lock file: %s", err.Error())
	}

	file2 := NewFile(test_filename, os.O_RDWR, 0770)
	err = file2.Lock()

	defer file2.Close()

	if err == nil {
		t.Errorf("Should have failed to obtain an already locked file")
	}
}

func TestFileSharedLock(t *testing.T) {
	file := NewFile(test_filename, os.O_RDONLY|os.O_CREATE, 0770)
	err := file.Lock()

	defer file.Close()

	if err != nil {
		t.Errorf("Failed to lock file: %s", err.Error())
	}

	file2 := NewFile(test_filename, os.O_RDONLY, 0770)
	err = file2.Lock()

	defer file2.Close()

	if err != nil {
		t.Errorf("Should have obtain locked file because is a shared lock")
	}
}

func TestFileWriteString(t *testing.T) {
	file := NewFile(test_filename, os.O_RDWR|os.O_CREATE, 0770)

	defer file.Close()
	assert_t.NoErr(t, file.Lock())

	written, err := file.WriteString(test_str0)

	if written != len(test_str0) {
		t.Error("Failed to write all bytes from string")
	}

	if err != nil {
		t.Errorf("Failed to write bytes: %s", err.Error())
	}

	// Check if temporary file changed because no commit was made to change the main file.
	str, err := os.ReadFile(file.temp.Name())

	assert_t.NoErr(t, err)

	if strings.Compare(string(str), test_str0) != 0 {
		t.Errorf("Temporary not holding the expected content")
	}
}

// func TestFileCommit(t *testing.T) {
// 	assert_t.NoErr(t, os.WriteFile(test_filename, []byte("START"), 0770))

// 	file := NewFile(test_filename, os.O_RDWR|os.O_CREATE, 0770)

// 	defer file.Close()
// 	assert_t.NoErr(t, file.Lock())

// 	_, err := file.WriteString("END")
// 	assert_t.NoErr(t, err)

// 	// str := make([]byte, 256)
// 	// written, err = file.Read(str)

// 	// if written == 0 {
// 	// 	t.Error("Failed to read bytes")
// 	// }

// 	// if err != nil {
// 	// 	t.Errorf("Failed to read bytes: %s", err.Error())
// 	// }

// 	// if string(str) == "END" {
// 	// 	t.Errorf("Should have changed file")
// 	// }

// 	// err = file.Commit()

// 	// if err != nil {
// 	// 	t.Errorf("Failed to commit: %s", err.Error())
// 	// }

// 	// str = make([]byte, 256)
// 	// written, err = file.Read(str)

// 	// if written == 0 {
// 	// 	t.Error("Failed to read bytes")
// 	// }

// 	// if err != nil {
// 	// 	t.Errorf("Failed to read bytes: %s", err.Error())
// 	// }

// 	// if string(str) == "END" {
// 	// 	t.Errorf("Should have changed file")
// 	// }
// }
