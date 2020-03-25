package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fefit/fetc/commands"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "fetc"
	app.Usage = "fetc is the fet's command line tool"
	app.Commands = []*cli.Command{
		commands.Init(),
		commands.Watch(),
	}
	app.Action = func(c *cli.Context) error {
		fmt.Println("thank you for using fet template engineer, use 'fetc -h' for helps.")
		return nil
	}
	app.Version = "0.0.1"
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
