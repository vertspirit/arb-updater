package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"
	cfg "arb-updater/pkg/config"
	arb "arb-updater/pkg/arb"
)

var (
	Version string
	Build   string
	params  *cfg.Arguments
)

func init() {
	testing.Init()

	params = cfg.ParseArguments()
	if params.PrintVersion {
		fmt.Printf("Arb-Updater Version %s\nBuilt %s\n", Version, Build)
		os.Exit(0)
	}
	if params.LocaleArbFile == "" {
		slog.Error("The '-locale-arb/-l' parameter is required. Run 'arb-updater -h' for details.")
		os.Exit(1)
	}
	if params.OutputArbFile == "" {
		params.OutputArbFile = params.LocaleArbFile
	}
}

func main() {
	update := arb.NewArbUpdate(params.TemplateArbFile, params.LocaleArbFile, params.OutputArbFile)
	ents, err := update.ReadAllEntries()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	var merged []byte
	if params.FullOn || params.Corutines == 0 {
		if merged, err = ents.MergeFullOn(params.IsSort); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	} else {
		ctx := context.Background()
		if merged, err = ents.Merge(ctx, params.Corutines, params.IsSort); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}

	if params.PrintOnly {
		arb.PrintEntries(merged)
		return
	}

	if err = arb.SaveEntriesToFile(update, merged); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
