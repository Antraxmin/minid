package main

import (
	"fmt"
	"os"

	"github.com/Antraxmin/minid/pkg/container"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "minid",
		Usage: "A lightweight container runtime",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Create and run a container",
				Action: func(ctx *cli.Context) error {
					config := &container.Config{
						Name:    ctx.String("name"),
						Command: ctx.Args().Slice(),
					}
					return container.Run(config)
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Usage: "Container name",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
