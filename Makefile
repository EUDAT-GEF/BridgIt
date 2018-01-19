pack:
	mkdir -p build
	mkdir -p build/def
	docker build -t bridgitcompile:linux .
	docker run --rm -v $(PWD)/build:/go/src/github.com/EUDAT-GEF/BridgIt/build bridgitcompile:linux
	cp def/config.json ./build/def
	tar -cvzf bridgit.tar.gz build/*
	rm -rf build
