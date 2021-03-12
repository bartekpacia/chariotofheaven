#!/bin/bash

while true; do
	gpio -g write 20 1
	echo "wrote 1"
	sleep 0.5
	gpio -g write 20 0
	echo "wrote 0"
	sleep 0.5
done
