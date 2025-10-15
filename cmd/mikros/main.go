package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/log"

	"github.com/mikros-dev/mikros-cli/internal/commands"
	"github.com/mikros-dev/mikros-cli/internal/settings"
)

func main() {
	var (
		ctx     = context.Background()
		options = []fang.Option{
			fang.WithoutVersion(),
			fang.WithoutCompletions(),
		}
	)

	cfg, err := settings.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := fang.Execute(ctx, commands.EntryPoint(cfg), options...); err != nil {
		os.Exit(1)
	}
}
