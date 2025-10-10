package escudo

import (
	"fmt"
	"path"
)

func (escudo *Escudo) lockpath() string {
	return path.Join(escudo.path, "global"+lockext)
}

func (escudo *Escudo) journalspath() string {
	return path.Join(escudo.path, "journals")
}

func (escudo *Escudo) journalpath(pid int) string {
	name := fmt.Sprintf("%d%s", pid, jourext)
	return path.Join(escudo.journalspath(), name)
}
