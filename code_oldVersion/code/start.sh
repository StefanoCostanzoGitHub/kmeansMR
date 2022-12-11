#!/bin/bash

if [ $# -ne 4 ] 
then 
    echo "Usage: ${0} ALGO NMAPPERS MAXITER THRESHOLD(1 = 0.001%)"
    exit 1
fi

if [ $2 -lt 0 -o $2 -gt 99 ]
then
	echo "The number of mappers isn't correct or is too high."
	exit 1
fi


if [ $3 -lt 0 -o $3 -gt 1000 ]
then
	echo "The number of max iterations isn't correct or is too high."
	exit 1
fi

if [ $4 -lt 0 -o $4 -gt 999 ]
then
	echo "DeltaThreshold must be included in [0,100]."
	exit 1
fi

if [ $1 -eq 1 ]
then
    NUMMAP=${2} MAXITER=${3} THRESHOLD=${4} docker compose up
else
    echo -e "The algorithm must be 1 or 2"
fi
