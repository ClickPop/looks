build:
	@for os in "darwin" "linux" "windows" ; do \
		for arch in "arm64" "amd64" ; do \
			echo "Building binary for $$os on $$arch" ; \
			if [ $$os = "windows" ] ; then \
				env GOOS=$$os GOARCH=$$arch go build -o bin/looks-$$os-$$arch.exe github.com/clickpop/looks ; \
			else \
				env GOOS=$$os GOARCH=$$arch go build -o bin/looks-$$os-$$arch github.com/clickpop/looks ; \
			fi \
		done \
	done