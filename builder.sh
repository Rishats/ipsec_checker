#!/bin/bash

go mod download

go build -o ipsec_checker

chmod +x ipsec_checker

sudo ./ipsec_checker