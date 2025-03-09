#!/bin/bash

go build -ldflags="-s -w" ./cmd/ipp
go build -ldflags="-s -w" ./cmd/pipcontroler

echo export PATH=$PATH:/home/PipControler > /home/.bashrc
