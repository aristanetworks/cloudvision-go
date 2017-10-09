#!/bin/sh

umask 066
exec >>/etc/haproxy/ssl/update.log 2>&1
date
set -ex
cd /etc/haproxy/ssl
ip=35.197.98.21
domains='jenkins gerrit'
for i in $domains; do
  dir=$i.corp.arista.io
  rm -rf $dir-prev
  if test -d $dir; then
    mv $dir $dir-prev
  fi
  scp -pri id_rsa_jenkins root@$ip:/etc/letsencrypt/live/$dir . || {
    rv=$?
    mv $dir-prev $dir
    exit $rv
  }
  cat $dir/fullchain.pem $dir/privkey.pem > $dir/haproxy-combined.pem
done

systemctl reload haproxy

# Ensure we have at least 20 days left on the cert before we exit successfully.
for i in $domains; do
  dir=$i.corp.arista.io
  openssl x509 -checkend $((20*86400)) -noout -in $dir/haproxy-combined.pem
done
