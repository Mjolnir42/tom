FROM debian:bullseye-slim

WORKDIR /tom

RUN apt update && apt upgrade -y
