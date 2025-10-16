package escudo

import (
	"fmt"
	"path"
)

func (file *File) lockpath() string {
	return path.Join(path.Dir(file.path), path.Base(file.path)+lockext)
}

func (file *File) temppath() string {
	return path.Join(path.Dir(file.path), path.Base(file.path)+tempext)
}

func (shield *Shield) lockpath() string {
	return path.Join(shield.path, "global"+lockext)
}

func (shield *Shield) journalspath() string {
	return path.Join(shield.path, "journals")
}

func (shield *Shield) journalpath(pid int) string {
	name := fmt.Sprintf("%d%s", pid, jourext)
	journalspath := shield.journalspath()
	return path.Join(journalspath, name)
}
