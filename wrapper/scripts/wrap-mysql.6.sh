#!/bin/bash
set -e

cd "${0%/*}"

#TODO: java.lang.ClassCastException: com.mysql.cj.core.exceptions.CJException cannot be cast to com.mysql.cj.core.exceptions.WrongArgumentException
DRIVER_NAME="mysql" \
DRIVER_GROUP="mysql" \
DRIVER_ARTIFACT="mysql-connector-java" \
DRIVER_VERSION="6.0.6" \
DRIVER_CLASS="com.mysql.cj.jdbc.Driver" \
DRIVER_BUILD_TIME_CLASSES="" \
DRIVER_CLASS_INITIALIZATION="com.mysql.cj.jdbc.AbandonedConnectionCleanupThread -H:ReflectionConfigurationFiles=`pwd`/mysql.6.json" \
DRIVER_OUTPUT_DIR="`cd ../.. && pwd`/$DRIVER_NAME" \
  ../wrap-driver.sh

  #NATIVE_IMAGE_CONFIG="--allow-incomplete-classpath --report-unsupported-elements-at-runtime" 