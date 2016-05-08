FROM golang:1.6-onbuild
MAINTAINER lennart@mooxmirror.io

EXPOSE 80
ENV BLOGREPO https://github.com/lnsp/bloggy-blueprint.git
ENV CERTFILE ""
ENV KEYFILE ""

CMD app -reset -blog personal-blog -repo $BLOGREPO -c $CERTFILE -k $KEYFILE
