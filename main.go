package main

import (
	"fmt"
	"log/slog"
	"os"
	cfg "arb-updater/pkg/config"
	arb "arb-updater/pkg/arb"
)

var (
	Version string
	Build   string
	params  *cfg.Arguments
)

func init() {
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
		if merged, err = ents.Merge(params.Corutines, params.IsSort); err != nil {
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
