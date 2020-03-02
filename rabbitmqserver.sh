docker run -d --network testing --hostname rabbitmqhost \
   --name rabbitmqserver -p 15672:15672 -p 5672:5672 rabbitmq:3.7-management-alpine