make:
	go build -o build/lib tlds.go urlextract.go

lint:
	go fmt .
