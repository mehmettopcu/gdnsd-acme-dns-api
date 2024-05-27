FROM alpine:latest AS builder_gdnsd

ARG GDNS_VER=3.8.2

ENV GDNS_OPT="--prefix=/ --datarootdir=/usr --with-buildinfo=${GDNS_VER}-0"
ENV GDNS_BUILD_DEPENDENCY="perl perl-libwww ragel libev-dev autoconf automake libtool userspace-rcu-dev libcap-dev libmaxminddb-dev perl-test-harness perl-test-harness-utils libsodium-dev git"

RUN apk update \
&& apk add gcc g++ make patch file openssl ${GDNS_BUILD_DEPENDENCY} \
&& addgroup -S -g 101 gdnsd \
&& adduser -S -H -D -u 100 -s /sbin/nologin gdnsd gdnsd \
&& mkdir /usr/src \
&& cd /usr/src \
&& wget https://github.com/gdnsd/gdnsd/releases/download/v${GDNS_VER}/gdnsd-${GDNS_VER}.tar.xz \
&& tar xJf gdnsd-${GDNS_VER}.tar.xz

RUN cd /usr/src \
&& cd gdnsd-${GDNS_VER} \
&& autoreconf -vif \
&& ./configure ${GDNS_OPT} \
&& make \
&& make install

RUN find / -name gdnsd && ls -l /var/lib


# Build stage
FROM golang:alpine AS builder_go

COPY go_app /app

WORKDIR /app

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Build gdnsd-image
FROM alpine:latest as gdnsd

RUN apk --no-cache add libev libsodium libmaxminddb userspace-rcu libcap bind-tools

COPY --from=builder_gdnsd /etc/gdnsd            /etc/gdnsd
COPY --from=builder_gdnsd /libexec/gdnsd        /libexec/gdnsd
COPY --from=builder_gdnsd /usr/doc/gdnsd        /usr/doc/gdsnd
COPY --from=builder_gdnsd /var/lib/gdnsd        /var/lib/gdnsd
COPY --from=builder_gdnsd /sbin/gdnsd           /sbin/gdnsd
COPY --from=builder_gdnsd /bin/gdnsd_geoip_test /bin/gdnsd_geoip_test
COPY --from=builder_gdnsd /run/gdnsd            /run/gdnsd

RUN find / -name gdnsd

EXPOSE 53

CMD [ "/sbin/gdnsd", "start" ]

# Build api-image
FROM alpine:latest as api

RUN apk add --no-cache curl bind-tools

COPY --from=builder_gdnsd /bin/gdnsdctl         /bin/gdnsdctl
COPY --from=builder_go /app/app /bin/gdnsd_acme_api

EXPOSE 8080

CMD [ "/bin/gdnsd_acme_api"]

# Build test
FROM nginx:alpine as checkme

ARG VHOST=www.example.local

# Configure the Nginx settings
RUN echo "server {" > /etc/nginx/conf.d/default.conf
RUN echo "    listen 80;" >> /etc/nginx/conf.d/default.conf
RUN echo "    server_name ${VHOST};" >> /etc/nginx/conf.d/default.conf
RUN echo "    location /checkme {" >> /etc/nginx/conf.d/default.conf
RUN echo "        add_header Content-Type text/html;" >> /etc/nginx/conf.d/default.conf
RUN echo "        return 200 '\'\${HOSTNAME}\' > Health check passed!';" >> /etc/nginx/conf.d/default.conf
RUN echo "    }" >> /etc/nginx/conf.d/default.conf
RUN echo "}" >> /etc/nginx/conf.d/default.conf

