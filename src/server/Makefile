docker:
	docker build -t smvs/server:latest .

# !WARNING: DO NOT CALL THIS ONE OUTSIDE OF A DOCKER IMAGE
# This is to reduce image size as well as harden it
# Removes all the go dependencies and the package manager
cleanServer: server
	rm -f /sbin/apk; \
	rm -rf /etc/apk; \
	rm -rf /lib/apk; \
	rm -rf /usr/share/apk; \
	rm -rf /var/lib/apk; \
	rm -rf /go/bin; \
	rm -rf /go/pkg

server: cert
	go build . \

cert:
	cd certs; \
	./gen.sh

clean:
	rm server certs/*.pem certs/ca-cert.srl