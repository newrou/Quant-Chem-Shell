#!/bin/bash

#echo "Path: $1"
echo "Id: $2"
echo "Stat: $3"
#babel -ican $1/$2/mol.can -oxyz --gen3D > $1/$2/v.xyz

n=0
for OUTPUT in $(seq $3)
do
    let n++
    echo "    Generate v$n.xyz"
    babel -ican $1$2/mol.can -oxyz --gen3D > $1$2/v$n.xyz
done

curl "http://localhost:8080/set-status?id=$2&status=prepared"
