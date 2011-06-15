#!/bin/bash

gomake && ./getblids.sh && ./blids && ./metadata.py && mv aktkartepolengrau.jpg_metadata.txt pl_metadata.txt && mv aktkartegergrau.jpg_metadata.txt de_metadata.txt
