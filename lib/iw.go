package mpiw

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	mp "github.com/hnw/go-mackerel-plugin-helper"
	"github.com/hnw/wsoui"
)

// IwPlugin mackerel plugin
type IwPlugin struct {
	prefix string
}

var iwDevHeaderPattern = regexp.MustCompile(
	`^\s+Interface (\w+)`,
)

var iwDevStationDumpHeaderPattern = regexp.MustCompile(
	`^Station ([0-9A-Fa-f]{2}(?::[0-9A-Fa-f]{2}){5}) \(on \w+\)`,
)

var iwDevStationDumpRxBytesPattern = regexp.MustCompile(
	`^\s+rx bytes:\s+(\d+)`,
)

var iwDevStationDumpTxBytesPattern = regexp.MustCompile(
	`^\s+tx bytes:\s+(\d+)`,
)

var iwDevStationDumpInactiveTimePattern = regexp.MustCompile(
	`^\s+inactive time:\s+(\d+)\s+ms`,
)

var iwDevStationDumpSignalDbmPattern = regexp.MustCompile(
	`^\s+signal:\s+(-\d+)\s*\[-\d+,\s*-\d+\]\s*dBm`,
)

// MetricKeyPrefix interface for PluginWithPrefix
func (p IwPlugin) MetricKeyPrefix() string {
	if p.prefix == "" {
		p.prefix = "iw"
	}
	return p.prefix
}

// GraphDefinition interface for mackerelplugin
func (p IwPlugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"interface.#": {
			Label: "Connected Wi-Fi clients",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "connected", Label: "connected"},
			},
		},
		"client.#": {
			Label: "Wi-Fi clients connectivity",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "connected", Label: "connected", Imputation: "zero"},
			},
		},
		"client_transfer_bytes.#": {
			Label: "Wi-Fi transfer",
			Unit:  "bytes/sec",
			Metrics: []mp.Metrics{
				{Name: "rxBytes", Label: "rxBytes", Scale: 0.01666, Diff: true, Imputation: "zero"},
				{Name: "txBytes", Label: "txBytes", Scale: 0.01666, Diff: true, Imputation: "zero"},
			},
		},
		"client_inactive_time.#": {
			Label: "Wi-Fi session inactive time",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "inactiveTime", Label: "inactiveTime", Imputation: "lastValue"},
			},
		},
		"client_signal_power.#": {
			Label: "Wi-Fi signal power",
			Unit:  "float",
			Metrics: []mp.Metrics{
				{Name: "signalDbm", Label: "signal (-dBm)"},
			},
		},
	}
}

// FetchMetrics interface for mackerelplugin
func (p IwPlugin) FetchMetrics() (map[string]interface{}, error) {
	ifNames, err := getWifiInterfaces()
	if err != nil {
		return nil, err
	}

	metrics := make(map[string]interface{})

	for _, ifName := range ifNames {
		cnt := 0
		stats, err := getInterfaceStats(ifName)
		if err != nil {
			return nil, err
		}
		for macaddr, m := range stats {
			abbr, _ := wsoui.LookUp(macaddr)
			readableMac := macaddr
			if len(abbr) > 0 {
				readableMac = abbr + "_" + readableMac
			}
			metrics["client."+readableMac+".connected"] = uint64(1)
			cnt++
			for k, v := range m {
				if k == "inactiveMsec" {
					// convert msec to sec
					metrics["client_inactive_time."+readableMac+".inactiveTime"] = float64(v) / 1000.0
				} else if k == "signalDbm" {
					// negate value
					metrics["client_signal_power."+readableMac+".signalDbm"] = float64(-v)
				} else {
					metrics["client_transfer_bytes."+readableMac+"."+k] = uint64(v)
				}
			}
		}
		metrics["interface."+ifName+".connected"] = uint64(cnt)
	}
	return metrics, nil
}

func getWifiInterfaces() ([]string, error) {
	out, err := getIwDev()
	if err != nil {
		return nil, err
	}
	return parseIwDev(out), nil
}

// Parse output from 'iw dev' and return name of interfaces
func parseIwDev(out string) []string {
	ifNames := make([]string, 0)
	for _, line := range strings.Split(out, "\n") {
		if matches := iwDevHeaderPattern.FindStringSubmatch(line); matches != nil {
			ifNames = append(ifNames, matches[1])
		}
	}
	return ifNames
}

// Run 'iw dev' and return command output
func getIwDev() (string, error) {
	cmd := exec.Command("iw", "dev")
	cmd.Env = append(os.Environ(), "LANG=C")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// Get statistics for specified network interface
func getInterfaceStats(ifName string) (map[string]map[string]int64, error) {
	out, err := getIwDevStationDump(ifName)
	if err != nil {
		return nil, err
	}
	return parseIwDevStationDump(out), nil
}

// Parse output from 'iw dev <ifname> station dump'
func parseIwDevStationDump(out string) map[string]map[string]int64 {
	stats := make(map[string]map[string]int64)
	macaddr := ""
	for _, line := range strings.Split(out, "\n") {
		if matches := iwDevStationDumpHeaderPattern.FindStringSubmatch(line); matches != nil {
			macaddr = strings.Replace(matches[1], ":", "-", -1)
			stats[macaddr] = make(map[string]int64)
		} else if matches := iwDevStationDumpRxBytesPattern.FindStringSubmatch(line); matches != nil {
			stats[macaddr]["rxBytes"], _ = strconv.ParseInt(matches[1], 10, 64)
		} else if matches := iwDevStationDumpTxBytesPattern.FindStringSubmatch(line); matches != nil {
			stats[macaddr]["txBytes"], _ = strconv.ParseInt(matches[1], 10, 64)
		} else if matches := iwDevStationDumpInactiveTimePattern.FindStringSubmatch(line); matches != nil {
			stats[macaddr]["inactiveMsec"], _ = strconv.ParseInt(matches[1], 10, 64)
		} else if matches := iwDevStationDumpSignalDbmPattern.FindStringSubmatch(line); matches != nil {
			stats[macaddr]["signalDbm"], _ = strconv.ParseInt(matches[1], 10, 64)
		}
	}
	return stats
}

// Run 'iw dev <ifname> station dump' and return command output
func getIwDevStationDump(ifName string) (string, error) {
	cmd := exec.Command("iw", "dev", ifName, "station", "dump")
	cmd.Env = append(os.Environ(), "LANG=C")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// Do the plugin
func Do() {
	optPrefix := flag.String("metric-key-prefix", "", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()
	p := IwPlugin{
		prefix: *optPrefix,
	}
	helper := mp.NewMackerelPlugin(p)
	helper.Tempfile = *optTempfile
	helper.Run()
}
