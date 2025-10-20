package escudo

import (
	"fmt"
	"path"
)

// filename.escl
func (file *File) lockpath() string {
	return path.Join(path.Dir(file.path), path.Base(file.path)+lockext)
}

// filename.esct
func (file *File) temppath() string {
	return path.Join(path.Dir(file.path), path.Base(file.path)+tempext)
}

// .escudo/global.escl
func (shield *Shield) lockpath() string {
	return path.Join(shield.path, "global"+lockext)
}

// .escudo/journals/
func (shield *Shield) journalspath() string {
	return path.Join(shield.path, "journals")
}

// .escudo/journals/1234.escj
func (shield *Shield) journalpath(pid int) string {
	name := fmt.Sprintf("%d%s", pid, jourext)
	journalspath := shield.journalspath()
	return path.Join(journalspath, name)
}
