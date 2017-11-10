#!/usr/bin/env bash

set -e
set -x

# customize the scylla dashboards for different Arista clusters
# Usage: ./ccreate_arista_scylla_dash.sh /tmp/scylla_dash_dir HEAD 2017.1

TEMP_DASH_DIR=${1:-"/tmp/scylla_dash_dir"}
COMMIT=${2:-"HEAD"}
DASH_VERSION=${3:-"2017.1"}

SCYLLA_DASH_REPO="https://github.com/scylladb/scylla-grafana-monitoring"
JSON_FILES="scylla-dash scylla-dash-io-per-server scylla-dash-per-server"
TEMP_REPO_DIR=`mktemp -d`

function do_arista_customizations() {

	FILENAME=$1
	CLUSTER=`echo $FILENAME | cut -d"." -f2`

	sed -i -e 's/"datasource": "prometheus"/"datasource": "scylla-'"$CLUSTER"'-prom"/g' $FILENAME
	sed -i -e 's/'"${DASH_VERSION}"'/'"$CLUSTER"'/g' $FILENAME

	# enable dashboard overwrite
	# add a message with update details
	sed -i '$ s/}$/,\
"overwrite": true,\
"message": "dashboard updated to '"${DASH_VERSION}"'"\
}/' $FILENAME

	# use proper mountpoint /persist2 rather than /var/lib/scylla
	sed -i 's/\/var\/lib\/scylla/\/persist/g' $FILENAME

	# TODO(krishna:) need to figure out a way to change the json
	# select md2 as the drive by default
	# select bond0 network interface as default
}

mkdir -p ${TEMP_DASH_DIR}

git clone ${SCYLLA_DASH_REPO} ${TEMP_REPO_DIR}
cd ${TEMP_REPO_DIR}
git checkout ${COMMIT}

# create copies, one set per cluster
for f in ${JSON_FILES}
do
	cp grafana/${f}.${DASH_VERSION}.json ${TEMP_DASH_DIR}/${f}.prod.json
	cp grafana/${f}.${DASH_VERSION}.json ${TEMP_DASH_DIR}/${f}.dev.json
done

# customize the json files per cluster - datasource, tags, title etc
for f in `ls ${TEMP_DASH_DIR}/*.json`
do
	do_arista_customizations $f
done

rm -rf ${TEMP_REPO_DIR}
echo "Generated Arista Scylla dashboards in ${TEMP_DASH_DIR}"

