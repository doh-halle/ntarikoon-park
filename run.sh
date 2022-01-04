#!/bin/bash

go build -o ntarikoonpark cmd/web/*.go
./ntarikoonpark -dbname=ntarikoon_park -dbuser=hallecraft -dbpassword=manager123