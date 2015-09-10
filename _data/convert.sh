#!/bin/sh

for HEADER_FILE in fields/*; do
	BASE_FILE=$(basename ${HEADER_FILE})
	../csv2json/csv2json "$HEADER_FILE" "txt/${BASE_FILE}.txt" "json/${BASE_FILE}.json"
done
