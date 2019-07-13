set -e
cd "${0%/*}"
curl -fLo graal.tar.gz https://glare.now.sh/oracle/graal/graalvm-ce-darwin
tar -zxf graal.tar.gz
`pwd`/graalvm-ce-19.1.0/Contents/Home/bin/gu install native-image
if ! [ -x "$(command -v mvn)" ]; then
  HOMEBREW_NO_AUTO_UPDATE=1 brew install maven
fi
