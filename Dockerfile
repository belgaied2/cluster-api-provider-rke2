FROM rke2-ubuntu:0.0.1
RUN mkdir -p /opt/rke2-artifacts
COPY files/rke2-images.linux-amd64.tar.zst /opt/rke2-artifacts/
COPY files/rke2.linux-amd64.tar.gz /opt/rke2-artifacts/
COPY files/sha256sum-amd64.txt /opt/rke2-artifacts/
WORKDIR /opt
COPY files/install.sh ./
CMD ["/lib/systemd/systemd"]
