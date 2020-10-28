FROM alpine:3.12.1
LABEL Description="This image is used to host the Audio len service" \ 
      Maintainer="airenass@gmail.com" Vendor="Airenas Vaiciunas" Version="1.0" 

WORKDIR /app

COPY data/audio-len /app/
COPY config.yaml /app/

RUN apk add sox ffmpeg

ENTRYPOINT ["./audio-len"]
