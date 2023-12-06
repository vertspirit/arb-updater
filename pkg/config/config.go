package config

import (
	"flag"
)

type Arguments struct {
	TemplateArbFile string
	LocaleArbFile   string
	OutputArbFile   string
	Corutines       int
	FullOn          bool
	IsSort          bool
	PrintOnly       bool
	PrintVersion    bool
}

func ParseArguments() *Arguments {
	params := &Arguments{}
	flag.StringVar(&params.TemplateArbFile, "template-arb", "intl_en.arb", "The template arb file to update locale arb file.")
	flag.StringVar(&params.TemplateArbFile, "t", "intl_en.arb", "The template arb file to update locale arb file. (shorten)")
	flag.StringVar(&params.LocaleArbFile, "locale-arb", "", "The locale arb file to update by template file.")
	flag.StringVar(&params.LocaleArbFile, "l", "", "The locale arb file to update by template file. (shorten)")
	flag.StringVar(&params.OutputArbFile, "output-arb", "", "The output arb file which the merged entries saved to.")
	flag.StringVar(&params.OutputArbFile, "o", "", "The output arb file which the merged entries saved to. (shorten)")
	flag.IntVar(&params.Corutines, "corutines", 2048, "The max jobs to enforce concurrently. Default is 2048.")
	flag.IntVar(&params.Corutines, "c", 2048, "The max jobs to enforce concurrently. Default is 2048. (shorten)")
	flag.BoolVar(&params.FullOn, "full-on", false, "Run the jobs without max corutines limit. The same as set corutines to zero.")
	flag.BoolVar(&params.FullOn, "f", false, "Run the jobs without max corutines limit. The same as set corutines to zero. (shorten)")
	flag.BoolVar(&params.IsSort, "sort", false, "Sort the merged entries by name.")
	flag.BoolVar(&params.IsSort, "s", false, "Sort the merged entries by name. (shorten)")
	flag.BoolVar(&params.PrintOnly, "print-only", false, "Print the merged entries in console only.")
	flag.BoolVar(&params.PrintOnly, "p", false, "Print the merged entries in console only (shorten)")
	flag.BoolVar(&params.PrintVersion, "version", false, "Print the version information.")
	flag.BoolVar(&params.PrintVersion, "v", false, "Print the version information. (shorten)")
	flag.Parse()
	return params
}
