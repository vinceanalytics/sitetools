
serve:
	 go run main.go
release:
	 go run main.go -s -r  in out/
css:
	cd static && npm run css && cd -