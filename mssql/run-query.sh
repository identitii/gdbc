cd "${0%/*}"
./libgdbc-mssql-query "com.microsoft.sqlserver.jdbc.SQLServerDriver" "jdbc:sqlserver://localhost:1433;databaseName=test" "sa" "yourStrong(!)Password" "SELECT 1;"