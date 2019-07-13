#!/bin/sh

# set -x

# hacky script to test gitlab-ci jobs using docker

parse_yaml() {
   local prefix=$2
   local s='[[:space:]]*' w='[a-zA-Z0-9_]*' fs=$(echo @|tr @ '\034')
   sed -ne "s|^\($s\)\($w\)$s:$s\"\(.*\)\"$s\$|\1$fs\2$fs\3|p" \
        -e "s|^\($s\)\($w\)$s:$s\(.*\)$s\$|\1$fs\2$fs\3|p"  $1 |
   awk -F$fs '{
      indent = length($1)/2;
      vname[indent] = $2;
      for (i in vname) {if (i > indent) {delete vname[i]}}
      if (length($3) > 0) {
         vn=""; for (i=0; i<indent; i++) {vn=(vn)(vname[i])("_")}
         printf("%s%s%s=\"%s\"\n", "'$prefix'",vn, $2, $3);
      }
   }'
}

set -e
#set -x

eval $(parse_yaml .gitlab-ci.yml "config_")

SCRIPT="config_$1_script"
IMAGE="config_$1_image"
CACHE="/root/.m2/"

docker run -v "`pwd`/.dockercache:$CACHE" -v "`pwd`:/work" -w "/work" -it ${!IMAGE} bash -c "${!SCRIPT}"