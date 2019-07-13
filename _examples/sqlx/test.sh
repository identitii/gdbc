cd "${0%/*}"
DYLD_LIBRARY_PATH=../../postgresql:../../mssql:../../mysql:../../sqlite go test -v -count=1