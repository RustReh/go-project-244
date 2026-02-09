package main

import (
    "fmt"
    "log"
    "os"
    "code"
    "github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name:  "gendiff",
        Usage: "Compares two configuration files and shows a difference.",
        Commands: []*cli.Command{
            {
                Name:    "compare",
                Aliases: []string{"cmp"},
                Usage:   "Compare two files",
                Flags: []cli.Flag{
                    &cli.StringFlag{
                        Name:    "format",
                        Aliases: []string{"f"},
                        Value:   "stylish",
                        Usage:   "output format",
                    },
                },
                Action: func(c *cli.Context) error {
                    if c.NArg() != 2 {
                        return cli.ShowSubcommandHelp(c)
                    }
                    path1 := c.Args().Get(0)
                    path2 := c.Args().Get(1)
                    format := c.String("format")

                    result, err := code.GenDiff(path1, path2, format)
                    if err != nil {
                        return err
                    }
                    fmt.Println(result)
                    return nil
                },
            },
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}
