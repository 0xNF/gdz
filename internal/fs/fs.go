/// Package fs contains library functions for getting stuff from the host System

package fs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// PrettyPrint prints golang structures to console in a user friendly way.
// Why isnt this a built in feature?
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

// Split takes a filepath and returns the Base Directory - i.e., C:/Windows/System32/run32.dll returns "System32", C:/Users/jsmith/Documents/ returns "Documents"
func Split(path string) string {
	var fullPath string = os.ExpandEnv(path)
	baseDir := filepath.Base(fullPath)
	return baseDir
}

// Expand expands a path to its full path on disk
func Expand(path string) string {
	return os.ExpandEnv(path)
}

// Exists tests wether the path exists
func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path/to/whatever does not exist
		return false
	}
	return true
}

// Only call this after determinig the path exists
func isDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// // Zips "./input" into "./output.zip"
// func makeZip(fpath string, bpath string, which string, verbose bool) string {
// 	grabHubkits()

// 	t := time.Now() //It will return time.Time object with current timestamp
// 	tUnixMilli := int64(time.Nanosecond) * t.UnixNano() / int64(time.Millisecond)
// 	zipName := fmt.Sprintf("GravioDiagnostics_%s_%d.zip", which, tUnixMilli)
// 	file, err := os.Create(zipName)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()

// 	w := zip.NewWriter(file)
// 	defer w.Close()

// 	walker := func(path string, info os.FileInfo, err error) error {
// 		if verbose {
// 			// fmt.Printf("Crawling: %#v\n", path)
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		if info.IsDir() {
// 			return nil
// 		}
// 		file, err := os.Open(path)
// 		if err != nil {
// 			return err
// 		}
// 		defer file.Close()

// 		// Ensure that `path` is not absolute; it should not start with "/".
// 		// This snippet happens to work because I don't use
// 		// absolute paths, but ensure your real-world code
// 		// transforms path into a zip-root relative path.
// 		mpath := strings.Replace(path, os.ExpandEnv(bpath), split(bpath), 1)
// 		fmt.Println("fpath: " + path)
// 		fmt.Println("bpath: " + os.ExpandEnv(bpath))
// 		fmt.Println("rpl: " + mpath)
// 		f, err := w.Create(mpath)
// 		if err != nil {
// 			return err
// 		}

// 		_, err = io.Copy(f, file)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	}
// 	err = filepath.Walk(fpath, walker)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return zipName
// }

// MakeZip2 creates a zip file next to the exe and fills it with the data from the files supplied by the Paths.
// For entries in Paths, if the entry is a Folder, the whole folder is taken
// If the entry is a File, then only that file is taken
// Returns the

// Adds {path} to {zip}, underneath the {basefolder} in the zip file.
// e.g., (%programdata%/gravio, hubkit38, ~/gdz.zip) adds %programdata%/gravio to gdz.zip as Hubkit38/[..data...]
