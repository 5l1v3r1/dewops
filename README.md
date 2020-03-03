# ahtapod

Decentralized word processing system. 
Client/Server model with decentralized nodes.
Leader is elected using RAFT. Leader communicates with client. <br>

Work splitting and work result communication are done via RabbitMQ.

### Run
`$ docker pull rabbitmq:3.7-management-alpine` <br>
`$ ./compute/build.sh` <br>
`$ ./MQ/build.sh` <br>
`$ docker-compose -f docker-compose.yml up -d` <br>

### Stop
`$  docker-compose -f docker-compose.yml stop rabbitmq server1 server2 server3 mqsetup
` <br>
`$ ./clear.sh`

### Version: _0.0.1 Abant_
Input as a text document. <br>
Split the data across multiple nodes. <br>
Process data and return **word count** to the user.


### Inspiration
Inspired by Apache Spark and Apache Hadoop.