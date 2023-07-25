Webhook
=====

A simple auto deploy tools.


# Requirement

* setup GOPATH and GOROOT environment variables
* go 1.18+
* ubuntu / linux / debian

# Installation (for Ubuntu/Linux)

1. install binary

```shell
go install github.com/juxuny/webhook@latest
```

2. create a linux service config file in /etc/systemd/system

```shell
[Unit]
Description=Webhook
After=network.target

[Service]
Type=simple
User=juxuny
Restart=on-failure
RestartSec=5s
ExecStart=/mnt/tools/gopath/bin/webhook serve -c /home/juxuny/.config/webhook/config.yaml

[Install]
WantedBy=multi-user.target
```

*Notice*: You should modify the real path for yourself.

3. create a YAML config file, look like this:

```yaml
port: 20089
url-prefix: /hook
auth:
  - name: nuc
    token: {{token}}
deployments:
  - name: test-api-dev
    bash-interpreter: /bin/bash
    work-dir: .
    variables:
      - name: tag
        nullable: false
        type: string
    scripts:
      - test.sh

```

4. enable the system service

```shell
sudo systemctl daemon-reload
sudo systemctl enable webhook
```

5. check running status

```shell
sudo service webhook status
```

6. trigger the deployment

```shell
curl -X POST -d 'tag=test-api-dev' http://127.0.0.1:20089/hook/test-api-dev
```