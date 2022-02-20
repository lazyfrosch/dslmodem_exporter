# DSL modem exporter for Prometheus

This is still very experimental to 

## Features

Currently, this exporter is build for Zyxel VDSL Modems and tested on VMG1312-B10A.

Maybe we could add more devices and even other vendors in the future.

## Example

```
# HELP dslmodem_collection_seconds Retrieval time for the DSL statistics from the modem.
# TYPE dslmodem_collection_seconds gauge
dslmodem_collection_seconds 1.036668786
# HELP dslmodem_link_attainable_rate_down_bytes Attainable data rate for downstream in bits/s.
# TYPE dslmodem_link_attainable_rate_down_bytes gauge
dslmodem_link_attainable_rate_down_bytes 1.36375e+08
# HELP dslmodem_link_attainable_rate_up_bytes Attainable data rate for upstream in bits/s.
# TYPE dslmodem_link_attainable_rate_up_bytes gauge
dslmodem_link_attainable_rate_up_bytes 3.8745e+07
# HELP dslmodem_link_attenuation_down_dbm Total attenuation for downstream in Decibel.
# TYPE dslmodem_link_attenuation_down_dbm gauge
dslmodem_link_attenuation_down_dbm 7.7
# HELP dslmodem_link_attenuation_up_dbm Total attenuation for upstream in Decibel.
# TYPE dslmodem_link_attenuation_up_dbm gauge
dslmodem_link_attenuation_up_dbm 6.6
# HELP dslmodem_link_rate_down_bytes Rate of the downstream link in bits/s.
# TYPE dslmodem_link_rate_down_bytes gauge
dslmodem_link_rate_down_bytes 1.16799e+08
# HELP dslmodem_link_rate_up_bytes Rate of the upstream link in bits/s.
# TYPE dslmodem_link_rate_up_bytes gauge
dslmodem_link_rate_up_bytes 3.6997e+07
# HELP dslmodem_link_receive_power_down_dbm Current receiving power on downstream in Decibel milliwatt.
# TYPE dslmodem_link_receive_power_down_dbm gauge
dslmodem_link_receive_power_down_dbm 6.4
# HELP dslmodem_link_receive_power_up_dbm Current receiving power on upstream in Decibel milliwatt.
# TYPE dslmodem_link_receive_power_up_dbm gauge
dslmodem_link_receive_power_up_dbm -10.2
# HELP dslmodem_link_snr_margin_down_db Signal to noise margin for downstream in Decibel.
# TYPE dslmodem_link_snr_margin_down_db gauge
dslmodem_link_snr_margin_down_db 11.8
# HELP dslmodem_link_snr_margin_up_db Signal to noise margin for upstream in Decibel.
# TYPE dslmodem_link_snr_margin_up_db gauge
dslmodem_link_snr_margin_up_db 9.9
# HELP dslmodem_link_transmit_power_down_dbm Current transmitting power to downstream in Decibel milliwatt.
# TYPE dslmodem_link_transmit_power_down_dbm gauge
dslmodem_link_transmit_power_down_dbm 14
# HELP dslmodem_link_transmit_power_up_dbm Current transmitting power to upstream in Decibel milliwatt.
# TYPE dslmodem_link_transmit_power_up_dbm gauge
dslmodem_link_transmit_power_up_dbm -3.7
# HELP dslmodem_link_uptime_seconds Since when the link is established.
# TYPE dslmodem_link_uptime_seconds gauge
dslmodem_link_uptime_seconds 45120
# HELP dslmodem_status_info Metadata.
# TYPE dslmodem_status_info gauge
dslmodem_status_info{mode="VDSL2 Annex B",profile="Profile 17a",status="Showtime",traffic_type="PTM Mode"} 1
# HELP dslmodem_up If we are connected to the modem.
# TYPE dslmodem_up gauge
dslmodem_up 1
```

## License

Apache License 2.0, see the full [LICENSE](LICENSE).
