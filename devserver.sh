#! /bin/zsh

function start() {
  echo "starting"
  go build -o wbdevbuild
  killall wbdevbuild
  ./wbdevbuild &
  echo "started"
}

start

fswatch --event Updated -0 thing-namer.go templates/index.html | while read -d "" event; do echo $event changed; start; done

