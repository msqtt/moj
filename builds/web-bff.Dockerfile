FROM golang:1.22-alpine as build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /target/web-bff moj/web-bff
RUN cp apps/web-bff/app.env /target

FROM alpine:3.20
WORKDIR /app
COPY --from=build /target/ .
EXPOSE 8080
CMD ["/app/web-bff"]