#!/bin/bash

if [ ! -f "./cc-test-reporter" ]; then
  curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > ./cc-test-reporter
  chmod +x ./cc-test-reporter
fi

./cc-test-reporter before-build
for pkg in $(go list ./... | grep -v vendor); do
    go test -coverprofile=$(echo $pkg | tr / -).cover $pkg
done
echo "mode: set" > c.out
grep -h -v "^mode:" ./*.cover >> c.out
rm -f *.cover

./cc-test-reporter after-build
