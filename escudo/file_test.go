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
	// Setup.
	file := NewFile(test_filename, os.O_RDONLY, 0770)
	file2 := NewFile(test_filename, os.O_WRONLY, 0770)
	file3 := NewFile(test_filename, os.O_RDWR, 0770)

	// Test.
	if (file.flag & os.O_RDONLY) != os.O_RDONLY {
		t.Errorf("Flag should be os.O_RDONLY, instead is: %d", file.flag)
	}

	if (file2.flag & os.O_RDWR) != os.O_RDWR {
		t.Errorf("Flag should have changed to os.O_RDWR, but is: %d", file2.flag)
	}

	if (file3.flag & os.O_RDWR) != os.O_RDWR {
		t.Errorf("Flag should be os.O_RDWR, but is: %d", file3.flag)
	}
}

func TestFileExclusiveLock(t *testing.T) {
	// Setup.
	file := NewFile(test_filename, os.O_RDWR|os.O_CREATE, 0770)
	file2 := NewFile(test_filename, os.O_RDWR, 0770)

	defer file.Close()
	defer file2.Close()

	// Test.
	err := file.Lock()

	if err != nil {
		t.Errorf("Failed to lock file: %s", err.Error())
	}

	err = file2.Lock()

	if err == nil {
		t.Errorf("Should have failed to obtain an already locked file")
	}
}

func TestFileSharedLock(t *testing.T) {
	// Setup.
	file := NewFile(test_filename, os.O_RDONLY|os.O_CREATE, 0770)
	file2 := NewFile(test_filename, os.O_RDONLY, 0770)

	defer file.Close()
	defer file2.Close()

	// Test.
	err := file.Lock()

	if err != nil {
		t.Errorf("Failed to lock file: %s", err.Error())
	}

	err = file2.Lock()

	if err != nil {
		t.Errorf("Should have obtain locked file because is a shared lock")
	}
}

func TestFileWriteString(t *testing.T) {
	// Setup.
	file := NewFile(test_filename, os.O_RDWR|os.O_CREATE, 0770)

	defer file.Close()
	assert_t.NoErr(t, file.Lock())

	// Test.
	written, err := file.WriteString(test_str0)

	if written != len(test_str0) {
		t.Error("Failed to write all bytes from string")
	}

	if err != nil {
		t.Errorf("Failed to write bytes: %s", err.Error())
	}

	str, err := os.ReadFile(file.temp.Name()) // Read temporary file.

	assert_t.NoErr(t, err)

	if strings.Compare(string(str), test_str0) != 0 {
		t.Errorf("Temporary file is not holding the expected content")
	}
}

func TestFileCommit(t *testing.T) {
	// Setup.
	assert_t.NoErr(t, os.WriteFile(test_filename, []byte(test_str0), 0770))

	file := NewFile(test_filename, os.O_RDWR|os.O_CREATE, 0770)

	defer file.Close()
	assert_t.NoErr(t, file.Lock())

	_, err := file.WriteString(test_str1)

	assert_t.NoErr(t, err)

	// Test.
	err = file.Commit()

	if err != nil {
		t.Errorf("Failed to commit changes: %s", err.Error())
	}

	str, err := os.ReadFile(file.orig.Name()) // Read original file.

	assert_t.NoErr(t, err)

	if strings.Compare(string(str), test_str1) != 0 {
		t.Errorf("Orignal file is not holding the expected content")
	}
}
