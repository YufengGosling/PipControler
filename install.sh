#!/bin/bash

go build -o bin/ -ldflags="-s -w" ./cmd/...

cp -r scripts bin/scripts

echo export PATH=$PATH:/home/PipControler/bin/ > /home/.bashrc
