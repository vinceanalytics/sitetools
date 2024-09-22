
serve:
	 go run main.go -root=vince
eject:
	 go run main.go -root=vince -eject=out
release:
	 go run main.go -s -r  in out/
css:
	cd static && npm run css && cd -