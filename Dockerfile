# build images
FROM golang:1.13.8-alpine AS build-env
ADD . /home/irishman
WORKDIR /home/irishman
RUN go build -v -o /home/irishman/irishman

# Set one or more individual labels
LABEL irishman.version="v1.0.1"
LABEL vendor="Ran Jack"
LABEL irishman.release-date="2020.04.09"

# run images
FROM alpine
RUN apk add -U tzdata
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime
COPY --from=build-env /home/irishman/irishman  /usr/local/bin/irishman
COPY --from=build-env /home/irishman/config.yaml  /config.yaml

EXPOSE 8080
CMD ["irishman"]
