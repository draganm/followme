package main

import "github.com/urfave/cli/v2"

func main() {
	app := &cli.App{
		Action: func(c *cli.Context) error {
			return nils
		},
	}
	app.RunAndExitOnError()
}
