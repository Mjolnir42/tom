FROM debian:bullseye-slim

ADD . /tom
WORKDIR /tom

RUN apt update\
	&& apt upgrade -y\
	&& apt install -y\
		golang ca-certificates\
		git build-essential

RUN go mod download &&\
	make install_all
