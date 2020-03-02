# ahtapod

Decentralized word processing system. 
Client/Server model with decentralized nodes.
Leader is elected using RAFT. Leader communicates with client. <br>

Work splitting and work result communication are done via RabbitMQ.

### Run
Setup network and queues <br>
`$ docker network create testing` <br>
`$ docker pull rabbitmq:3.7-management-alpine` <br>
`$ ./MQ/build.sh` <br>

`$ ./compute/build1.sh` // create server image and run server instance <br> 
`$ ./compute/build2.sh` // run another server instance <br>
TODO: these scripts will be combined in a docker-compose.yml

### Version: _0.0.1 Abant_
Input as a text document. <br>
Split the data across multiple nodes. <br>
Process data and return **word count** to the user.


### Inspiration
Inspired by Apache Spark and Apache Hadoop.