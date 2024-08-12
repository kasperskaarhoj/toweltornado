#!/bin/sh

# GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.License=Free" -o binaries/TowelTornado.Mac-arm64-m1
GOOS=darwin GOARCH=arm64 go build -o binaries/TowelTornado.Mac-arm64-m1
GOOS=windows GOARCH=amd64 go build -o binaries/TowelTornado.Win-amd64.exe
GOOS=windows GOARCH=386 go build -o binaries/TowelTornado.Win-386.exe
GOOS=linux GOARCH=amd64 go build -o binaries/TowelTornado.Linux-amd64
GOOS=linux GOARCH=386 go build -o binaries/TowelTornado.Linux-386

cd binaries

zip TowelTornado.Mac.zip TowelTornado.Mac-arm64-m1
zip TowelTornado.Win.zip TowelTornado.Win-amd64.exe TowelTornado.Win-386.exe
zip TowelTornado.Linux.zip TowelTornado.Linux-amd64 TowelTornado.Linux-386

rm TowelTornado.Win-amd64.exe TowelTornado.Win-386.exe TowelTornado.Linux-amd64 TowelTornado.Linux-386 TowelTornado.Mac-arm64-m1

cd ..
