package xray

import (
	"encoding/json"
	"fmt"

	"ttcs/database"
)

type Config struct {
	Log      LogConfig       `json:"log,omitempty"`
	Inbounds []InboundConfig `json:"inbounds"`
}

type LogConfig struct {
	Loglevel string `json:"loglevel,omitempty"`
}

type InboundConfig struct {
	Port           int             `json:"port"`
	Protocol       string          `json:"protocol"`
	Settings       InboundSettings `json:"settings"`
	StreamSettings *StreamSettings `json:"streamSettings,omitempty"`
}

type InboundSettings struct {
	Clients []Client `json:"clients"`
}

type Client struct {
	ID      string `json:"id"`
	AlterID int    `json:"alterId"`
}

type StreamSettings struct {
	Network     string       `json:"network,omitempty"`
	Security    string       `json:"security,omitempty"`
	TCPSettings *TCPSettings `json:"tcpSettings,omitempty"`
}

type TCPSettings struct {
	AcceptProxyProtocol bool `json:"acceptProxyProtocol,omitempty"`
}

func BuildConfig(inbounds []database.Inbound) ([]byte, error) {
	config := Config{
		Log: LogConfig{
			Loglevel: "warning",
		},
		Inbounds: []InboundConfig{},
	}

	for _, inbound := range inbounds {
		if !inbound.Enabled {
			continue
		}

		if inbound.Protocol != "vmess" {
			continue
		}

		if inbound.Port <= 0 {
			return nil, fmt.Errorf(
				"invalid port for inbound %q",
				inbound.Remark,
			)
		}

		if inbound.UUID == "" {
			return nil, fmt.Errorf(
				"UUID is empty for inbound %q",
				inbound.Remark,
			)
		}

		item := InboundConfig{
			Port:     inbound.Port,
			Protocol: "vmess",
			Settings: InboundSettings{
				Clients: []Client{
					{
						ID:      inbound.UUID,
						AlterID: inbound.AlterID,
					},
				},
			},
			StreamSettings: &StreamSettings{
				Network:  inbound.Transmission,
				Security: "none",
			},
		}

		if inbound.Transmission == "tcp" {
			item.StreamSettings.TCPSettings = &TCPSettings{
				AcceptProxyProtocol: inbound.TCPAcceptProxyProtocol,
			}
		}

		config.Inbounds = append(
			config.Inbounds,
			item,
		)
	}

	data, err := json.MarshalIndent(
		config,
		"",
		"  ",
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to encode Xray config: %w",
			err,
		)
	}

	return data, nil
}
