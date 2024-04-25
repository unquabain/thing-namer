#! /bin/zsh

go build -C .. -o render_test/wizard_bacon_test_subject

./wizard_bacon_test_subject &

sleep 1

curl http://localhost:9099/wizardbacon.go > wizardbacon/wizardbacon.go

if ! go run main.go ; then
  killall wizard_bacon_test_subject
  exit -1
fi

killall wizard_bacon_test_subject


