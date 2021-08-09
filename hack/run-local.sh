export SERVER_PORT=31199
export LOG_LEVEL=trace
export KAFKA_BROKERS=192.168.1.3:9092
export TOPIC=iotpaas
export REDIS_HOST=192.168.1.21:6379
export REDIS_PASSWORD=test
export VERSION=1.0.3
export NAME=iotpass-message-producer

./build/microservice
