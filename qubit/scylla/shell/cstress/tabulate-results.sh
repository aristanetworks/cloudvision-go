#!/usr/bin/env bash

USAGE="Usage: $0 <test tag>"

TESTRUNTAG=${1?"missing test tag. ${USAGE}"}
BASEDIR="/tmp/cstress-runs/${TESTRUNTAG}"

# cassandra-stress metrics
OPRATE="op rate"
LATENCY_MED="latency median"
LATENCY_95P="latency 95th percentile"
LATENCY_99P="latency 99th percentile"
LATENCY_99_9P="latency 99.9th percentile"
LATENCY_MAX="latency max"
OPTIME="Total operation time"
PARTS="Total partitions"
CMD="CMD:"

# header for the output report
TABLEDATA="Tag~TotalTime~Parts~Op/s~Lat(med)~Lat(95)~Lat(99)~Lat(99.9)~Lat(max)~Timeouts~Command\n"

function extractMetrics {
	NODE=$1
	LOGFILE=$2
	TYPE=$3

        CMD_VALUE=$(awk '/^'"${CMD}"'/ { print($0) }' ${LOGFILE} | cut -d':' -f2- | tr -d '\r' )
        # this ensures that there aren't any whitespaces in the CMD string, again for easy spreadhseeting
        CMD_VALUE=`echo $CMD_VALUE | sed -e 's/\s/_/g'`
        OPRATE_VALUE=$(awk '/^'"${OPRATE}"'/ { print($4) }' ${LOGFILE} )
        L_MED_VALUE=$(awk '/^'"${LATENCY_MED}"'/ { print($4) }' ${LOGFILE} )
        L_95P_VALUE=$(awk '/^'"${LATENCY_95P}"'/ { print($5) }' ${LOGFILE} )
        L_99P_VALUE=$(awk '/^'"${LATENCY_99P}"'/ { print($5) }' ${LOGFILE} )
        L_99_9P_VALUE=$(awk '/^'"${LATENCY_99_9P}"'/ { print($5) }' ${LOGFILE} )
        L_MAX_VALUE=$(awk '/^'"${LATENCY_MAX}"'/ {print($4) }' ${LOGFILE} )
        TIME_VALUE=$(awk '/^'"${OPTIME}"'/ {print($5) }' ${LOGFILE} | tr -d '\r' )
        PARTS_VALUE=$(awk '/^'"${PARTS}"'/ {print($4) }' ${LOGFILE} )
        EXCEPTIONS=$(grep -i -e exception ${LOGFILE} | wc -l)
        echo "${NODE}~${TYPE}$SIZE~$TIME_VALUE~$PARTS_VALUE~$OPRATE_VALUE~$L_MED_VALUE~$L_95P_VALUE~$L_99P_VALUE~$L_99_9P_VALUE~$L_MAX_VALUE~$EXCEPTIONS~$CMD_VALUE\n"
}

function reportMetrics {
	TYPE=$1

	# find nodes
	NODES=""
	TYPE_FILES=`ls ${BASEDIR}/*.${TYPE}.*.log`
	for f in ${TYPE_FILES}
	do
		FULL=`basename $f`
		NODE=${FULL%%.${TYPE}*}
		NODES=`echo ${NODES} ${NODE}`
	done

	for node in ${NODES}
	do
		for SIZE in 512 1024 32768 1048576
	  	do
        		LOGFILE="${BASEDIR}/${node}.${TYPE}.${SIZE}.log"

			[[ ! -f ${LOGFILE} ]] && continue

			TABLEDATA="${TABLEDATA}`extractMetrics ${node} ${LOGFILE} ${TYPE}`"
		done
	done
}

reportMetrics write
reportMetrics read

printf "$TABLEDATA" | column -t -s '~' | sed -e 's/\s\+/ /g'
