FROM golang:latest

WORKDIR /usr/vid_trimmer
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./main.go

FROM alpine:latest as runner
RUN apk add  --no-cache ffmpeg
RUN apk --no-cache add ca-certificates curl python3
RUN ln -s /usr/bin/python3 /usr/bin/python
#/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
RUN curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl
RUN chmod a+rx /usr/local/bin/youtube-dl

WORKDIR /root/
COPY --from=0 /usr/vid_trimmer/app .
COPY --from=0 /usr/vid_trimmer/public ./public
RUN echo "let baseURL = \"${HOST}\"" > ./public/link.js
RUN mkdir file
EXPOSE $PORT
CMD ["./app"]
