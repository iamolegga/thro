package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/adrg/xdg"
	"github.com/gosimple/slug"
	"github.com/juju/fslock"
)

const appName = "thro"
const deferPrefix = "defer-"

var stateDir = path.Join(xdg.StateHome, appName)

func main() {
	prepareStateDir()

	args := os.Args[1:]

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmdStr := strings.Join(args, " ")
	lockfile := path.Join(stateDir, slug.Make(cmdStr))
	deferLockfile := path.Join(stateDir, deferPrefix+slug.Make(cmdStr))

	lock := fslock.New(lockfile)
	if err := lock.TryLock(); err == nil {
		defer os.Remove(lockfile)
		defer lock.Unlock()
		_ = cmd.Run()
		return
	}

	deferred := fslock.New(deferLockfile)
	if err := deferred.TryLock(); err == nil {
		_ = lock.Lock()
		defer os.Remove(deferLockfile)
		defer lock.Unlock()
		_ = deferred.Unlock()
		_ = cmd.Run()
	}
}

func prepareStateDir() {
	if _, err := os.Stat(stateDir); os.IsNotExist(err) {
		if err := os.MkdirAll(stateDir, os.ModePerm); err != nil {
			log.Fatal(fmt.Errorf("thro: unable to create state directory %s: %w", stateDir, err))
		}
	} else if err != nil {
		log.Fatalf("thro: unable to read state directory: %s", stateDir)
	}
}
