build:
	go build -o pgsw
install:
	mkdir -p /etc/pterowatch
	cp ./pgsw /usr/bin/pgsw
	cp -n data/pgsw.service /etc/systemd/system/
.DEFAULT: build