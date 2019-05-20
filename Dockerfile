FROM debian:latest

RUN mkdir /app
WORKDIR /app
ADD /bin/main /app/main

CMD ["./main"]