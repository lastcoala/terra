package main

import (
	"context"
	"log"
	"os"

	"github.com/lastcoala/terra/config"
	"github.com/lastcoala/terra/internal/app"
	"github.com/lastcoala/terra/pkg/handler"
	"github.com/urfave/cli/v3"
)

func main() {
	var configPath string
	var rest bool

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "rest",
				Value:       false,
				Usage:       "run rest server",
				Destination: &rest,
			}, &cli.StringFlag{
				Name:        "config",
				Value:       "config/config.yaml",
				Usage:       "config file path",
				Destination: &configPath,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cfg, err := config.LoadConfig(configPath, "TERRA")
			if err != nil {
				return err
			}

			app := app.NewApp(cfg)
			registry := handler.NewRegistry()

			if rest {
				restHandler := app.CreateRestServer()
				registry.Register("REST", restHandler)
			}

			registry.StartAll()
			registry.StopAll()

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
