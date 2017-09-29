package mpiw

import (
	"reflect"
	"testing"
)

func TestParseIwDev(t *testing.T) {
	stub := `phy#1
	Interface wlan1
		ifindex 10
		wdev 0x100000002
		addr 00:00:5e:e8:df:ab
		ssid SSID_n
		type AP
		channel 36 (5180 MHz), width: 40 MHz, center1: 5190 MHz
		txpower 14.00 dBm
phy#0
	Interface wlan0
		ifindex 9
		wdev 0x2
		addr 00:00:5e:e8:df:aa
		ssid SSID_g
		type AP
		channel 14 (2484 MHz), width: 20 MHz (no HT), center1: 2484 MHz
		txpower 10.00 dBm`

	interfaces := parseIwDev(stub)
	expected := []string{"wlan1", "wlan0"}
	if !reflect.DeepEqual(expected, interfaces) {
		t.Errorf("Expected %#v, got %#v", expected, interfaces)
	}
}

func TestParseIwDevStationDump(t *testing.T) {
	stub := `Station 00:00:5e:f0:07:87 (on wlan1)
	inactive time:	1600 ms
	rx bytes:	6216368
	rx packets:	50296
	tx bytes:	4006367
	tx packets:	30159
	tx retries:	1159
	tx failed:	2
	rx drop misc:	17
	signal:  	-50 [-49, -57] dBm
	signal avg:	-50 [-49, -56] dBm
	tx bitrate:	243.0 MBit/s MCS 14 40MHz
	rx bitrate:	216.0 MBit/s MCS 13 40MHz
	expected throughput:	53.557Mbps
	authorized:	yes
	authenticated:	yes
	associated:	yes
	preamble:	long
	WMM/WME:	yes
	MFP:		no
	TDLS peer:	no
	DTIM period:	2
	beacon interval:100
	short slot time:yes
	connected time:	195122 seconds
Station 00:00:5e:11:88:2e (on wlan1)
	inactive time:	8140 ms
	rx bytes:	12219038
	rx packets:	98718
	tx bytes:	242662420
	tx packets:	168089
	tx retries:	10161
	tx failed:	365
	rx drop misc:	20
	signal:  	-65 [-70, -67] dBm
	signal avg:	-65 [-69, -67] dBm
	tx bitrate:	6.0 MBit/s
	rx bitrate:	6.0 MBit/s
	expected throughput:	53.557Mbps
	authorized:	yes
	authenticated:	yes
	associated:	yes
	preamble:	long
	WMM/WME:	yes
	MFP:		no
	TDLS peer:	no
	DTIM period:	2
	beacon interval:100
	short slot time:yes
	connected time:	29035 seconds`

	stats := parseIwDevStationDump(stub)
	expected := map[string]map[string]int64{
		"00-00-5e-f0-07-87": map[string]int64{
			"rxBytes":      6216368,
			"txBytes":      4006367,
			"inactiveMsec": 1600,
			"signalDbm":    -50,
		},
		"00-00-5e-11-88-2e": map[string]int64{
			"rxBytes":      12219038,
			"txBytes":      242662420,
			"inactiveMsec": 8140,
			"signalDbm":    -65,
		},
	}
	if !reflect.DeepEqual(expected, stats) {
		t.Errorf("Expected %#v, got %#v", expected, stats)
	}
}
