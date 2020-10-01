package collect

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/0xNF/gdz/internal/events"
	"github.com/0xNF/gdz/internal/fs"
	"github.com/0xNF/gdz/internal/gdzerr"
)

/* Studio Paths */
const stWinPath38 string = "${UserProfile}\\AppData\\Local\\Packages\\InfoteriaPte.Ltd.GravioStudio_mrnz526z5qc9p"
const stWinPath40 string = "${UserProfile}\\AppData\\Local\\Packages\\InfoteriaPte.Ltd.GravioStudio4_mrnz526z5qc9p"

const stAndroidPath38 string = ""
const stAndroidPath40 string = ""

const stDockerPath38 string = ""
const stDockerPath40 string = ""

const stLinPath38 string = ""
const stLinPath40 string = ""

const stMacPath38 string = ""
const stMacPath40 string = ""

/* Hubkit Paths */
const hkWinPath38 string = "${ProgramData}\\Gravio"
const hkWinPath40 string = "${ProgramData}\\Hubkit"

const hkAndroidPath38 string = ""
const hkAndroidPath40 string = ""

const hkDockerPath38 string = ""
const hkDockerPath40 string = ""

const hkLinPath38 string = ""
const hkLinPath40 string = ""

const hkMacPath38 string = ""
const hkMacPath40 string = ""

/* Coordinator paths */
const coWinPath38 string = "${ProgramData}\\Gravio"
const coWinPath40 string = "${ProgramData}\\Hubkit"

const coAndroidPath38 string = ""
const coAndroidPath40 string = ""

const coDockerPath38 string = ""
const coDockerPath40 string = ""

const coLinPath38 string = ""
const coLinPath40 string = ""

const coMacPath38 string = ""
const coMacPath40 string = ""

var os2HubkitPathDict = map[string]map[string]string{
	"windows": {"3.8": hkWinPath38, "4.0": hkWinPath40},
	"darwin":  {"3.8": hkMacPath38, "4.0": hkMacPath40},
	"linux":   {"3.8": hkLinPath38, "4.0": hkLinPath40},
	"docker":  {"3.8": hkDockerPath38, "4.0": hkDockerPath40},
	"android": {"3.8": hkAndroidPath38, "4.0": hkAndroidPath40},
}

var os2StudioPathDict = map[string]map[string]string{
	"windows": {"3.8": stWinPath38, "4.0": stWinPath40},
	"darwin":  {"3.8": stMacPath38, "4.0": stMacPath40},
	"linux":   {"3.8": stLinPath38, "4.0": stLinPath40},
	"docker":  {"3.8": stDockerPath38, "4.0": stDockerPath40},
	"android": {"3.8": stAndroidPath38, "4.0": stAndroidPath40},
}

var os2LogsDict = map[string]map[string]string{
	"windows": {"3.8": coWinPath38, "4.0": coWinPath40},
	"darwin":  {"3.8": coMacPath38, "4.0": coMacPath40},
	"linux":   {"3.8": coLinPath38, "4.0": coLinPath40},
	"docker":  {"3.8": coDockerPath38, "4.0": coDockerPath40},
	"android": {"3.8": coAndroidPath38, "4.0": coAndroidPath40},
}

var timestamp string = ""

func Get(c *fs.Conf) []string {
	t := time.Now() //It will return time.Time object with current timestamp
	timestamp = fmt.Sprintf("%d", int64(time.Nanosecond)*t.UnixNano()/int64(time.Millisecond))
	gdzerr.Timestamp = timestamp
	var e error
	if e != nil {
		panic(e)
	}
	defer gdzerr.CloseError()

	MakeZip2()

	paths := []string{}
	return paths
}

func grabSystemLogs() map[string][]string {
	var d = make(map[string][]string)
	here, _, err := events.GetEventLogs()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", here)
	d["EventLogs"] = here
	return d
}

