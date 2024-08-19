FROM golang:1.22-alpine as build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /target/game moj/game
RUN cp apps/game/game.env /target

FROM alpine:3.20
WORKDIR /app
COPY --from=build /target/ .
EXPOSE 9090
CMD ["/app/game"]