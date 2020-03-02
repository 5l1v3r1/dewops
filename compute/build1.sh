docker build -t server .
docker run -p 8080:8081 --network="testing" --name server1 -it server