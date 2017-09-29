mackerel-plugin-iw
=====================

iw custom metrics plugin for mackerel.io agent.

## Build

```shell
$ go build
```

## Synopsis

```shell
$ mackerel-plugin-iw
x9x99x.interface.wlan1.connected	2	1506618443
x9x99x.interface.wlan0.connected	0	1506618443
x9x99x.client.Apple_00-00-5e-31-38-cc.connected	1	1506618443
x9x99x.client.SamsungE_00-00-5e-11-88-2e.connected	1	1506618443
x9x99x.client_transfer_bytes.Apple_00-00-5e-31-38-cc.rxBytes	699.320160	1506618443
x9x99x.client_transfer_bytes.Apple_00-00-5e-31-38-cc.txBytes	2140.143600	1506618443
x9x99x.client_transfer_bytes.SamsungE_00-00-5e-11-88-2e.rxBytes	0.000000	1506618443
x9x99x.client_transfer_bytes.SamsungE_00-00-5e-11-88-2e.txBytes	0.000000	1506618443
x9x99x.client_inactive_time.Apple_00-00-5e-31-38-cc.inactiveTime	0.210000	1506618443
x9x99x.client_inactive_time.SamsungE_00-00-5e-11-88-2e.inactiveTime	7.580000	1506618443
x9x99x.client_signal_power.Apple_00-00-5e-31-38-cc.signalDbm	64.000000	1506618443
x9x99x.client_signal_power.SamsungE_00-00-5e-11-88-2e.signalDbm	43.000000	1506618443
```

## Example of mackerel-agent.conf

```
[plugin.metrics.iw]
command = "/path/to/mackerel-plugin-iw"
```
