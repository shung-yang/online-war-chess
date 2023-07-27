FROM golang:1.19.4
COPY . /online_war_chess
WORKDIR /online_war_chess
RUN go mod tidy
CMD go run .