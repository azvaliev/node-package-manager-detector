package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type PackageManager string

const (
	PackageManagerNpm         PackageManager = "npm"
	PackageManagerYarn        PackageManager = "yarn"
	PackageManagerYarnClassic PackageManager = "yarn-classic"
	PackageManagerPnpm        PackageManager = "pnpm"
)

func main() {
	dir, err := getDir()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}

	packageManager, err := GetPackageManager(dir)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(packageManager)
	os.Exit(0)
}

func getDir() (dir string, err error) {
	flag.Func(
		"dir",
		"(optional) absolute path to directory to which to introspect. Default to cwd",
		func(suppliedDir string) error {
			if suppliedDir == "" {
				return nil
			}

			if !filepath.IsAbs(suppliedDir) {
				return fmt.Errorf("%s is not an absolute path", suppliedDir)
			}

			return nil
		},
	)
	flag.BoolFunc("help", "Print help message", func(_ string) error {
		flag.Usage()
		return nil
	})
	flag.Parse()

	if dir == "" {
		dir, err = os.Getwd()

		if err != nil {
			return dir, errors.Join(
				errors.New("Failed to determine current working directory. Please specify"),
				err,
			)
		}
	}

	return dir, nil
}
