FROM alpine:3.4
MAINTAINER elcuervo <elcuervo@elcuervo.net>

ENV GOSU_VERSION 1.9
ENV ENSURE_VERSION 0.1
RUN set -x \
      && apk add --no-cache --virtual .gosu-deps dpkg gnupg openssl \
      && dpkgArch="$(dpkg --print-architecture | awk -F- '{ print $NF }')" \
      && wget -O /usr/local/bin/gosu "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch"\
      && wget -O /usr/local/bin/gosu.asc "https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch.asc" \
      && wget -O /usr/local/bin/ensure "https://github.com/elcuervo/ensure/releases/download/$ENSURE_VERSION/ensure_linux_$dpkgArch" \
      && export GNUPGHOME="$(mktemp -d)" \
      && gpg --keyserver ha.pool.sks-keyservers.net --recv-keys B42F6819007F00F88E364FD4036A9C25BF357DD4 \
      && gpg --batch --verify /usr/local/bin/gosu.asc /usr/local/bin/gosu \
      && rm -r "$GNUPGHOME" /usr/local/bin/gosu.asc \
      && chmod +x /usr/local/bin/gosu \
      && chmod +x /usr/local/bin/ensure \
      && gosu nobody true \
      && apk del .gosu-deps

RUN addgroup ensure && adduser -S -G ensure ensure

ENTRYPOINT ["gosu", "ensure", "ensure"]
