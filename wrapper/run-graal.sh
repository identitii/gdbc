mvn compile

CLASSPATH="`mvn -q exec:exec -Dexec.executable=echo -Dexec.args='%classpath'`:`cd target/classes && pwd`"

native-image --no-server --report-unsupported-elements-at-runtime --allow-incomplete-classpath '-H:IncludeResources=.*' --no-fallback -H:+ReportExceptionStackTraces --initialize-at-build-time=com.sun.jna.platform.win32.Win32Exception,com.sun.jna.LastErrorException,org.postgresql.Driver,org.postgresql.core.Logger,org.postgresql.util.SharedTimer -cp $CLASSPATH -H:Name=simpletest SimpleTest
./simpletest "org.postgresql.Driver" "jdbc:postgresql://localhost/test?loggerLevel=DEBUG" "root" "password" "SELECT 1;"