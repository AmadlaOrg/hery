package file

import (
	"io"
	"os"
	"syscall"
	"time"
)

// IFile is the interface for the methods found in os/file
type IFile interface {
	Name() string
	Read(b []byte) (n int, err error)
	ReadAt(b []byte, off int64) (n int, err error)
	ReadFrom(r io.Reader) (n int64, err error)
	Write(b []byte) (n int, err error)
	WriteAt(b []byte, off int64) (n int, err error)
	WriteTo(w io.Writer) (n int64, err error)
	Seek(offset int64, whence int) (ret int64, err error)
	WriteString(s string) (n int, err error)
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	SyscallConn() (syscall.RawConn, error)
	Chmod(mode os.FileMode) error
	Chown(uid, gid int) error
	Close() error
	Fd() uintptr
	ReadDir(n int) ([]os.DirEntry, error)
	Readdirnames(n int) (names []string, err error)
	Sync() error
	Truncate(size int64) error
	Chdir() error
}
