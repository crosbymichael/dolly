#!/bin/bash

ssh root@162.243.132.79 supervisorctl stop web
ssh root@45.55.8.248 supervisorctl stop web
redis-cli -h 104.236.142.140 del nodes.frontend-1.fill
redis-cli -h 104.236.142.140 del nodes.frontend-2.fill
redis-cli -h 104.236.142.140 del nodes.frontend-3.fill
redis-cli -h 104.236.142.140 del nodes.frontend-1.avg
redis-cli -h 104.236.142.140 del nodes.frontend-2.avg
redis-cli -h 104.236.142.140 del nodes.frontend-3.avg

ssh root@45.55.8.248 supervisorctl start web
sleep 96
ssh root@162.243.132.79 supervisorctl start web
