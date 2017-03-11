OUTPUT=./_build
SRC=$(shell find . -iname "*.go")
LDFLAGS='-X main.pkgType="binary" -s -w'
RESOURCES=$(wildcard ./console/resources/*.html)

all: binaries msi deb

binaries: $(SRC)
	GOOS=darwin GOARCH=amd64 POSTFIX= make $(OUTPUT)/count-darwin-amd64
	GOOS=windows GOARCH=386 POSTFIX=.exe make $(OUTPUT)/count-windows-386.exe
	GOOS=windows GOARCH=amd64 POSTFIX=.exe make $(OUTPUT)/count-windows-amd64.exe
	GOOS=linux GOARCH=amd64 POSTFIX= make $(OUTPUT)/count-linux-amd64
	GOOS=linux GOARCH=386 POSTFIX= make $(OUTPUT)/count-linux-386

msi:
	wixl -a x86 packaging/msi/count-x86.wxs -o $(OUTPUT)/count-setup-x86.msi
	wixl -a x64 packaging/msi/count-x64.wxs -o $(OUTPUT)/count-setup-x64.msi

deb:
	mkdir -p $(OUTPUT)/x86-deb/DEBIAN/
	mkdir -p $(OUTPUT)/x86-deb/usr/bin/
	cp $(OUTPUT)/count-linux-386 $(OUTPUT)/x86-deb/usr/bin/count
	cp packaging/deb/control-x86 $(OUTPUT)/x86-deb/DEBIAN/control
	dpkg-deb --build $(OUTPUT)/x86-deb $(OUTPUT)/count-x86.deb
	mkdir -p $(OUTPUT)/x64-deb/DEBIAN/
	mkdir -p $(OUTPUT)/x64-deb/usr/bin/
	cp $(OUTPUT)/count-linux-amd64 $(OUTPUT)/x64-deb/usr/bin/count
	cp packaging/deb/control-x64 $(OUTPUT)/x64-deb/DEBIAN/control
	dpkg-deb --build $(OUTPUT)/x64-deb $(OUTPUT)/count-x64.deb

$(OUTPUT)/count-$(GOOS)-$(GOARCH)$(POSTFIX): $(SRC)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $@ -ldflags=$(LDFLAGS) github.com/marekchen/count

install: resources
	GOOS=$(GOOS) go install github.com/marekchen/count


resources:
	(cd console; $(MAKE))

ccount:
	rm -rf $(OUTPUT)

.PHONY: test msi deb install ccount resources
