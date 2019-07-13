#!/bin/bash
set -e

cd "${0%/*}"

#TODO: com.oracle.svm.core.util.UserError$UserException: only four jsr nesting levels are supported
DRIVER_NAME="mysql" \
DRIVER_GROUP="mysql" \
DRIVER_ARTIFACT="mysql-connector-java" \
DRIVER_VERSION="5.1.47" \
DRIVER_CLASS="com.mysql.jdbc.Driver" \
DRIVER_BUILD_TIME_CLASSES="com.mysql.jdbc.Driver,com.mysql.jdbc.log.StandardLogger,com.mysql.fabric.jdbc.FabricMySQLDriver" \
DRIVER_OUTPUT_DIR="`cd ../.. && pwd`/$DRIVER_NAME" \
NATIVE_IMAGE_CONFIG="-H:IncludeResourceBundles=com.mysql.jdbc.LocalizedErrorMessages -H:ReflectionConfigurationFiles=`pwd`/mysql.5.json" \
  ../wrap-driver.sh