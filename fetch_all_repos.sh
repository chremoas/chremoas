#!/usr/bin/env bash

REPOS="
role-srv
sig-cmd
esi-srv
auth-srv
role-cmd
auth-cmd
auth-web
purge-cmd
perms-srv
perms-cmd
lookup-cmd
discord-gateway
services-common
auth-esi-poller
chremoas-ctl
chremoas-helm
discord-tools
"

cd ..

for repo in $REPOS; do
  git clone "https://github.com/chremoas/${repo}.git"
done