FROM enclaive/gramine-os:impish-1.3-425e8450 AS builder

# fix for sgx-dcap-pccs.postinst
RUN mkdir -p /etc/init/

RUN apt-get update \
    && apt-get install -y sgx-dcap-pccs \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /opt/intel/sgx-dcap-pccs

USER pccs

RUN npm install

# final build

FROM enclaive/gramine-os:impish-1.3-425e8450

RUN apt-get update \
    && apt-get install -y jq nodejs \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /opt/intel/sgx-dcap-pccs

COPY --from=builder /opt/intel/sgx-dcap-pccs/ /opt/intel/sgx-dcap-pccs/
COPY ./start.sh .

RUN groupadd -g 103 pccs \
    && useradd -M -d /opt/intel/sgx-dcap-pccs/ -u 102 -g pccs pccs \
    && chown pccs:pccs .

USER pccs

ENTRYPOINT [ "/opt/intel/sgx-dcap-pccs/start.sh" ]
