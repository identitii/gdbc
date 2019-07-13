#!/bin/bash
set -e

cd "${0%/*}"

DRIVER_NAME="postgresql" \
DRIVER_GROUP="org.postgresql" \
DRIVER_ARTIFACT="postgresql" \
DRIVER_VERSION="42.2.5.jre6" \
DRIVER_CLASS="org.postgresql.Driver" \
DRIVER_BUILD_TIME_CLASSES="com.sun.jna.platform.win32.Win32Exception,com.sun.jna.LastErrorException,org.postgresql.Driver,org.postgresql.core.Logger,org.postgresql.util.SharedTimer" \
DRIVER_OUTPUT_DIR="`cd ../.. && pwd`/$DRIVER_NAME" \
NATIVE_IMAGE_CONFIG="--report-unsupported-elements-at-runtime --allow-incomplete-classpath" \
  ../wrap-driver.sh
