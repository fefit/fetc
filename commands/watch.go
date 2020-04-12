package commands

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	pfet "github.com/fefit/fet"
	"github.com/fefit/fet/types"
	"github.com/fefit/fetc/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli/v2"
)

// run the command
func runWatch() error {
	var (
		conf    *types.FetConfig
		err     error
		watcher *fsnotify.Watcher
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
		err = filepath.Walk(fet.TemplateDir, func(curPath string, fi os.FileInfo, err error) error {
			if fi.Mode().IsDir() && !utils.IsSpecialDorf(curPath) && !fet.NeedIgnore(curPath) {
				return watcher.Add(curPath)
			}
			return nil
		})
		if err != nil {
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
					tpl := fet.RealTmplPath(name)
					if utils.IsSpecialDorf(tpl) {
						break
					}
					if op == fsnotify.Chmod {
						// ignore
					} else if op == fsnotify.Remove {
						ctpl := fet.RealCmplPath(tpl)
						if fet.NeedIgnore(ctpl) {
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
										if utils.ContainsStr(deps, tpl) {
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
						if isNeedAddSelf && !fet.NeedIgnore(tpl) {
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
			return runWatch()
		},
	}
}
