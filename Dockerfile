FROM golang:1.23.5-alpine

RUN apk add --no-cache \
    nodejs \
    npm 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest

RUN npm install
RUN templ generate
RUN npx @tailwindcss/cli -i ./static/css/style.css -o ./static/css/tailwind.css

ENV GOINSECURE=true

RUN GOOS=linux GOARCH=amd64 go build -o ./cmd/main ./cmd

EXPOSE 8080

CMD ["./cmd/main"]