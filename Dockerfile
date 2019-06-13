FROM golang:1.12
EXPOSE 80
COPY ./bin/service /usr/local/bin/
CMD ["service"]
