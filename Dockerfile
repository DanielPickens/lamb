FROM alpine:3.16.2

RUN mkdir /app



#COPY ./bin/cmd/main.go /app/cmd/main.go
RUN chmod a+x /app/cmd

WORKDIR /app