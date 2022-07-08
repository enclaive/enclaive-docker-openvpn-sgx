#!/bin/bash

docker-compose run -v "$HOME/.config/gramine/:/root/.config/gramine/" --entrypoint /manifest/manifest.sh ovpn
