FROM golang:1.18-alpine

ENV HOST=0.0.0.0
ENV PORT=8000

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go generate
RUN go build .

EXPOSE 8000

CMD [ "./bot" ]
