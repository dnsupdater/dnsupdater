FROM alpine:3.18
COPY bin/dnsu /dnsu/dnsu
COPY update-dns.sh /dnsu/update-dns.sh
# ENTRYPOINT ["tail"]
# CMD ["-f","/dev/null"]
