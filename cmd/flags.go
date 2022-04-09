package cmd

import (
	"flag"
	"log"

	"github.com/artnoi43/stubborn/domain/entity"
)

type Flags struct {
	ConfigFile string
	Outbound   entity.Outbound
	outbound   string // Just for parsing
}

var defaultOutbound = entity.OutboundDoT

func (f *Flags) Parse() {
	flag.StringVar(&f.ConfigFile, "c", ConfLocation, "Path to configuration file")
	flag.StringVar(&f.outbound, "out", "", "Outbound (choose either \"DOT\" or \"DOH\")")
	flag.Parse()

	if len(f.outbound) == 0 {
		f.Outbound = defaultOutbound
	} else {
		o := entity.Outbound(f.outbound)
		if o.IsValid() {
			f.Outbound = o
		} else {
			log.Fatalf("invalid outbound config: %s", o)
		}
	}
}
