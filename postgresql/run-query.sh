cd "${0%/*}"
#sudo dtruss 
./libgdbc-postgresql-query "org.postgresql.Driver" "jdbc:postgresql://127.0.0.1:5432/test?loggerLevel=DEBUG" "root" "password" "SELECT 1;"