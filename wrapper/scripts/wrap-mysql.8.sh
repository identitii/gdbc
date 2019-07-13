#!/bin/bash
set -e

cd "${0%/*}"

#TODO: Can't find com.mysql.cj.protocol.x.XMessage
DRIVER_NAME="mysql" \
DRIVER_GROUP="mysql" \
DRIVER_ARTIFACT="mysql-connector-java" \
DRIVER_VERSION="8.0.16" \
DRIVER_CLASS="com.mysql.cj.jdbc.Driver" \
DRIVER_BUILD_TIME_CLASSES="com.mysql.cj.jdbc.Driver,com.mysql.cj.jdbc.AbandonedConnectionCleanupThread" \
NATIVE_IMAGE_CONFIG="--allow-incomplete-classpath --report-unsupported-elements-at-runtime" \
DRIVER_OUTPUT_DIR="`cd ../.. && pwd`/$DRIVER_NAME" \
  ../wrap-driver.sh

  #-H:ClassInitialization=com.mysql.cj.jdbc.NonRegisteringDriver:rerun