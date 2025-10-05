package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	code "github.com/darkartx/go-project-244"

	cli "github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Usage: "Compares two configuration files and shows a difference.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Value:   "stylish",
				Usage:   "output format",
				Aliases: []string{"f"},
			},
		},
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "filepath1",
			},
			&cli.StringArg{
				Name: "filepath2",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			filepath1 := cmd.StringArg("filepath1")

			if len(filepath1) == 0 {
				return errors.New("filepath1 requred")
			}

			filepath2 := cmd.StringArg("filepath2")

			if len(filepath2) == 0 {
				return errors.New("filepath2 requred")
			}

			format := cmd.String("format")

			diff, err := code.GenDiff(filepath1, filepath2, format)

			if err != nil {
				return err
			}

			fmt.Println(diff)

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
