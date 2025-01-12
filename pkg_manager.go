package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetPackageManager(dir string) (PackageManager, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("Could not read entries in directory %s\n%s\n", dir, err.Error())
	}

	if _, hasPackageLock := Find(entries, EntryIsNpmLock); hasPackageLock {
		return PackageManagerNpm, nil
	}

	if _, hasPnpmLock := Find(entries, EntryIsPnpmLock); hasPnpmLock {
		return PackageManagerPnpm, nil
	}

	if lockfile, hasYarnLock := Find(entries, EntryIsYarnLock); hasYarnLock {
		pathToLockfile := filepath.Join(dir, lockfile.Name())
		metaContent, err := ReadFirst10Lines(pathToLockfile)
		if err != nil {
			return "", fmt.Errorf("Could not determine if using yarn classic or new yarn\n%s\n", err.Error())
		}

		if LockfileIsYarnClassic(metaContent) {
			return PackageManagerYarnClassic, nil
		} else {
			return PackageManagerYarn, nil
		}

	}

	// No lockfile present, default to npm
	return PackageManagerNpm, nil
}

var YarnLockfileFilenames = []string{"yarn.lock"}
var EntryIsYarnLock = IsOneOfFilename(YarnLockfileFilenames)

const YarnClassicSignature = "yarn lockfile v1"

func LockfileIsYarnClassic(contents string) bool {
	for _, line := range strings.Split(contents, "\n") {
		lineIncludesV1Signature := strings.Contains(strings.ToLower(line), YarnClassicSignature)
		if lineIncludesV1Signature {
			return true
		}
	}

	return false
}

var PnpmLockfileFilenames = []string{"pnpm-lock.yaml", "pnpm-lock.yml"}
var EntryIsPnpmLock = IsOneOfFilename(PnpmLockfileFilenames)

var NpmLockfileFilenames = []string{"package-lock.json"}
var EntryIsNpmLock = IsOneOfFilename(NpmLockfileFilenames)

func IsOneOfFilename(checkFilenames []string) func(os.DirEntry) bool {
	return func(entry os.DirEntry) bool {
		entryFilename := entry.Name()

		_, found := Find(
			checkFilenames,
			func(checkFilename string) bool {
				return checkFilename == entryFilename
			},
		)
		return found
	}
}

func Find[T any](arr []T, checkCond func(item T) bool) (T, bool) {
	condPassed := false

	for _, item := range arr {
		if condPassed = checkCond(item); condPassed {
			return item, condPassed
		}
	}

	var zero T
	return zero, false
}

func ReadFirst10Lines(filename string) (contents string, err error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		err = fmt.Errorf("Error opening file: %s", err.Error())
		return
	}
	defer file.Close() // Ensure the file is closed when we're done

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read and print the first 10 lines
	for lineCount := 0; lineCount <= 10 || scanner.Scan(); lineCount += 1 {
		contents += scanner.Text()
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("Error reading file: %s", err.Error())
	}

	return contents, nil
}
