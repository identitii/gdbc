mvn compile

CLASSPATH="`mvn -q exec:exec -Dexec.executable=echo -Dexec.args='%classpath'`:`cd target/classes && pwd`"

java -cp $CLASSPATH SimpleTest "org.postgresql.Driver" "jdbc:postgresql://localhost/test?loggerLevel=DEBUG" "root" "password" "SELECT 1;"