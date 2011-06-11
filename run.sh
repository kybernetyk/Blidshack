#!/bin/bash

gomake && ./getblids.sh && ./blids && ./metadata.py
