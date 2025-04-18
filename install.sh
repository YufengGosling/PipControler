#!/bin/bash

go build -o bin/ -ldflags="-s -w" ./cmd/...

echo export PATH=$PATH:/home/PipCtrler/bin/ > /home/.bashrc
