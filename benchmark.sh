cd "${0%/*}"

set -x

if [ "$1" == "" ]; then
    echo "Usage: ./benchmark.sh (postgresql|mysql|mssql|oracle|sqlite etc)"
    exit 1
fi

export GOFLAGS=-mod=vendor

LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$1 DYLD_LIBRARY_PATH=$1 go test -v -bench=. -count=1 -tags="$1"