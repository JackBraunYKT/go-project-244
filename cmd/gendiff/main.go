package main

import (
	"code"
	"code/formatters"
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	// var Format formatters.Format = formatters.Stylish
	Format := formatters.JSON

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

		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.NArg() < 1 {
				return fmt.Errorf("ошибка: не указан путь к файлу или директории")
			}

			filepath1 := cmd.Args().First()
			filepath2 := cmd.Args().Get(1)

			result, err := code.GenDiff(filepath1, filepath2, Format)
			if err != nil {
				return err
			}

			fmt.Printf("%s", *result)

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "ошибка: %v\n", err)
		os.Exit(1)
	}
}
