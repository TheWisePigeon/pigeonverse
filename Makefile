build:
	@go build -o app main.go

install-deps:
	@go mod download

start-tailwind-compilation:
	@npx tailwindcss -i ./assets/input.css -o ./assets/output.css --minify --watch

build-css:
	@npx tailwindcss -i ./assets/input.css -o ./assets/output.css --minify
