#!/bin/bash
while read p; do
  curl -sfLO https://github.com/rancher/rke2/releases/download/v1.26.0%2Brke2r1/$p
done <artifact-list.txt
