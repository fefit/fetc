package commands

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	pfet "github.com/fefit/fet"
	"github.com/fefit/fet/types"
	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli/v2"
)

// check director of file if exists
func checkDorfExists(pathname string) (notexist bool, err error) {
	if _, err := os.Stat(pathname); err != nil {
		if os.IsNotExist(err) {
			return true, fmt.Errorf("the file or directory '%s' is not exist", pathname)
		}
		return false, err
	}
	return false, nil
}

// e.g .vscode .a.html.swap
func isSpecialDorf(dorf string) bool {
	return strings.HasPrefix(path.Base(dorf), ".")
}

var watcher *fsnotify.Watcher

func watchDir(curPath string, fi os.FileInfo, err error) error {
	if fi.Mode().IsDir() && !isSpecialDorf(curPath) {
		return watcher.Add(curPath)
	}
	return nil
}

func contains(arr []string, key string) bool {
	for _, cur := range arr {
		if cur == key {
			return true
		}
	}
	return false
}

// run the command
func run() error {
	var (
		conf *types.FetConfig
		err  error
	)
	if conf, err = pfet.LoadConf("fet.config.json"); err != nil {
		return err
	}
	if fet, mErr := pfet.New(conf); mErr == nil {
		// first, compile all files, get the includes and extends map
		fileDeps, err := fet.CompileAll()
		if err != nil {
			fmt.Println("compile error:", err.Error())
		}
		// create watcher
		watcher, _ = fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()
		if err := filepath.Walk(fet.TemplateDir, watchDir); err != nil {
			fmt.Println("watch error:", err)
		}
		done := make(chan bool)
		//
		go func() {
			for {
				select {
				// watch for events
				case event := <-watcher.Events:
					name, op := event.Name, event.Op
					tpl := fet.GetTemplateFile(name)
					if isSpecialDorf(tpl) {
						break
					}
					if op == fsnotify.Chmod {
						// ignore
					} else if op == fsnotify.Remove {
						ctpl := fet.GetCompileFile(tpl)
						if fet.IsIgnoreFile(ctpl) {
							// do nothing
						} else {
							// delete the compile file
							fmt.Println("delete compiled file:", ctpl)
							os.Remove(ctpl)
						}
						fileDeps.Delete(ctpl)
					} else {
						files := []string{}
						isNeedAddSelf := true
						if op == fsnotify.Create {
							// add self file
						} else {
							fmt.Println("changes:", tpl)
							fileDeps.Range(func(key, value interface{}) bool {
								if curTpl, ok := key.(string); ok {
									if deps, ok := value.([]string); ok {
										if contains(deps, tpl) {
											if curTpl == tpl {
												isNeedAddSelf = false
											}
											files = append(files, curTpl)
										}
									}
								}
								return true
							})
						}
						if isNeedAddSelf && !fet.IsIgnoreFile(tpl) {
							files = append(files, tpl)
						}
						var wg sync.WaitGroup
						wg.Add(len(files))
						for _, curTpl := range files {
							go func(tpl string, conf *types.FetConfig) {
								fet, _ = pfet.New(conf)
								_, deps, err := fet.Compile(tpl, true)
								if err != nil {
									fmt.Println("compile failure:", err.Error())
								} else {
									fileDeps.Store(tpl, deps)
								}
								wg.Done()
							}(curTpl, conf)
						}
						wg.Wait()
					}
					// watch for errors
				case err := <-watcher.Errors:
					fmt.Println("watch error:", err)
				}
			}
		}()
		<-done
	} else {
		err = mErr
	}
	return err
}

// Watch command
func Watch() *cli.Command {
	return &cli.Command{
		Name:    "watch",
		Aliases: []string{"w"},
		Usage:   "watch the file fet template files changes and compile them",
		Action: func(c *cli.Context) error {
			return run()
		},
	}
}
