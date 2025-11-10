package escudo

import (
	"os"
	"testing"
)

func TestFileFlag(t *testing.T) {
	file := NewFile("../test/file.txt", os.O_RDONLY, 0770)

	if (file.flag & os.O_RDONLY) != os.O_RDONLY {
		t.Errorf("Flag should be os.O_RDONLY, instead is: %d", file.flag)
	}

	file = NewFile("../test/file.txt", os.O_WRONLY, 0770)

	if (file.flag & os.O_RDWR) != os.O_RDWR {
		t.Errorf("Flag should have changed to os.O_RDWR, but is: %d", file.flag)
	}

	file = NewFile("../test/file.txt", os.O_RDWR, 0770)

	if (file.flag & os.O_RDWR) != os.O_RDWR {
		t.Errorf("Flag should be os.O_RDWR, but is: %d", file.flag)
	}
}

func TestFileExclusiveLock(t *testing.T) {
	file := NewFile("../test/file.txt", os.O_RDWR|os.O_CREATE, 0770)
	err := file.Lock()
	defer file.Close()

	if err != nil {
		t.Errorf("Failed to lock file: %s", err.Error())
	}

	file2 := NewFile("../test/file.txt", os.O_RDWR, 0770)
	err = file2.Lock()
	defer file2.Close()

	if err == nil {
		t.Errorf("Should have failed to obtain an already locked file")
	}
}

func TestFileSharedLock(t *testing.T) {
	file := NewFile("../test/file.txt", os.O_RDONLY|os.O_CREATE, 0770)
	err := file.Lock()
	defer file.Close()

	if err != nil {
		t.Errorf("Failed to lock file: %s", err.Error())
	}

	file2 := NewFile("../test/file.txt", os.O_RDONLY, 0770)
	err = file2.Lock()
	defer file2.Close()

	if err != nil {
		t.Errorf("Should have obtain locked file because is a shared lock")
	}
}

func TestFileCommit(t *testing.T) {
	var START_STR = "Start"
	var END_STR = "End"
	var str = make([]byte, 256)
	var err error
	var written int
	var file *File

	os.WriteFile("../test/file.txt", []byte(START_STR), 0770)
	file = NewFile("../test/file.txt", os.O_RDWR|os.O_CREATE, 0770)
	file.Lock()
	defer file.Close()

	written, err = file.WriteString(END_STR)

	if written == 0 {
		t.Error("Failed to write bytes")
	}

	if err != nil {
		t.Errorf("Failed to write bytes: %s", err.Error())
	}

	err = file.Commit()

	if err != nil {
		t.Errorf("Failed to commit: %s", err.Error())
	}

	written, err = file.Read(str)

	if written == 0 {
		t.Error("Failed to read bytes")
	}

	if err != nil {
		t.Errorf("Failed to read bytes: %s", err.Error())
	}

	if string(str) == END_STR {
		t.Errorf("Should have changed file")
	}
}
