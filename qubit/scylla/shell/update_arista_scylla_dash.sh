#!/usr/bin/env bash

set -e
set -x

# update the scylla dashboards from a directory to grafana

# For production grafana use,
#  ./update_arista_scylla_dash.sh <path> ether-dev@arista.com:ether http://0.0.0.0:3000

JSON_PATH=${1}
GRAFANA_AUTH=${2:-"admin:"}
GRAFANA_URL=${3:-"http://0.0.0.0:3000"}

[[ -z $JSON_PATH ]] && echo "JSON_PATH must be provided" && exit -1

declare -A clusters
clusters[prod]="http://sa101:9090"
clusters[dev]="http://sa105:9091"

for clusterType in ${!clusters[@]}
do
# setup the datasource
curl -XPOST -i -u $GRAFANA_AUTH $GRAFANA_URL/api/datasources \
	-H "Content-Type: application/json" \
	--data-binary '
	{"name":"'"scylla-$clusterType-prom"'",
	 "type":"prometheus",
	 "url":"'"${clusters[$clusterType]}"'",
	 "access":"proxy", "basicAuth":false}'
done

# setup the dashboards
for f in `ls $JSON_PATH/*.json`
do
curl -XPOST -i -u $GRAFANA_AUTH $GRAFANA_URL/api/dashboards/db \
	--data-binary @${f} -H "Content-Type: application/json"
done

echo "Succesfully updated Scylla dashboards"
