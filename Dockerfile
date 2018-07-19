FROM debian:stretch-slim

ADD ./build/bin /bin

ENTRYPOINT ["/bin/lookout"]
CMD [ "serve" ]
