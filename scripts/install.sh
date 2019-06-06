#!/bin/bash

BASEDIR=$(dirname $BASH_SOURCE)

mkdir -p /opt/custom/smf/share
cp $BASEDIR/../smf/lofs_overlay.xml /opt/custom/smf
cp $BASEDIR/../bin/smos/lofs_overlay /opt/custom/smf/share

mkdir -p /usbkey/crud/root
