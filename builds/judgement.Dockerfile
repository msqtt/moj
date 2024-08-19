FROM golang:1.22-alpine as build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /target/judgement moj/judgement
RUN cp apps/judgement/judgement.env /target

FROM alpine:3.20
WORKDIR /app
COPY --from=build /target/ .
EXPOSE 9090
CMD ["/app/judgement"]