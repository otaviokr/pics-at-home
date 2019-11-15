#!/bin/sh

PICS_AT_HOME_SERVER=$1

sed -i -e "s/@@URL_PICS_AT_HOME@@/${PICS_AT_HOME_SERVER}/g" ./client.html

chromium-browser --kiosk ./client.html
