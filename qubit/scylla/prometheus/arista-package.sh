#!/usr/bin/env bash

# This script is used to create RPM packages for prometheus server.
# Prometheus server is used to hold stats from Scylla server. The
# stats in prometheus server are used by Scylla dashboards.
# This is a standalone script to be used when a new prometheus
# server release is needed or when the prometheus server configuration
# is updated.

cleanup_exit() {
        OK=${1:-1}
        rm -rf $TMP_BIN_DIR > /dev/null 2>&1
        rm -rf $TMP_CONFIG_DIR > /dev/null 2>&1
        rm -rf ./*.rpm > /dev/null 2>&1
	rm -rf prometheus-${PROMETHEUS_VERSION}.linux-amd64
	if [[ "$1" == "0" ]]; then
		rm -rf prometheus-${PROMETHEUS_VERSION}.linux-amd64.tar.gz
	fi
        exit $OK
}

set -e

PROMETHEUS_VERSION=1.6.3
# version of config RPMs must be incremented when
# a) Prometheus version changes
# b) Configuration is updated
DEV_CONFIG_VERSION=2.0
PROD_CONFIG_VERSION=2.0

VENDOR="Prometheus"
URL="http://gerrit/#/admin/projects/ardc-config/qubit/scylla/prometheus"
LICENSE="ASL 2.0"
ARCH="x86_64"
DESCRIPTION="The Prometheus monitoring system and time series database."
SRC_PROMETHEUS_DIR=prometheus-${PROMETHEUS_VERSION}.linux-amd64

# download and extract ${PROMETHEUS_VERSION}
if [[ ! -f prometheus-${PROMETHEUS_VERSION}.linux-amd64.tar.gz ]]; then
	wget https://github.com/prometheus/prometheus/releases/download/v${PROMETHEUS_VERSION}/prometheus-${PROMETHEUS_VERSION}.linux-amd64.tar.gz
fi
if [[ ! -d ${SRC_PROMETHEUS_DIR} ]]; then
	tar zxvf prometheus-${PROMETHEUS_VERSION}.linux-amd64.tar.gz
fi

# Create Binary RPM
fpm -f -s dir -t rpm -m 'ether-dev@arista.com' -n prometheus --no-depends \
	--license "$LICENSE" \
	--vendor "$VENDOR" \
	--url "$URL" \
	--description "$DESCRIPTION" \
	-a $ARCH \
	--version $PROMETHEUS_VERSION \
	./${SRC_PROMETHEUS_DIR}/prometheus=/usr/bin/prometheus \
	./${SRC_PROMETHEUS_DIR}/promtool=/usr/bin/promtool \
	./${SRC_PROMETHEUS_DIR}/consoles=/usr/share/prometheus/consoles \
	./${SRC_PROMETHEUS_DIR}/console_libraries=/usr/share/prometheus/console_libraries || cleanup_exit 1

# Create prometheus-Dev RPM
fpm -f -s dir -t rpm -m 'ether-dev@arista.com' -n prometheus-Dev --depends prometheus \
	--license "$LICENSE" \
	--vendor "$VENDOR" \
	--url "$URL" \
	--description "Configuration for monitoring Scylla Dev cluster" \
	-a $ARCH \
	--version "$DEV_CONFIG_VERSION" \
	--after-install ./configs/post_install_dev.sh \
	--before-remove ./configs/pre_uninstall_dev.sh \
	--after-remove ./configs/post_uninstall_dev.sh \
	./configs/prometheus_dev.yml=/etc/prometheus/prometheus_dev.yml \
	./configs/scylla_servers_dev.yml=/etc/prometheus/scylla_servers_dev.yml \
	./configs/node_exporter_servers_dev.yml=/etc/prometheus/node_exporter_servers_dev.yml \
	./configs/prometheus-dev.service=/usr/lib/systemd/system/prometheus-dev.service || cleanup_exit 1

# Create prometheus-Prod RPM
fpm -f -s dir -t rpm -m 'ether-dev@arista.com' -n prometheus-Prod --depends prometheus \
	--license "$LICENSE" \
	--vendor "$VENDOR" \
	--url "$URL" \
	--description "Configuration for monitoring Scylla Prod cluster" \
	-a $ARCH \
	--version "$PROD_CONFIG_VERSION" \
	--after-install ./configs/post_install_prod.sh \
	--before-remove ./configs/pre_uninstall_prod.sh \
	--after-remove ./configs/post_uninstall_prod.sh \
	./configs/prometheus_prod.yml=/etc/prometheus/prometheus_prod.yml \
	./configs/scylla_servers_prod.yml=/etc/prometheus/scylla_servers_prod.yml \
	./configs/node_exporter_servers_prod.yml=/etc/prometheus/node_exporter_servers_prod.yml \
	./configs/prometheus-prod.service=/usr/lib/systemd/system/prometheus-prod.service || cleanup_exit 1

mkdir -p RPMS
rm -rf RPMS/*
mv ./*.rpm RPMS

echo "Created RPMS"
ls -l RPMS | awk '{print($9);}'
cleanup_exit 0

