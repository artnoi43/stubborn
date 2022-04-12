package handler

import (
	"github.com/artnoi43/stubborn/domain/entity"
	"github.com/artnoi43/stubborn/domain/usecase/clients/dohclient"
	"github.com/artnoi43/stubborn/domain/usecase/clients/dotclient"
)

type Config struct {
	Protocol      string           `mapstructure:"protocol" json:"protocol"`
	QueryAllTypes bool             `mapstructure:"query_all_types" json:"queryAllTypes"`
	HostsFile     string           `mapstructure:"hosts_file" json:"hostsFile"`
	DoT           dotclient.Config `mapstructure:"dot" json:"dot"`
	DoH           dohclient.Config `mapstructure:"doh" json:"doh"`
	// entity.Outbound is not in config file,
	// instead stubborn convert Outbound (string) to
	// corresponding int enum value.
	Outbound       string          `mapstructure:"outbound" json:"outbound"`
	EntityOutbound entity.Outbound `mapstructure:"entity_outbound" json:"-"`
}
