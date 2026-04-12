package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name: "gendiff",

		Usage: "Compares two configuration files and shows a difference",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "Format",
				Aliases:     []string{"f"},
				Usage:       "output format",
				DefaultText: "stylish",
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
		os.Exit(1)
	}
}
