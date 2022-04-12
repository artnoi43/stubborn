package cmd

import (
	"flag"

	"github.com/artnoi43/stubborn/domain/entity"
)

type Flags struct {
	ConfigFile     string
	TableFile      string
	EntityOutbound entity.Outbound
	outbound       string // Just for parsing
}

// Parse only parses the command-line flags.
// Validity check is done in package config or handler.
func (f *Flags) Parse() {
	flag.StringVar(&f.ConfigFile, "c", ConfLocation, "Path to configuration file")
	flag.StringVar(&f.TableFile, "t", TableLocation, "Path to local hosts file")
	flag.StringVar(&f.outbound, "o", "", "Outbound (choose either \"DOT\" or \"DOH\")")
	flag.Parse()

	if f.outbound != "" {
		f.EntityOutbound = entity.OutboundFromString(f.outbound)
	}
}
