FROM busybox:1.28.4-glibc

COPY build/${CI_PROJECT_NAME} /bin/${CI_PROJECT_NAME}

RUN chmod +x bin/${CI_PROJECT_NAME}

CMD ["bin/${CI_PROJECT_NAME}"]
