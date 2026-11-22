package workspace

import (
	"os"
	"syscall"
	"time"
)

// changeTime returns the status change time (ctime) of the file, or its modification time if it has none
func changeTime(info os.FileInfo) time.Time {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return info.ModTime()
	}
	return time.Unix(stat.Ctimespec.Unix())
}
