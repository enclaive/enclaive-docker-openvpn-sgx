FROM enclaive/gramine-os:impish-1.3-425e8450

RUN apt-get update \
    && apt-get install -y git nano build-essential golang

RUN git clone --branch=master https://github.com/gramineproject/gramine \
    && cd ./gramine/ \
    && git checkout 425e8450

COPY ./src/ ./provisioner.sh /

RUN CGO_CFLAGS_ALLOW="-L/usr/local/lib/x86_64-linux-gnu" go build -buildmode=pie

EXPOSE 4433/tcp

ENTRYPOINT [ "/provisioner.sh" ]
