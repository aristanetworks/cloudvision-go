#!/usr/bin/env bash

# read the tests from tests.list and run them using ansible

set -e

USAGE="Usage: $0 <test tag> <test nodes> <contact nodes> [duration]"

TEST_TAG="${1?${USAGE}}"
TEST_NODES="${2?${USAGE}}"
CONTACT_NODES="${3?${USAGE}}"
TEST_DURATION=${4-"15m"}
TS_CMD="\`date +%F_%T\`"
TEST_BASE_DIR="cstress-runs"
LOCAL_TEST_BASE_DIR="/tmp/${TEST_BASE_DIR}"
REMOTE_TEST_BASE_DIR="${TEST_BASE_DIR}"
TEST_LIST_INPUT="tests.list"
MAIN_LOGFILE="${LOCAL_TEST_BASE_DIR}/${TEST_TAG}/main.log"
ANSIBLE_ROOT="../../ansible"
TESTSUITE_CMDS="${LOCAL_TEST_BASE_DIR}/${TEST_TAG}/cmds.file"

mkdir -p `dirname $MAIN_LOGFILE`
> ${MAIN_LOGFILE}

> ${TESTSUITE_CMDS}
chmod a+x ${TESTSUITE_CMDS}

echo "Generating commands in test suite..."
echo "#!/usr/bin/env bash" >> ${TESTSUITE_CMDS}

cat ${TEST_LIST_INPUT} | while read -r TEST_PARAMS || [[ -n ${TEST_PARAMS} ]]
do
  # skip commented out tests
  [[ ${TEST_PARAMS} =~ ^# ]] && continue

  NAME=`echo $TEST_PARAMS | awk '{print($1)}'`

  CMD1="echo \"${TS_CMD} TestBegin ${NAME}\" >> $MAIN_LOGFILE"
  CMD2="${ANSIBLE_ROOT}/runner ${ANSIBLE_ROOT}/run_cstress_test.yml ${TEST_NODES}, \"${TEST_PARAMS} LOCAL_TEST_BASE_DIR=${LOCAL_TEST_BASE_DIR} REMOTE_TEST_BASE_DIR=${REMOTE_TEST_BASE_DIR} TEST_DURATION=${TEST_DURATION} TEST_TAG=${TEST_TAG} CONTACT_POINTS=${CONTACT_NODES}\""
  CMD3="[[ \$? != 0 ]] && ( echo \"${TS_CMD} TesSuite FAIL\" >> ${MAIN_LOGFILE} ) && exit 1"
  CMD4="echo \"${TS_CMD} TestEnd ${NAME}\" >> $MAIN_LOGFILE"

  echo $CMD1 >> ${TESTSUITE_CMDS}
  echo $CMD2 >> ${TESTSUITE_CMDS}
  echo $CMD3 >> ${TESTSUITE_CMDS}
  echo $CMD4 >> ${TESTSUITE_CMDS}
done
echo "echo \"${TS_CMD} TestSuite PASS\" >> $MAIN_LOGFILE" >> ${TESTSUITE_CMDS}

echo "Running commands in test suite..."
${TESTSUITE_CMDS}
echo "Completed commands in test suite. Check ${MAIN_LOGFILE}."
