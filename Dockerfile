FROM alpine:3.16.2

RUN mkdir /app

RUN mkdir -p /app/dashboard
RUN mkdir -p /app/scripts
COPY ./statics /app/statics
COPY ./dashboard/build /app/dashboard/build
COPY ./cmd/db /app/db
#COPY ./bin/api-server /app/api-server
RUN chmod a+x /app/cmd

WORKDIR /app