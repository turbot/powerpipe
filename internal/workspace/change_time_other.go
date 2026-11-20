//go:build !darwin && !linux

package workspace

import (
	"os"
	"time"
)

// changeTime returns the modification time of the file: this system has no status change time
func changeTime(info os.FileInfo) time.Time {
	return info.ModTime()
}
