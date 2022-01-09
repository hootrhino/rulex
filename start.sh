#!/bin/bash

##
## 启动脚本
##
nohup ./rulex run &>run_log_$(date '+%Y_%m_%d_%H_%M_%S').log &
