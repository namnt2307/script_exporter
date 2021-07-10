# SCRIPT_EXPORTER
## Description
Script exporter will get exit code from executed script and return prometheus metrics "script_success{name="script_name"}". If the script run without any error, prometheus will return "script_success=0", otherwise the metric for failed script will be "script_success=1"  
### Example config file

``````````````
scripts:
  - name: main.py
    dir: ./main.py
    executor: /usr/bin/python3
  - name: a.sh
    dir: ./a.sh
    executor: /bin/bash
``````````````
### Run
````
go run main.go /path/to/config.yml
```` 
### Setup as systemd service
````
[Unit]
Description=Script exporter
Wants=network-online.target
After=network-online.target
[Service]
Type=simple
ExecStart=/root/script_exporter/main /root/script_exporter/config.yml
Restart=on-failure
RestartSec=10
[Install]
WantedBy=default.target
````
Version:
- Go: 1.16