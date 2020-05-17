FROM golang:latest

WORKDIR /usr/vid_trimmer
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./main.go

FROM alpine:latest as runner
RUN apk add  --no-cache ffmpeg
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /usr/vid_trimmer/app .
RUN mkdir file
EXPOSE $PORT
CMD ["./app"]
