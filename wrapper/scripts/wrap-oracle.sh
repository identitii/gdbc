#!/bin/bash
set -e

cd "${0%/*}"

#TODO:Caused by: java.sql.SQLException: Access denied for user 'root'@'172.17.0.1' (using password: NO)
DRIVER_NAME="oracle" \
DRIVER_GROUP="com.oracle" \
DRIVER_ARTIFACT="ojdbc6" \
DRIVER_VERSION="12.1.0.1-atlassian-hosted" \
DRIVER_REPOSITORY="http://repo.spring.io/plugins-release/" \
DRIVER_BUILD_TIME_CLASSES="" \
NATIVE_IMAGE_CONFIG="--allow-incomplete-classpath -H:IncludeResourceBundles=oracle.net.mesg.Message --initialize-at-run-time=oracle.sql.LnxLibServer,oracle.sql.LoadCorejava -H:ReflectionConfigurationFiles=`pwd`/oracle.json -H:ClassInitialization=oracle.jdbc.driver.OracleTimeoutThreadPerVM:rerun" \
DRIVER_OUTPUT_DIR="`cd ../.. && pwd`/$DRIVER_NAME" \
  ../wrap-driver.sh