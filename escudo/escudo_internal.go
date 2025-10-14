package escudo

import (
	"fmt"
	"path"
)

func (shield *Shield) lockpath() string {
	return path.Join(shield.path, "global"+lockext)
}

func (shield *Shield) journalspath() string {
	return path.Join(shield.path, "journals")
}

func (shield *Shield) journalpath(pid int) string {
	name := fmt.Sprintf("%d%s", pid, jourext)
	return path.Join(shield.journalspath(), name)
}
