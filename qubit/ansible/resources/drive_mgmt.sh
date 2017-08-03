#!/usr/bin/bash

shopt -s nocasematch

declare -A drive_handlers

# setup drive handler routines for different vendors
drive_handlers[MICRON_CHECK_FW]=check_micron_needs_fwupdate
drive_handlers[INTEL_CHECK_FW]=check_intel_needs_fwupdate

# Assumes that the server has only SCSI (-S option) SSD drives
ssd_drives=`lsblk -S -l -n -o NAME`

# drive_info_columns and drive_info_fields must match
drive_info_fields="Serial Number|Device Model|Firmware Version|User Capacity"
drive_info_columns="- - - -"

# Micron SSD tool doesn't have an option to check
# if a newer FW version is available. We need to
# check the download website for newer versions,
# download them and do the FW update.
function check_micron_needs_fwupdate() {
	echo "CHECKWEB"
}

function micron_fwupdate() {
	echo "FW update not implemented"
}

function check_intel_needs_fwupdate() {
	drive=$1
	vendor=$2
	serial=$3

	status=`isdct show -d FirmwareUpdateAvailable -intelssd $serial | grep FirmwareUpdateAvailable`
	if [[ "$status" =~ .+"Intel SSD contains current firmware".+ ]]
	then
		echo "CURRENT"
	else
		echo "NEEDS_UPDATE"
	fi
}

function intel_fwupdate() {
   echo "FW update not implemented"
}

function fw_needs_update() {
  drive=$1
  vendor=$2
  serial=$3

  handler=${drive_handlers[${vendor}_CHECK_FW]}
  $handler $drive $vendor $serial
}

function fw_update() {
  drive=$1
  vendor=$2
  serial=$3

  handler=${drive_handlers[${vendor}_UPDATE_FW]}
  $handler $drive $vendor $serial
}

function report() {
  for drive in ${ssd_drives}
  do
    vendor="INTEL"

    drive_info=`smartctl -i /dev/${drive} | egrep "${drive_info_fields}"`
    drive_info_line=`echo ${drive_info} | paste ${drive_info_columns}`

    serial=`echo "${drive_info}" | grep 'Serial Number' | awk -F: '{print($2)}' | sed -e 's/^\s\+//g'`
    if [[ ${drive_info_line} =~ Micron ]]
    then
      vendor="MICRON"
    fi

    fw_needs_update=`fw_needs_update $drive $vendor $serial`

    echo ${drive} "FW: ${fw_needs_update}" ${drive_info_line}
  done
}

function check_update() {
  echo "FW update not implemented"
}

function usage() {
  echo "Usage: $0 report|update"
  echo
  echo "report: get a summary of drive information"
  echo "update: do a FW update on the drives, if necessary"
}

# -- main --

if [[ $# -ne 1 ]]
then
  usage
  exit
fi

case $1 in
  'report') report;;
  'update') check_update;;
  '-h') usage;;
  *) echo "unsupported command";;
esac

