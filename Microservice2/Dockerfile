FROM golang

ENV AZURE_STORAGE_ACCOUNT ""
ENV AZURE_STORAGE_ACCESS_KEY ""
ENV AZURE_STORAGE_CONTAINER ""

COPY . /go
WORKDIR /go

RUN apt-get update && apt-get install gcc
RUN go get -u -d github.com/Azure/azure-sdk-for-go/...
RUN go get github.com/gorilla/mux
RUN go get -u github.com/jinzhu/gorm
RUN go get -u github.com/mattn/go-sqlite3

WORKDIR /go/src
CMD go run main.go