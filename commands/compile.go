package commands

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	pfet "github.com/fefit/fet"
	"github.com/fefit/fet/types"
	"github.com/urfave/cli/v2"
)

// run the command
func runCompile(c *cli.Context) error {
	var (
		conf *types.FetConfig
		err  error
	)
	// make sure config is ok
	if conf, err = pfet.LoadConf("fet.config.json"); err != nil {
		return err
	}
	// always disable debug
	conf.Debug = false
	// compile files
	if fet, mErr := pfet.New(conf); mErr == nil {
		// need compile some files
		fileList := c.Args()
		totalNum := fileList.Len()
		allFiles := []string{}
		isCompileAll := totalNum == 0
		if totalNum > 0 {
			for i := 0; i < totalNum; i++ {
				var files []string
				cur := fileList.Get(i)
				if files, err = fet.GetCompileFiles(fet.RealTmplPath(cur)); err != nil {
					return err
				}
				allFiles = append(allFiles, files...)
			}
		} else {
			/* compile all files */
			if allFiles, err = fet.GetCompileFiles(fet.TemplateDir); err != nil {
				return err
			}
		}
		colorful := color.New(color.FgYellow).PrintfFunc()
		logLn := func(info string, args ...interface{}) {
			colorful(info, args...)
			fmt.Println()
		}
		totalCount := len(allFiles)
		startTime := time.Now()
		if totalCount > 0 {
			logLn("Total Files: %d", totalCount)
			if isCompileAll {
				_, err = fet.CompileAll()
			} else {
				for _, file := range allFiles {
					if _, _, err = fet.Compile(file, true); err != nil {
						return err
					}
				}
			}
			if err == nil {
				endTime := time.Now()
				logLn("All Files Compiled success: %v", endTime.Sub(startTime))
			}
		} else {
			logLn("No files need compile")
		}
	} else {
		err = mErr
	}
	return err
}

// Compile command
func Compile() *cli.Command {
	return &cli.Command{
		Name:    "compile",
		Aliases: []string{"c"},
		Usage:   "compile the files",
		Action: func(c *cli.Context) error {
			return runCompile(c)
		},
	}
}
