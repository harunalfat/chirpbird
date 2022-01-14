FROM redis:6.2.6

RUN mkdir -p /redis

WORKDIR /redis

COPY sentinel.conf .
COPY entrypoint.sh /usr/local/bin/

RUN chown redis:redis /redis/* && \
    chmod +x /usr/local/bin/entrypoint.sh

CMD ["entrypoint.sh"]