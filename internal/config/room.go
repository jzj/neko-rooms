package config

import (
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"m1k1o/neko_rooms/internal/utils"
)

type Room struct {
	NAT1To1IPs   []string
	EphemeralMin uint16
	EphemeralMax uint16

	TraefikDomain       string
	TraefikEntrypoint   string
	TraefikCertresolver string
	TraefikNetwork      string
}

func (Room) Init(cmd *cobra.Command) error {
	cmd.PersistentFlags().String("epr", "59000-59999", "limits the pool of ephemeral ports that ICE UDP connections can allocate from")
	if err := viper.BindPFlag("epr", cmd.PersistentFlags().Lookup("epr")); err != nil {
		return err
	}

	cmd.PersistentFlags().StringSlice("nat1to1", []string{}, "sets a list of external IP addresses of 1:1 (D)NAT and a candidate type for which the external IP address is used")
	if err := viper.BindPFlag("nat1to1", cmd.PersistentFlags().Lookup("nat1to1")); err != nil {
		return err
	}

	cmd.PersistentFlags().String("traefik_domain", "neko.lan", "traefik: domain on which will be container hosted")
	if err := viper.BindPFlag("traefik_domain", cmd.PersistentFlags().Lookup("traefik_domain")); err != nil {
		return err
	}

	cmd.PersistentFlags().String("traefik_entrypoint", "web-secure", "traefik: router entrypoint")
	if err := viper.BindPFlag("traefik_entrypoint", cmd.PersistentFlags().Lookup("traefik_entrypoint")); err != nil {
		return err
	}

	cmd.PersistentFlags().String("traefik_certresolver", "lets-encrypt", "traefik: certificate resolver for router")
	if err := viper.BindPFlag("traefik_certresolver", cmd.PersistentFlags().Lookup("traefik_certresolver")); err != nil {
		return err
	}

	cmd.PersistentFlags().String("traefik_network", "traefik", "traefik: docker network name")
	if err := viper.BindPFlag("traefik_network", cmd.PersistentFlags().Lookup("traefik_network")); err != nil {
		return err
	}

	return nil
}

func (s *Room) Set() {
	s.NAT1To1IPs = viper.GetStringSlice("nat1to1")

	// if not specified, get public
	if len(s.NAT1To1IPs) == 0 {
		ip, err := utils.GetIP()
		if err == nil {
			s.NAT1To1IPs = append(s.NAT1To1IPs, ip)
		}
	}

	min := uint16(59000)
	max := uint16(59100)
	epr := viper.GetString("epr")
	ports := strings.SplitN(epr, "-", -1)
	if len(ports) > 1 {
		start, err := strconv.ParseUint(ports[0], 10, 16)
		if err == nil {
			min = uint16(start)
		}

		end, err := strconv.ParseUint(ports[1], 10, 16)
		if err == nil {
			max = uint16(end)
		}
	}

	if min > max {
		s.EphemeralMin = max
		s.EphemeralMax = min
	} else {
		s.EphemeralMin = min
		s.EphemeralMax = max
	}

	s.TraefikDomain = viper.GetString("traefik_domain")
	s.TraefikEntrypoint = viper.GetString("traefik_entrypoint")
	s.TraefikCertresolver = viper.GetString("traefik_certresolver")
	s.TraefikNetwork = viper.GetString("traefik_network")
}