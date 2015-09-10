#!/bin/sh

for FILE in countryInfo timeZones; do
	if [ ! -f "txt/${FILE}.txt" ]; then
		wget -O "txt/${FILE}.txt" "http://download.geonames.org/export/dump/${FILE}.txt"
	fi
done

for FILE in allCountries; do
	if [ ! -f "zip/${FILE}.zip" ]; then
		wget -O "zip/${FILE}.zip" "http://download.geonames.org/export/dump/${FILE}.zip"
	fi
	if [ ! -f "txt/${FILE}.txt" ]; then
		unzip "zip/${FILE}.zip" -d txt/
	fi
done
