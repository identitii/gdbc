#!/bin/bash
set -e

cd "${0%/*}"

DRIVER_NAME="mssql" \
DRIVER_GROUP="com.microsoft.sqlserver" \
DRIVER_ARTIFACT="mssql-jdbc" \
DRIVER_VERSION="7.3.0.jre8-preview" \
DRIVER_CLASS="com.microsoft.sqlserver.jdbc.SQLServerDriver" \
DRIVER_OUTPUT_DIR="`cd ../.. && pwd`/$DRIVER_NAME" \
DRIVER_BUILD_TIME_CLASSES="com.microsoft.sqlserver.jdbc.SQLServerDriver,com.microsoft.sqlserver.jdbc.SQLServerResource,com.microsoft.sqlserver.jdbc.SQLServerDriverPropertyInfo" \
NATIVE_IMAGE_CONFIG="--allow-incomplete-classpath --report-unsupported-elements-at-runtime --initialize-at-run-time=com.microsoft.sqlserver.jdbc.SQLServerADAL4JUtils -H:IncludeResourceBundles=com.microsoft.sqlserver.jdbc.SQLServerResource" \
  ../wrap-driver.sh

  #--allow-incomplete-classpath --report-unsupported-elements-at-runtime