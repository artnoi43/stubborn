package handler

import (
	"github.com/artnoi43/stubborn/domain/entity"
	"github.com/artnoi43/stubborn/domain/usecase/dotclient"
)

type outbound string

type Config struct {
	Protocol      string           `mapstructure:"protocol" json:"protocol"`
	QueryAllTypes bool             `mapstructure:"query_all_types" json:"queryAllTypes"`
	HostsFile     string           `mapstructure:"hosts_file" json:"hostsFile"`
	Outbound      entity.Outbound  `mapstructure:"outbound" json:"outbound"`
	DoT           dotclient.Config `mapstructure:"dot" json:"dot"`
}

const (
	OutboundDoH outbound = "DOH"
	OutboundDot outbound = "DOT"
)