// Gets paths to available hubkits on this device
// returns Version:path, e.g,
// { 3.8: %programdata/hubkit, 4.0: %programdata%/gravio}
func grabHubkits() map[string]string {
	var d = make(map[string]string)
	rt := runtime.GOOS
	for k, v := range os2HubkitPathDict[rt] {
		expanded := fs.Expand(v)
		if fs.Exists(expanded) {
			d[k] = expanded
		}
	}
	return d
}

// Gets paths to available Studios on this device
func grabStudios() map[string]string {
	var d = make(map[string]string)
	rt := runtime.GOOS
	for k, v := range os2StudioPathDict[rt] {
		expanded := fs.Expand(v)
		if fs.Exists(expanded) {
			d[k] = expanded
		}
	}
	return d
}

// CollectPaths gets all the paths for zipping on the client machine
func collectPaths() map[string]map[string][]string {
	pathmap := map[string]map[string][]string{} /* dict of BaseFolder:[version]:[paths] on this machine to hoover up. ex: {Hubkits: 4.0: [program/hubkit], 3.8: [program/gravio]} */

	/* Grab Hubkit data on the machine */
	var hkKey = "Hubkits"
	for k, v := range grabHubkits() {
		if _, ok := pathmap[hkKey]; !ok {
			pathmap[hkKey] = map[string][]string{}
			pathmap[hkKey][k] = []string{v}
		} else {
			pathmap[hkKey][k] = append(pathmap[hkKey][k], v)
		}
	}

	/* TODO Grab Coordinator data if available */

	/* Grab Studio Data if available */
	var sKey = "Studios"
	for k, v := range grabStudios() {
		if _, ok := pathmap[sKey]; !ok {
			pathmap[sKey] = map[string][]string{}
			pathmap[sKey][k] = []string{v}
		} else {
			pathmap[sKey][k] = append(pathmap[sKey][k], v)
		}
	}

	/* TODO Grab System Logs */
	var lKey = "SystemLogs"
	for k, v := range grabSystemLogs() {
		for _, v2 := range v {
			if _, ok := pathmap[lKey]; !ok {
				pathmap[lKey] = map[string][]string{}
			}
			if _, ok := pathmap[lKey][k]; !ok {
				pathmap[lKey][k] = []string{v2}
			} else {
				pathmap[lKey][k] = append(pathmap[lKey][k], v2)
			}
		}
	}

	return pathmap

}

// MakeZip2 makes a new zip file with all the relevant content
func MakeZip2() string {
	/* make a zip with the current timestamp */
	zipName := fmt.Sprintf("GravioDiagnostics_%s.zip", timestamp)
	file, err := os.Create(zipName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	fpaths := collectPaths()
	fs.PrettyPrint(fpaths)

	for bdir, lst := range fpaths {
		for cdir, items := range lst {
			for _, item := range items {
				addToZip(item, bdir, cdir, w)
			}
		}
	}

	return zipName
}

var skipThese []string = []string{
	filepath.Join("_mrnz526z5qc9p", "AC"), /* too much useless cryto stuff from Microsoft, doesn't matter */
}

func addToZip(basePath string, baseFolder string, childFolder string, zipFile *zip.Writer) {

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		expanded := fs.Expand(path)
		/* skip items in the blacklist */
		for _, skip := range skipThese {
			if strings.Contains(expanded, skip) {
				return nil
			}
		}
		file, err := os.Open(path)
		if err != nil {
			gdzerr.WriteErr(err.Error())
			return err
		}
		defer file.Close()

		// Ensure that `path` is not absolute; it should not start with "/".
		// This snippet happens to work because I don't use
		// absolute paths, but ensure your real-world code
		// transforms path into a zip-root relative path.
		expandedBase := os.ExpandEnv(basePath)
		mainDir := strings.Replace(expandedBase, expandedBase, fs.Split(basePath), 1)
		mmpath := strings.Replace(path, expandedBase, mainDir, 1)
		mpath := filepath.Join(baseFolder, childFolder, mmpath)
		f, err := zipFile.Create(mpath)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}

	if !fs.Exists(basePath) {
		return
	}
	err := filepath.Walk(basePath, walker)
	if err != nil {
		panic(err)
	}

}
