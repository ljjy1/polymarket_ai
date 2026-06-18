FROM alpine:latest
MAINTAINER zhufuyi "g.zhufuyi@gmail.com"

# add curl, used for http service checking, can be installed without it if deployed in k8s
RUN apk add curl

COPY configs/ /app/configs/
COPY be /app/be
RUN chmod +x /app/be

# http port
EXPOSE 8080


WORKDIR /app

CMD ["./be", "-c", "configs/be.yml"]
# if you use the Configuration Center, be.yml is changed to the Configuration Center configuration.
#CMD ["./be", "-c", "configs/be.yml", "-enable-cc"]
