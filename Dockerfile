FROM golang:1.6-onbuild
MAINTAINER lennart@mooxmirror.io

EXPOSE 80
ENV repo https://github.com/lnsp/bloggy-blueprint.git

CMD app -reset -blog personal-blog -repo $repo