#!/bin/bash
# /app/web을 백그라운드로 실행
/app/web &

# cron 데몬을 포그라운드 모드로 실행하여 컨테이너가 계속 실행되도록 유지
cron -f
