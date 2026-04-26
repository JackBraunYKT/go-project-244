package main

import (
	"code"
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	// var Format formatters.Format = formatters.Stylish
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
			if cmd.NArg() < 2 {
				return fmt.Errorf("ошибка: не указан путь к файлу или директории")
			}

			args := cmd.Args()

			filepath1 := args.Get(0)
			filepath2 := args.Get(1)
			format := cmd.String("format")

			result, err := code.GenDiff(filepath1, filepath2, format)
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
