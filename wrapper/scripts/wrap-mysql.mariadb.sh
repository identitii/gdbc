#!/bin/bash
set -e

cd "${0%/*}"

#TODO:Caused by: java.sql.SQLException: Access denied for user 'root'@'172.17.0.1' (using password: NO)
DRIVER_NAME="mysql" \
DRIVER_GROUP="org.mariadb.jdbc" \
DRIVER_ARTIFACT="mariadb-java-client" \
DRIVER_VERSION="2.4.1" \
DRIVER_CLASS="org.mariadb.jdbc.Driver" \
DRIVER_BUILD_TIME_CLASSES="" \
NATIVE_IMAGE_CONFIG="-H:ReflectionConfigurationFiles=`pwd`/mysql.mariadb.json -H:ClassInitialization=sun.java2d.opengl.OGLRenderQueue:rerun,sun.awt.AWTAutoShutdown:rerun,sun.java2d.opengl.OGLContext:rerun,sun.awt.CGraphicsEnvironment:rerun,java.awt.GraphicsEnvironment:rerun" \
DRIVER_OUTPUT_DIR="`cd ../.. && pwd`/$DRIVER_NAME" \
  ../wrap-driver.sh

  #-H:ClassInitialization=com.mysql.cj.jdbc.NonRegisteringDriver:rerun

  #--allow-incomplete-classpath --report-unsupported-elements-at-runtime