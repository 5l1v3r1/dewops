docker build -t mq .
docker run -d --network testing --name mq mq