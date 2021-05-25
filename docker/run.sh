#!/bin/bash

read_var() {
    VAR=$(grep $1 $2 | xargs)
    IFS="=" read -ra VAR <<< "$VAR"
    echo ${VAR[1]}
}

MF_QUEUE_SYSTEM=$(read_var MF_QUEUE_SYSTEM .env)

if [ $MF_QUEUE_SYSTEM = "rabbitmq" ]; then
    docker-compose -f rabbitmq.yml up
elif [ $MF_QUEUE_SYSTEM = "nats" ]; then
    docker-compose -f nats.yml up
else
    # The default queue system is NATS
    docker-compose -f docker-compose.yml up
fi
