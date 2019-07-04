package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/AlecAivazis/survey"
	"github.com/fefit/fet/types"
	"github.com/urfave/cli"
)

// validate if empty
func validIfEmpty(name string) func(val interface{}) error {
	return func(val interface{}) error {
		if value, ok := val.(string); !ok || strings.TrimSpace(value) == "" {
			return fmt.Errorf("'%s' must be a string not empty", name)
		}
		return nil
	}
}

// config the fet
func configForm() {
	qs := []*survey.Question{
		{
			Name: "mode",
			Prompt: &survey.Select{
				Message: "please choose the compile mode of fet:",
				Options: []string{"Smarty", "Gofet"},
				Default: "Smarty",
			},
		},
		{
			Name: "leftDelimiter",
			Prompt: &survey.Input{
				Message: "please set the 'leftDelimiter' of fet:",
				Default: "{%",
			},
			Validate: validIfEmpty("leftDelimiter"),
		},
		{
			Name: "rightDelimiter",
			Prompt: &survey.Input{
				Message: "please set the 'rightDelimiter' of fet:",
				Default: "%}",
			},
			Validate: validIfEmpty("rightDelimiter"),
		},
		{
			Name: "ucaseField",
			Prompt: &survey.Confirm{
				Message: "do you need to set the first character of fields(struct/map etc.) to uppercase?",
			},
		},
		{
			Name: "templateDir",
			Prompt: &survey.Input{
				Message: "the directory of your fet template files:",
				Default: "templates",
			},
			Validate: survey.Required,
		},
		{
			Name: "compileDir",
			Prompt: &survey.Input{
				Message: "the directory of your fet compiled files:",
				Default: "views",
			},
			Validate: survey.Required,
		},
		{
			Name: "ignores",
			Prompt: &survey.Input{
				Message: "the fet files or directories do not need compile, use ',' split them(use golang glob):",
				Default: "inc/*",
			},
		},
	}
	answers := struct {
		Mode           string
		LeftDelimiter  string
		RightDelimiter string
		UcaseField     bool
		TemplateDir    string
		CompileDir     string
		Ignores        string
	}{}
	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	config := types.FetConfig{}
	if answers.Mode == "Smarty" {
		config.Mode = types.Smarty
	} else {
		config.Mode = types.Gofet
	}
	config.LeftDelimiter = answers.LeftDelimiter
	config.RightDelimiter = answers.RightDelimiter
	config.TemplateDir = answers.TemplateDir
	config.CompileDir = answers.CompileDir
	if answers.Ignores != "" {
		dorf := strings.Split(answers.Ignores, ",")
		ignores := []string{}
		for _, cur := range dorf {
			name := strings.TrimSpace(cur)
			if name != "" {
				ignores = append(ignores, name)
			}
		}
		config.Ignores = ignores
	}
	confdata, _ := json.Marshal(config)
	err = ioutil.WriteFile("fet.config.json", confdata, 0644)
	if err != nil {
		fmt.Println("Sorry,write the config file error:" + err.Error())
	} else {
		fmt.Printf("Your 'fet' config file was created successfully!")
	}
}

// Init command
func Init() cli.Command {
	return cli.Command{
		Name:    "init",
		Aliases: []string{"i"},
		Usage:   "initialize the fet's configs",
		Action: func(c *cli.Context) error {
			if len(c.Args()) > 0 {
				fmt.Println("the command 'init' do not receive any argument")
			} else {
				configForm()
			}
			return nil
		},
	}
}
