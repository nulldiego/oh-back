FROM golang:1.22 as build
WORKDIR /app
COPY . .
RUN go build -o /oh-back .

FROM scratch
COPY --from=build /oh-back /oh-back
EXPOSE 8080
CMD ["/oh-back"]