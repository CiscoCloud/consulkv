FROM gliderlabs/alpine

MAINTAINER Chris Aubuchon <Chris.Aubuchon@gmail.com>

COPY . /go/src/github.com/CiscoCloud/consulkv
RUN apk add --update go git mercurial \
	&& cd /go/src/github.com/CiscoCloud/consulkv \
	&& export GOPATH=/go \
	&& go get \
	&& go build -o /bin/consulkv \
	&& rm -rf /go \
	&& apk del --purge go git mercurual

ENTRYPOINT [ "/bin/consulkv" ]
