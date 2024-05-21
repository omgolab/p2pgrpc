build:
	npm i -g pnpm
	rm -rf dist-bundle
	mkdir -p dist-bundle/
	${MAKE} -j fe_build # go_build be_build # building parallelly
	cd dist-bundle && \
	zip -r dist.zip .

fe_build:
	cd frontend && \
	pnpm build && \
	cp -r dist/. ../dist-bundle/

go_build:
	cd backend && \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o app.bin && \
	cp app.bin ../dist-bundle/ 

be_build:
	cd backend && \
	pnpm build && \
	mkdir -p ../dist-bundle/backend && \
	cp -r dist/. ../dist-bundle/backend/
