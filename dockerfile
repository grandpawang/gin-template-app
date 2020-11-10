FROM alpine:latest
ENV TZ=Asia/Shanghai
EXPOSE 50052 
WORKDIR /root/
COPY ./cmd/cloud ./cmd/favicon.ico ./cmd/config.toml ./
# CMD ["/bin/sh"]
# make docker && docker run -d -p 50052:50052 --name=gbbmn-cloud gbbmn-cloud
CMD ["./cloud", "--conf", "/root/config.toml", "--icon", "/root/favicon.ico"]

