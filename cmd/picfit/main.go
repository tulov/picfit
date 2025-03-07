package main

import (
	"fmt"
	"os"

	"github.com/thoas/picfit/constants"
	"github.com/thoas/picfit/server"
	"github.com/thoas/picfit/signature"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "picfit"
	app.Version = fmt.Sprintf("%s [git:%s:%s]\ncompiled using %s at %s", constants.Version,
		constants.Branch, constants.Revision, constants.Compiler, constants.BuildTime)
	app.Author = "thoas"
	app.Email = "florent.messa@gmail.com"
	app.Usage = "Display, manipulate, transform and cache your images"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Config file path",
			EnvVar: "PICFIT_CONFIG_PATH",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "version",
			ShortName: "v",
			Usage:     "Retrieve the version number",
			Action: func(c *cli.Context) {
				fmt.Printf("picfit %s\n", constants.Version)
			},
		},
		{
			Name:      "signature",
			ShortName: "s",
			Usage:     "Verify that your client application is generating correct signatures",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "key",
					Usage:    "The signing key",
					Required: true,
				},
				cli.StringFlag{
					Name:     "method",
					Usage:    "HTTP method",
					Required: false,
					Value:    "GET",
				},
				cli.StringFlag{
					Name:     "url",
					Usage:    "Full url without query string and host",
					Required: false,
					Value:    "/",
				},
				cli.StringFlag{
					Name:     "path",
					Usage:    "Path to file for upload to server",
					Required: false,
				},
			},
			Action: func(c *cli.Context) {
				key := c.String("key")
				if key == "" {
					fmt.Fprintf(os.Stderr, "You must provide a key\n")
					os.Exit(1)
				}
				method := c.String("method")
				url := c.String("url")
				path := c.String("path")

				queryString := ""
				if len(c.Args()) > 0 {
					queryString = c.Args()[0]
				}

				sig, err := signature.SignRaw(key, method, url, path, queryString)

				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				appended := signature.AppendSign(queryString, sig)

				fmt.Fprintf(os.Stdout, "Query String: %s\n", queryString)
				fmt.Fprintf(os.Stdout, "Signature: %s\n", sig)
				fmt.Fprintf(os.Stdout, "Signed Query String: %s\n", appended)
			},
		},
	}
	app.Action = func(c *cli.Context) {
		config := c.String("config")

		if config != "" {
			if _, err := os.Stat(config); err != nil {
				fmt.Fprintf(os.Stderr, "Can't find config file `%s`\n", config)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Can't find config file\n")
			os.Exit(1)
		}

		err := server.Run(config)

		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}

	app.Run(os.Args)
}
