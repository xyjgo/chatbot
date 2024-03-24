## Chatbot 

two binaries
> chatbot:  implements the main function: receive reviews, send reply according to rules defined, save to db  
> simserver:  simulates telegram bot api 

## Build & Run
### local dev
- install mongodb and run 
- local build chatbot and simserver
```shell
make
```
- local Run
```shell
cd cmd/chatbot/ && ./chatbot
cd cmd/simserver/ && ./simserver
```

### deploy with docker
- build chatbot and simserver images
- if these containers run on same server, replace ip address with container name in config.yml before running
```shell
docker build -f Dockerfile-Chatbot -t chatbot .
docker build -f Dockerfile-Simserver -t simserver .
mkdir mongo-data && chmod -R 777 mongo-data
docker-compose up -d
```