#!/usr/bin/env bash

mkdir -p /persist/prometheus/dev
systemctl daemon-reload
systemctl enable prometheus-dev
