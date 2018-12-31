// +build mage

package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

var (
	pkgName     = "usersd"
	dockerImage = "ntrrg/" + pkgName

	goFiles    = getGoFiles()
	goSrcFiles = getGoSrcFiles()
)

var Default = Build

func Build() error {
	if run, err := target.Path("dist/"+pkgName, goSrcFiles...); !run || err != nil {
		return err
	}

	env := map[string]string{"CGO_ENABLED": "0"}
	return sh.RunWith(env, "go", "build", "-o", "dist/"+pkgName)
}

type Clean mg.Namespace

func (Clean) Default() {
	sh.Rm("dist")
	sh.Run("docker", "image", "rm", dockerImage)
}

type Docker mg.Namespace

func (Docker) Default() error {
	return sh.RunV("docker", "build", "-t", dockerImage, ".")
}

// Development

var (
	coverageFile = "coverage.txt"
)

func Benchmark() error {
	return sh.RunV("go", "test", "-race", "-bench", ".", "-benchmem", "./...")
}

func CI() {
	mg.SerialDeps(Lint, QA, Test, Coverage.Default, Benchmark, Build)
}

func (Clean) Dev() {
	mg.Deps(Clean.Default)
	sh.Rm(coverageFile)
	sh.Run("docker", "image", "rm", dockerImage+":debug")
}

type Coverage mg.Namespace

func (Coverage) Default() error {
	mg.Deps(CoverageFile)
	return sh.RunV("go", "tool", "cover", "-func", coverageFile)
}

func (Coverage) Web() error {
	mg.Deps(CoverageFile)
	return sh.RunV("go", "tool", "cover", "-html", coverageFile)
}

func CoverageFile() error {
	if run, err := target.Path(coverageFile, goFiles...); !run || err != nil {
		return err
	}

	return sh.RunV("go", "test", "-race", "-coverprofile", coverageFile, "./...")
}

func (Docker) Debug() error {
	return sh.RunV("docker", "build", "--target", "build", "-t", dockerImage+":debug", ".")
}

type Docs mg.Namespace

func (Docs) Ref() {
	sh.RunV("godoc", "-http", ":6060", "-play")
}

func (Docs) Rest() error {
	return errors.New("WIP..")
}

func Format() error {
	args := []string{"-s", "-w", "-l"}
	args = append(args, goFiles...)
	return sh.RunV("gofmt", args...)
}

func Lint() error {
	args := []string{"-d", "-e", "-s"}
	args = append(args, goFiles...)
	return sh.RunV("gofmt", args...)
}

func Install() error {
	mg.Deps(Build)
	return nil
}

func QA() error {
	return sh.RunV("golangci-lint", "run")
}

func Test() error {
	return sh.RunV("go", "test", "-race", "-v", "./...")
}

// Helpers

func getGoFiles() []string {
	var goFiles []string

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, "vendor/") {
			return filepath.SkipDir
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		goFiles = append(goFiles, path)
		return nil
	})

	return goFiles
}

func getGoSrcFiles() []string {
	var goSrcFiles []string

	for _, path := range goFiles {
		if !strings.HasSuffix(path, "_test.go") {
			continue
		}

		goSrcFiles = append(goSrcFiles, path)
	}

	return goSrcFiles
}
