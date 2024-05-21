# PATH := $(PATH):$$(go env GOPATH)/bin
# PATH := $(PATH):$(HOME)/.pub-cache/bin
# PATH := $(PATH):$(pwd)/node_modules/.bin
# PROTO_ROOT_DIR = .

# all: gen

# _gendart:
# 	@echo "Generating Dart code..."
# 	@rm -rf lib/gen/dart
# 	@mkdir -p lib/gen/dart
# 	@protoc --dart_out=grpc:lib/gen/dart protos/**/*.proto google/protobuf/empty.proto

# _gengo:
# 	@echo "Generating Go code..."
# 	@rm -rf lib/gen/go
# 	@mkdir -p lib/gen/go
# 	# @protoc --go_out=. --go-grpc_out=. ./protos/**/*.proto
# 	@protoc --go_out=. --connect-go_out=. ./protos/**/*.proto
# 	@mv github.com/astronlabltd/promotebrain-protobuf/lib/gen/go lib/gen/

# _genphp:
# 	@echo "Generating PHP code..."
# 	@rm -rf lib/gen/php
# 	@mkdir -p lib/gen/php
# 	@protoc --php_out=lib/gen/php/ ./protos/**/*.proto

# _gents:
# 	@echo "Generating TS code..."
# 	@rm -rf lib/gen/ts
# 	@mkdir -p lib/gen/ts
# 	@protoc --plugin=./node_modules/.bin/protoc-gen-es --es_out lib/gen/ts/ --es_opt target=ts ./protos/**/*.proto
# 	@protoc --plugin=./node_modules/.bin/protoc-gen-connect-web --connect-web_out lib/gen/ts/ --connect-web_opt target=ts ./protos/**/*.proto

# _genjss:
# 	@echo "Generating JSON schema..."
# 	@rm -rf lib/gen/jsonschema
# 	@mkdir -p lib/gen/jsonschema
# 	@protoc --jsonschema_out=enforce_oneof,disallow_bigints_as_strings,all_fields_required:lib/gen/jsonschema ./protos/**/*.proto

# gen: _gengo _gendart _genphp _gents

proto-setup:
# steps:
# install protoc compiler and flutter sdk first
	@echo "setup go grpc"
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest
	@echo "setup PHP protobuf"
	composer require google/protobuf
	@echo "setup buf cli"
	pnpm i -g @bufbuild/buf

# @echo "install json schema plugin"
# go install github.com/chrusty/protoc-gen-jsonschema/cmd/protoc-gen-jsonschema@latest
# @echo "setup dart grpc"
# dart pub global activate protoc_plugin
# composer require twirp/quickstart
# rm -rf ./tmp
# docker run --rm -it -v $PWD/tmp:/tmp --platform=linux/amd64 curlimages/curl sh -c "curl -Ls https://git.io/twirphp | sh -s -- -b /tmp"
# @echo "setup TS protobuf"
# pnpm config set auto-install-peers true
# pnpm i -D @bufbuild/protoc-gen-connect-web @bufbuild/protoc-gen-es
# pnpm i @bufbuild/connect-web @bufbuild/protobuf

# release version in github main branch: make release tag=0.0.4
release: | buf
	@echo "Creating github release"
	@-git checkout main
	@-git add .
	@-git commit -am "released v$(tag)"
	@-git tag -a v$(tag) -m "v$(tag)"
	@-git push origin main
	@-git push origin v$(tag)

buf:
	pnpm i -g @bufbuild/buf
	# buf mod update ./protos
	buf dep update ./protos
	buf lint
	buf format -w
	rm -rf gen/
	# buf generate
	buf generate --include-imports
	
build: buf
	buf breaking --against '.git#branch=main'