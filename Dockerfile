FROM golang:1.23.2

WORKDIR /go/src

RUN ln -sf /bin/bash /bin/sh

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

#manter container rodando
CMD ["tail", "-f", "/dev/null"] 
