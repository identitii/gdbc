#!/bin/bash
set -e

cd "${0%/*}"

DRIVER_NAME="sqlite" \
DRIVER_GROUP="org.xerial" \
DRIVER_ARTIFACT="sqlite-jdbc" \
DRIVER_VERSION="3.27.2.1" \
DRIVER_CLASS="org.sqlite.JDBC" \
DRIVER_OUTPUT_DIR="`cd ../.. && pwd`/$DRIVER_NAME" \
NATIVE_IMAGE_CONFIG="" \
  ../wrap-driver.sh


#DRIVER_BUILD_TIME_CLASSES="com.identitii.gdbc.wrapper.CLIQuery,org.sqlite.JDBC,org.sqlite.core.NativeDB,org.sqlite.util.OSInfo,org.sqlite.core.DB\$ProgressObserver,org.sqlite.SQLiteJDBCLoader\$1,org.sqlite.SQLiteJDBCLoader" 