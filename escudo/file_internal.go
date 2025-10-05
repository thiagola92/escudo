package escudo

import (
	"path"
	"strings"
)

func (file *File) lockpath() string {
	var lockname strings.Builder

	lockname.WriteString(file.orig.Name())
	lockname.WriteString(LockExt)

	return path.Join(path.Dir(file.path), lockname.String())
}

func (file *File) temppath() string {
	var tempname strings.Builder

	tempname.WriteString(file.orig.Name())
	tempname.WriteString(TempExt)

	return path.Join(path.Dir(file.path), tempname.String())
}
