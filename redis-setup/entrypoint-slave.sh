#!/bin/sh

redis-server /redis/sentinel-slave.conf --sentinel --slaveof redis1 6379