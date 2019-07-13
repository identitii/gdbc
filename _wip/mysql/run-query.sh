cd "${0%/*}"
./libgdbc-mysql-query "com.mysql.jdbc.Driver" "jdbc:mysql://localhost/test" "root" "password" "SELECT 1;"