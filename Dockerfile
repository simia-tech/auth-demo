FROM scratch

COPY target/auth-demo /bin/auth-demo

ENV LISTEN_NETWORK="tcp" \
    LISTEN_ADDRESS=":8000"

EXPOSE 8000

ENTRYPOINT ["/bin/auth-demo"]
