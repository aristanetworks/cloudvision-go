#!/usr/bin/env bash

mkdir -p /persist/prometheus/prod
systemctl daemon-reload
systemctl enable prometheus-prod
