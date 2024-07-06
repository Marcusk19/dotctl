FROM golang:1.21.12
RUN apt-get update && \
    apt-get install -y git neovim make
WORKDIR /app
COPY go.mod go.sum /
RUN go mod download
COPY . .
RUN make install
ENTRYPOINT sh

