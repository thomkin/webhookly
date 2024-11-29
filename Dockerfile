FROM docker.io/golang:1.23.3-alpine 
WORKDIR /app 

ENTRYPOINT ["go", "run", "."]