#!/bin/bash
set -e
executable=${executable:-"go run main.go"}

ndays=$(${executable} --count-days)
day=0
while [[ day -lt ndays ]]; do
    echo "::group::Day ${day}"
    ${executable} --day="${day}" 2>&1
    echo "::endgroup::"
    echo
    ((++day))
done