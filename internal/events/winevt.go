/// Package winevt deals with getting Windows Event Log information

package events

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/0xNF/gdz/internal/fs"
)

// GetEventLogs returns paths of event logs, which can be added to thew big zip file
// returns (valid paths, missing files)
func GetEventLogs() ([]string, []string, error) {
	if runtime.GOOS == "windows" {
		return getWindowsLogs()
	} else if runtime.GOOS == "darwin" {
		return getMacLogs()
	} else if runtime.GOOS == "linux" {
		return getLinuxLogs()
	}
	return nil, nil, errors.New("Invalid runtime detected")
}

func getWindowsLogs() ([]string, []string, error) {
	var actual = []string{}
	var missing = []string{}

	var getThese = []string{
		"Application.evtx",
		"System.evtx",
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

	return actual, missing, nil
}

func getLinuxLogs() ([]string, []string, error) {
	var actual = []string{}
	var missing = []string{}
	return actual, missing, nil
}

func getMacLogs() ([]string, []string, error) {
	var actual = []string{}
	var missing = []string{}
	return actual, missing, nil
}
