cd "${0%/*}"
./libgdbc-sqlite-query "org.sqlite.JDBC" "jdbc:sqlite::memory:" "" "" "SELECT 1;"