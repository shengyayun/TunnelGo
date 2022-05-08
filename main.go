package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"tunnel/internal"
)

func main() {
	app := cli.NewApp()
	app.Usage = "Connect everywhere"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   "./tunnel.yaml",
			Usage:   "make a reference to tunnel.yaml.example",
		},
	}

	app.Action = func(c *cli.Context) error {
		path := c.String("config")
		config, err := internal.NewConfig(path)
		if err != nil {
			log.Fatalf("Failed to read config: %v\n", err)
		}
		for name, tunnel := range config.Tunnels {
			go func(name string, tunnel internal.Tunnel) {
				tunnel.Name = name
				tunnel.Run()
			}(name, tunnel)
		}
		select {}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
