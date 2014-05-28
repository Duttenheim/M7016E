#!/bin/bash
scp -i $1 ./code/bin/bitverseserver ubuntu@130.240.134.117:
scp -i $1 ./ping.bmp ubuntu@130.240.134.117:
scp -i $1 ./code/bin/bitverseserver ubuntu@130.240.134.120:
scp -i $1 ./ping.bmp ubuntu@130.240.134.120:
