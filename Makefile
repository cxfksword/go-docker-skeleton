build: fontend-build
	go build -o ./tmp/main .


fontend-build:
	cd view && npm run build && cd ..