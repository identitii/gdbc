#!/bin/bash
set -e
set -x

cd "${0%/*}"

echo "Wrapping driver $DRIVER_NAME (class: $DRIVER_CLASS). Output: $DRIVER_OUTPUT_DIR"

echo "MAVEN_OPTS=$MAVEN_OPTS"

mvn clean install -Dmaven.test.skip=true -Dorg.slf4j.simpleLogger.log.org.apache.maven.cli.transfer.Slf4jMavenTransferListener=WARN

rm -rf target/lib || true
mkdir -p target/lib

# mvn -Dtransitive=true -Doptional=true dependency:copy-dependencies -DoutputDirectory="lib" -DDRIVER_GROUP="$DRIVER_GROUP" -DDRIVER_ARTIFACT="$DRIVER_ARTIFACT" -DDRIVER_VERSION="$DRIVER_VERSION" -f pom-driver.xml

JAR_DIR=target/lib/$DRIVER_NAME
mkdir -p $JAR_DIR

mvn -DgroupId=$DRIVER_GROUP -DartifactId=$DRIVER_ARTIFACT -Dversion=$DRIVER_VERSION -Ddest=$JAR_DIR/ -Dtransitive=true -DremoteRepositories=$DRIVER_REPOSITORY dependency:get

JARS="`ls -m $JAR_DIR/*.jar`"
JARS="${JARS//, /:}"

# Get the (non-optional) classpath by using a fake intermediate maven project
CLASSPATH="`mvn -f pom2.xml -q exec:exec -Dexec.executable=echo -Dexec.args='%classpath'`:$JARS"

if ! [ -x "$(command -v zip)" ]; then
  echo "zip isnt installed, trying to install via apt"
  apt-get update && apt-get install -y zip
fi

zip $JARS META-INF/MANIFEST.MF # fails > 1 jar

# Take out class dir from the beginning, and add to the end. That way, we can override whatever we like. I think. Maybe?
# CLASS_DIR=`cd target/classes && pwd`
# CLASSPATH="${CLASSPATH//$CLASS_DIR:/}:$CLASS_DIR"

# pushd target/lib
# unzip -q *.jar
# find . | grep "\.class" | sed 's/\.class//' | sed 's/\.\///' | sed 's/\//\./g' | jq -R -s 'split("\n") | map(select(length > 0))| map({name:.})' > ../reflect-config.json
# popd

export LIBRARY_NAME="libgdbc-$DRIVER_NAME"

# if [ -d "./graalvm-ce-19.0.0" ]; then
    # use local if there...
    export GRAAL_HOME=`pwd`/graalvm-ce-19.1.0
# else
#     export GRAAL_HOME=~/Downloads/graalvm-ee-19.0.0
# fi
export JAVA_HOME=$GRAAL_HOME/Contents/Home
export PATH=$JAVA_HOME/bin:$PATH

REFLECTION_CONFIG="config/default.json,$REFLECTION_CONFIG"

# TODO: is there an output directory option?
# NATIVE_IMAGE_CONFIG="--no-fallback --no-server --initialize-at-build-time=\"$DRIVER_CLASS,$DRIVER_BUILD_TIME_CLASSES,com.identitii.gdbc.wrapper.DriverWrapper\" -H:ReflectionConfigurationFiles=$REFLECTION_CONFIG"

#-H:ClassInitialization=$DRIVER_CLASS_INITIALIZATION
# echo "$NATIVE_IMAGE_CONFIG"
#--report-unsupported-elements-at-runtime --allow-incomplete-classpath 

#--no-server 

# -H:+GeneratePIC should allow us to use upx to compress the library. but it doesn't
# -H:IncludeResources=".*" 
# native-image -J-ea $NATIVE_IMAGE_CONFIG --enable-all-security-services -H:+AddAllCharsets -H:IncludeResources=".*" --no-fallback --no-server -H:+ReportExceptionStackTraces --initialize-at-build-time="$DRIVER_BUILD_TIME_CLASSES"                                          -cp "$CLASSPATH" -H:Name="$LIBRARY_NAME-query" SimpleTest 
native-image -J-ea $NATIVE_IMAGE_CONFIG --enable-all-security-services -H:+AddAllCharsets --no-fallback --no-server -H:+ReportExceptionStackTraces --initialize-at-build-time="$DRIVER_BUILD_TIME_CLASSES,com.identitii.gdbc.wrapper.DriverWrapper" -cp "$CLASSPATH" -H:Name="$LIBRARY_NAME" --shared --verbose

LIBRARY_FILE=$LIBRARY_NAME.dylib
if [ -f "$LIBRARY_FILE" ]; then
    # on mac... update the dylib install name to load from current directory

    otool -L $LIBRARY_FILE
    # install_name_tool -add_rpath @executable_path/. $LIBRARY_FILE
    # install_name_tool -add_rpath @executable_path/lib $LIBRARY_FILE
    # install_name_tool -add_rpath @executable_path/postgres $LIBRARY_FILE
    install_name_tool -id @executable_path/$LIBRARY_FILE $LIBRARY_FILE

    otool -L $LIBRARY_FILE

    # TODO: Linux?
fi

[ -d $DRIVER_OUTPUT_DIR ] || mkdir -p $DRIVER_OUTPUT_DIR
mv $LIBRARY_NAME* $DRIVER_OUTPUT_DIR
mv graal_isolate* $DRIVER_OUTPUT_DIR
cp *.go $DRIVER_OUTPUT_DIR
sed -i -e "s/postgresql/$DRIVER_NAME/g" $DRIVER_OUTPUT_DIR/*.go
rm $DRIVER_OUTPUT_DIR/*.go-e || true

