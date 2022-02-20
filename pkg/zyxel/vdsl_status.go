package zyxel

import (
	"bytes"
	"dslmodem-exporter/pkg/units"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
)

type VDSLStatus struct {
	Status                 string `text:"VDSL Training Status"`
	Mode                   string
	Profile                string `text:"VDSL Profile"`
	TrafficType            string
	LinkUptime             Uptime
	LineRateUp             units.Mbps
	LineRateDown           units.Mbps
	ActualDataRateUp       units.Mbps
	ActualDataRateDown     units.Mbps
	SNRMarginUp            units.Decibel
	SNRMarginDown          units.Decibel
	ActualDelayUp          Delay
	ActualDelayDown        Delay
	TransmitPowerUp        units.DecibelMilliwatt
	TransmitPowerDown      units.DecibelMilliwatt
	ReceivePowerUp         units.DecibelMilliwatt
	ReceivePowerDown       units.DecibelMilliwatt
	AttenuationUp          units.Decibel
	AttenuationDown        units.Decibel
	AttainableDataRateUp   units.Mbps
	AttainableDataRateDown units.Mbps
	// TODO: band status?
	// TODO: counters?
}

const (
	groupStatus = iota
	groupPortDetails
	groupBandStatus
	groupCounters

	knownGroups = groupCounters + 1
)

// UnmarshalText can parse the text output from Zyxel DSL status and convert it to native data in VDSLStatus.
//
// For examples see directory `testdata`.
func (s *VDSLStatus) UnmarshalText(data []byte) error {
	var group uint8
	var groups [knownGroups][]byte
	var parsedHeader bool

	input := bytes.NewBuffer(data)

	// Read all lines and group them
	for {
		line, err := input.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}

		if group >= knownGroups {
			// TODO: maybe store rest of lines as extra data
			log.Warnf("found more than %d groups in data", group)
			break
		}

		if groups[group] == nil {
			groups[group] = []byte{}
		}

		if isSeparator(bytes.TrimSpace(line)) {
			if !parsedHeader {
				// first header will not increment groups
				parsedHeader = true
				continue
			}

			group++
			continue
		}

		groups[group] = append(groups[group], line...)
	}

	var err error

	// Unmarshal each data block
	for i, groupData := range groups {
		switch i {
		case groupStatus:
			err = s.unmarshalStatusText(groupData)
		case groupPortDetails:
			err = s.unmarshalPortDetailsText(groupData)
		case groupBandStatus:
			// TODO: implement?
		case groupCounters:
			// TODO: implement:
		default:
			panic("group not implemented")
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *VDSLStatus) unmarshalStatusText(data []byte) error {
	for _, line := range bytes.Split(data, []byte("\n")) {
		if len(line) == 0 {
			continue
		}

		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) != 2 {
			log.WithField("line", string(line)).Warnf("unknown line")
			continue
		}

		var (
			key   = string(bytes.TrimSpace(parts[0]))
			value = bytes.TrimSpace(parts[1])
		)

		switch key {
		case "VDSL Training Status":
			s.Status = string(value)
		case "Mode":
			s.Mode = string(value)
		case "VDSL Profile":
			s.Profile = string(value)
		case "Traffic Type":
			s.TrafficType = string(value)
		case "Link Uptime":
			err := s.LinkUptime.UnmarshalText(value)
			if err != nil {
				return err
			}
		default:
			log.WithFields(log.Fields{
				"key":   key,
				"value": string(value),
			}).Warn("unknown field in status")
		}
	}

	return nil
}

func (s *VDSLStatus) unmarshalPortDetailsText(data []byte) error {
	var firstLineParsed bool

	for _, line := range bytes.Split(data, []byte("\n")) {
		if len(line) == 0 {
			continue
		}

		// ignore first line
		if !firstLineParsed {
			firstLineParsed = true
			continue
		}

		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) != 2 {
			log.WithField("line", string(line)).Warnf("unknown line")
			continue
		}

		var (
			key    = string(bytes.TrimSpace(parts[0]))
			values = bytes.SplitN(bytes.TrimSpace(parts[1]), []byte("  "), 2)
		)

		// Trim each value
		for i := range values {
			values[i] = bytes.TrimSpace(values[i])
		}

		var errs [2]error

		switch key {
		case "Line Rate":
			errs[0] = s.LineRateUp.UnmarshalText(values[0])
			errs[1] = s.LineRateDown.UnmarshalText(values[1])
		case "Actual Net Data Rate":
			errs[0] = s.ActualDataRateUp.UnmarshalText(values[0])
			errs[1] = s.ActualDataRateDown.UnmarshalText(values[1])
		case "SNR Margin":
			errs[0] = s.SNRMarginUp.UnmarshalText(values[0])
			errs[1] = s.SNRMarginDown.UnmarshalText(values[1])
		case "Actual Delay":
			errs[0] = s.ActualDelayUp.UnmarshalText(values[0])
			errs[1] = s.ActualDelayDown.UnmarshalText(values[1])
		case "Transmit Power":
			errs[0] = s.TransmitPowerUp.UnmarshalText(values[0])
			errs[1] = s.TransmitPowerDown.UnmarshalText(values[1])
		case "Receive Power":
			errs[0] = s.ReceivePowerUp.UnmarshalText(values[0])
			errs[1] = s.ReceivePowerDown.UnmarshalText(values[1])
		case "Total Attenuation":
			errs[0] = s.AttenuationUp.UnmarshalText(values[0])
			errs[1] = s.AttenuationDown.UnmarshalText(values[1])
		case "Attainable Net Data Rate":
			errs[0] = s.AttainableDataRateUp.UnmarshalText(values[0])
			errs[1] = s.AttainableDataRateDown.UnmarshalText(values[1])
		case "Trellis Coding", "Actual INP":
			// ignore
		default:
			log.WithFields(log.Fields{
				"key":    key,
				"values": []string{string(values[0]), string(values[1])},
			}).Warn("unknown field in status")
		}

		// check both errors
		for _, err := range errs {
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isSeparator(line []byte) (result bool) {
	if len(line) == 0 {
		return false
	}

	for _, b := range line {
		if b != '=' {
			return false
		}
	}

	return true
}
