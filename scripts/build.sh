#!/bin/bash

if test -f ./cq_server; then
	rm ./cq_server
fi

go build .

