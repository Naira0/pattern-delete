package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	recursive bool
	pattern   string
	extension string
)

func init() {
	log.SetPrefix("Error: ")
	log.SetFlags(0)

	flag.BoolVar(&recursive, "r", false, "If set it will recursively delete files.")
	flag.StringVar(&extension, "e", "", "The extension to look for. If not set it will ignore extensions")
	flag.StringVar(&pattern, "p", "", "The pattern to delete for. If not set it will either default to the given extension flag or it wont do anything at all.")

	flag.Parse()
}

func handleError(msg string, err error) {
	if err != nil {
		log.Fatal(msg, err.Error())
	}
}

func hasPattern(file string) bool {

	if pattern == "" {
		return false
	}

	match, err := regexp.MatchString(pattern, file[:strings.Index(file, ".")])

	handleError("Could not validate the given regexp", err)

	return match
}

func hasExtension(file string) bool {

	if strings.HasSuffix(file, "."+extension) && extension != "" {
		return true
	}
	return false
}

func traverseRecursive(cwd string) {

	var deleted int
	var couldNotDelete int

	walk := filepath.Walk(cwd, func(p string, f os.FileInfo, err error) error {

		if f.IsDir() {
			return nil
		}

		var delError error

		if hasExtension(f.Name()) && hasPattern(f.Name()) {
			delError = os.Remove(p)
		} else if (pattern != "" && extension != "") && (!hasPattern(f.Name()) || !hasExtension(f.Name())) {
			return nil
		} else if hasExtension(f.Name()) || hasPattern(f.Name()) {
			delError = os.Remove(p)
		} else {
			return nil
		}

		if delError != nil {
			couldNotDelete++
			return delError
		}

		deleted++

		return err
	})

	handleError("could not find files", walk)

	fmt.Println("Deleted", deleted, "files")

	if couldNotDelete > 0 {
		fmt.Println("Failed to delete", couldNotDelete, "Files")
	}
}

func traverseFolder(cwd string) {

	dir, err := ioutil.ReadDir(cwd)

	handleError("Could not read directory", err)

	var deleted int
	var couldNotDelete int

	for _, file := range dir {

		if file.IsDir() {
			continue
		}

		var err error

		if hasExtension(file.Name()) && hasPattern(file.Name()) {
			err = os.Remove(path.Join(cwd, file.Name()))
		} else if (pattern != "" && extension != "") && (!hasPattern(file.Name()) || !hasExtension(file.Name())) {
			continue
		} else if hasExtension(file.Name()) || hasPattern(file.Name()) {
			err = os.Remove(path.Join(cwd, file.Name()))
		} else {
			continue
		}

		if err != nil {
			couldNotDelete++
			continue
		}

		deleted++
	}

	fmt.Println("Deleted", deleted, "files")

	if couldNotDelete > 0 {
		fmt.Println("Failed to delete", couldNotDelete, "Files")
	}
}

func main() {

	if extension == "" && pattern == "" {
		log.Fatal("You must provide either -p or -e flags. use -help for more info.")
	}

	cwd, err := os.Getwd()

	handleError("could not get the current working directory", err)

	if recursive {
		traverseRecursive(cwd)
	} else {
		traverseFolder(cwd)
	}
}
