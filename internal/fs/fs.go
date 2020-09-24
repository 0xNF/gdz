/// Package fs contains library functions for getting stuff from the host System

package fs

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const winPath38 string = "${ProgramData}\\Gravio"
const winPath40 string = "${ProgramData}\\Hubkit"

const androidPath38 string = ""
const androidPath40 string = ""

const dockerPath38 string = ""
const dockerPath40 string = ""

const linPath38 string = ""
const linPath40 string = ""

const macPath38 string = ""
const macPath40 string = ""

var os2PathDict = map[string]map[string]string{
	"windows": {"3.8": winPath38, "4.0": winPath40},
	"darwin":  {"3.8": macPath38, "4.0": macPath40},
	"linux":   {"3.8": linPath38, "4.0": linPath40},
	"docker":  {"3.8": dockerPath38, "4.0": dockerPath40},
	"android": {"3.8": androidPath38, "4.0": androidPath40},
}

func split(path string) string {
	var fullPath string = os.ExpandEnv(path)
	baseDir := filepath.Base(fullPath)
	return baseDir
}

func Get(c *Conf) []string {
	paths := []string{}
	for _, v := range c.Versions {
		rt := runtime.GOOS
		path := os2PathDict[rt][v]
		p, err := getStuff(path, path, fmt.Sprintf("%s_%s", rt, v), c.Verbose)
		if err != nil {
			panic(err)
		}
		paths = append(paths, p)
	}
	return paths
}

func getStuff(fpath string, bpath string, which string, verbose bool) (string, error) {
	var fullPath string = os.ExpandEnv(fpath)
	if exists(fullPath) {
		zpath := makeZip(fullPath, bpath, which, verbose)
		fmt.Printf("Finished creating Gravio Diagnostics Zip File. This file can be found at %s\n", zpath)
		return zpath, nil
	}
	return "", fmt.Errorf("Couldn't find Gravio data at %s", winPath38)
}

func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path/to/whatever does not exist
		return false
	}
	return true
}

// Zips "./input" into "./output.zip"
func makeZip(fpath string, bpath string, which string, verbose bool) string {
	t := time.Now() //It will return time.Time object with current timestamp
	tUnixMilli := int64(time.Nanosecond) * t.UnixNano() / int64(time.Millisecond)
	zipName := fmt.Sprintf("GravioDiagnostics_%s_%d.zip", which, tUnixMilli)
	file, err := os.Create(zipName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		if verbose {
			// fmt.Printf("Crawling: %#v\n", path)
		}
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Ensure that `path` is not absolute; it should not start with "/".
		// This snippet happens to work because I don't use
		// absolute paths, but ensure your real-world code
		// transforms path into a zip-root relative path.
		mpath := strings.Replace(path, os.ExpandEnv(bpath), split(bpath), 1)
		fmt.Println("fpath: " + path)
		fmt.Println("bpath: " + os.ExpandEnv(bpath))
		fmt.Println("rpl: " + mpath)
		f, err := w.Create(mpath)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	err = filepath.Walk(fpath, walker)
	if err != nil {
		panic(err)
	}
	return zipName
}
