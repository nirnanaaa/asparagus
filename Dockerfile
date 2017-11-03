FROM busybox:ubuntu-14.04

MAINTAINER "Florian Kasper <florian@xpandmmi.com>"

EXPOSE 9092
WORKDIR /app
ENV PATH=/app:$PATH
COPY build/asparagus /app/
ENTRYPOINT ["asparagus"]
