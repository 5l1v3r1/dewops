# ahtapod

Decentralized word processing system. 
Client/Server model with decentralized nodes.
Leader is elected using RAFT. Leader communicates with client. <br>

Work splitting and work result communication are done via RabbitMQ.

### Version: _0.0.1 Abant_
Input as a text document. <br>
Split the data across multiple nodes. <br>
Process data and return **word count** to the user.


### Inspiration
Inspired by Apache Spark and Apache Hadoop.