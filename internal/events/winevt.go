/// Package winevt deals with getting Windows Event Log information

package events

import (
	"fmt"
	"runtime"

	"github.com/0xNF/gdz/internal/fs"
)

// GetEventLogs returns paths of event logs, which can be added to thew big zip file
// returns (valid paths, missing files)
func GetEventLogs() ([]string, []string) {
	var actual = []string{}
	var missing = []string{}
	if runtime.GOOS == "windows" {
		var getThese = []string{
			"Application.evtx",
			"Systen.evtx",
		}
		wpath := "C:\\Windows\\System32\\winevt\\Logs\\%s"
		for _, val := range getThese {
			fpath := fmt.Sprintf(wpath, val)
			if fs.Exists(fpath) {
				actual = append(actual, fpath)
			} else {
				missing = append(missing, fpath)
			}
		}
	}
	return actual, missing
}

// GetLinuxLogs NF TODO
func GetLinuxLogs() ([]string, []string) {
	var actual = []string{}
	var missing = []string{}
	// if runtime.GOOS == "linux" {
	// 	var getThese = []string{
	// 		"/var/log/syslog",
	// 		"",
	// 	}
	// 	wpath := "C:\\Windows\\System32\\winevt\\Logs\\%s"
	// 	for _, val := range getThese {
	// 		fpath := fmt.Sprintf(wpath, val)
	// 		if fs.Exists(fpath) {
	// 			actual = append(actual, fpath)
	// 		} else {
	// 			missing = append(missing, fpath)
	// 		}
	// 	}
	// }
	return actual, missing
}
