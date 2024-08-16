
serve:
	 go run main.go -s  in out/
release:
	 go run main.go -s -r  in out/
css:
	cd static && npm run css && cd -