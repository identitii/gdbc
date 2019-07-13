#!/bin/bash
set -e

cd "${0%/*}"

DRIVER_NAME="h2" \
DRIVER_GROUP="com.h2database" \
DRIVER_ARTIFACT="h2" \
DRIVER_VERSION="1.4.199" \
DRIVER_CLASS="org.h2.Driver" \
DRIVER_BUILD_TIME_CLASSES="" \
DRIVER_OUTPUT_DIR="`cd .. && pwd`/$DRIVER_NAME" \
NATIVE_IMAGE_CONFIG="--report-unsupported-elements-at-runtime" \
  ../wrap-driver.sh