FROM busybox:1.28.4-glibc

COPY build/main /bin/irishman

RUN chmod +x bin/irishman

CMD ["bin/irishman"]
