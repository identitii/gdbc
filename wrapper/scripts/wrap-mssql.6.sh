#!/bin/bash
set -e

cd "${0%/*}"

DRIVER_NAME="mssql" \
DRIVER_GROUP="com.microsoft.sqlserver" \
DRIVER_ARTIFACT="mssql-jdbc" \
DRIVER_VERSION="6.4.0.jre7" \
DRIVER_CLASS="com.microsoft.sqlserver.jdbc.SQLServerDriver" \
DRIVER_OUTPUT_DIR="`cd ../.. && pwd`/$DRIVER_NAME" \
NATIVE_IMAGE_CONFIG="--allow-incomplete-classpath --report-unsupported-elements-at-runtime --initialize-at-run-time=com.microsoft.sqlserver.jdbc.SQLServerADAL4JUtils" \
  ../wrap-driver.sh

  #--allow-incomplete-classpath --report-unsupported-elements-at-runtime


