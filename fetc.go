package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/fefit/fetc/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	var (
		fetVersion []byte
		err        error
	)
	fetVersion, err = exec.Command("bash", "-c", "go list -m -u github.com/fefit/fet|awk '{print $2}'").Output()
	usage := "'fetc' is the command line tool of 'fet' template engine"
	if err == nil {
		usage += ", cur version of 'fet' is " + string(fetVersion)
	}
	// create the command line
	app := cli.NewApp()
	app.Name = "fetc"
	app.Usage = usage
	app.Commands = []*cli.Command{
		commands.Init(),
		commands.Watch(),
		commands.Compile(),
	}
	app.Action = func(c *cli.Context) error {
		fmt.Println("welcome to use 'fet' template engine, use 'fetc -h' for helps.")
		return nil
	}
	app.Version = "0.1.4"
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
