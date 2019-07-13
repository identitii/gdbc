
#!/bin/bash
set -e

cd "${0%/*}"

DRIVER_NAME="postgresql" \
DRIVER_GROUP="postgresql" \
DRIVER_ARTIFACT="postgresql" \
DRIVER_VERSION="9.1-901-1.jdbc4" \
DRIVER_CLASS="org.postgresql.Driver" \
DRIVER_BUILD_TIME_CLASSES="" \
DRIVER_OUTPUT_DIR="`cd ../.. && pwd`/$DRIVER_NAME" \
  ../wrap-driver.sh
