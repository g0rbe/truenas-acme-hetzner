LDFLAGS = -s
LDFLAGS += -w
LDFLAGS += -extldflags "-static"
# CGO not used in this project, but leave this for future reference
LDFLAGS += -linkmode "external"

GOFLAGS = -tags "netgo,osusergo" -ldflags='$(LDFLAGS)'

build: clean
	go build -o release/tah-linux-amd64 $(GOFLAGS) .
	cd release && sha512sum * | gpg --local-user daniel@elmasy.com -o checksum.txt --clearsign

clean:
	@if [ -e "./release/tah-linux-amd64" ]; then rm -rf "./release/tah-linux-amd64"	; fi
	@if [ -e "./release/checksum.txt" ];	then rm -rf "./release/checksum.txt"	; fi