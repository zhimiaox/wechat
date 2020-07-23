#!/bin/sh

set -e

if [ "$1" = 'wechat_server_linux_amd64' -a "$(id -u)" = '0' ]; then
    exec su-exec zhimiao "$0" "$@"
fi

exec $@