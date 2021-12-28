package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/radovskyb/watcher"
)

var gameInstance *exec.Cmd

func main() {
	_, file, _, _ := runtime.Caller(0)
	dir := path.Dir(file)
	fmt.Println(dir)

	compileAndStartProcess(false)
	watchForHotReloadRequested()
}

func compileAndStartProcess(reloadPreviousState bool) {
	cmnd := exec.Command("go", "build", "-o", "game.exe", "main.go")
	cmnd.Run()
	
	if  reloadPreviousState {
		gameInstance = exec.Command("game.exe", "-hotpid", strconv.Itoa(os.Getpid()), "-hot")
	} else {
		gameInstance = exec.Command("game.exe", "-hotpid", strconv.Itoa(os.Getpid()))
	}

	gameInstance.Stderr = os.Stderr
	gameInstance.Stdout = os.Stdout
	gameInstance.Start()
}

func watchForHotReloadRequested() {
	w := watcher.New()

	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write)

	_, file, _, _ := runtime.Caller(0)
	dir := path.Dir(file)

	filename := path.Join(dir, "../hotstate.json")
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		handle, _ := os.Create(filename)
		if handle != nil {
			handle.Close()
		}
	}

	_ = w.Add(filename)

	go func() {
		for {
			select {
			case <-w.Event:	
				gameInstance.Process.Kill()
				compileAndStartProcess(true)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}