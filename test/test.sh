#!/usr/bin/env bash

: <<'END'
	testSizeP := flag.Int("t", 200, "number of request")
	LimiterFlagP := flag.Bool("l", false, "Limiter Flag")
END

sudo ulimit -n 6049
sudo sysctl -w kern.ipc.somaxconn=1024



cd /Users/user/IdeaProjects/rate-limiter/test
go build rateTest.go

var=0

for req in 1 10 50 100 150
do
  diffSum=0
  for i in {1..100}
  do

  var1=0
  inst1=$(/Users/user/IdeaProjects/rate-limiter/test/rateTest "-t"=$req "-l"=false)
  var1=$(echo "$var1 + $inst1" | bc)

  sleep 3

  var2=0
  inst2=$(/Users/user/IdeaProjects/rate-limiter/test/rateTest "-t"=$req "-l"=true)
  var2=$(echo "$var2 + $inst2" | bc)

  diffVar=$(echo "$var2 - $var1" | bc)
  diffSum=$(echo "$diffSum + $diffVar" | bc)

  sleep 3
  done

  #printf "%.3f\n" "$diffSum/10.0"|bc -l
  printf "%.3f\n" $(bc -l <<< "$diffSum/10")

done

