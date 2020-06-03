package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Author = "yeqown@gmail.com"
	app.Copyright = "2020@yeqown"

	mountCommands(app)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
