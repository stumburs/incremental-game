# Temporary solution that will most likely stay permanent

VERSION=1.0.0

linux:
	go build -ldflags="-s -w" src/incremental-game.go
	strip incremental-game
	upx incremental-game
windows:
	env CGO_ENABLED=1 GOOS=windows CC=x86_64-w64-mingw32-gcc go build -ldflags="-s -w" src/incremental-game.go
	strip incremental-game.exe
	upx incremental-game.exe

all: linux windows

tar: all
	tar -czvf incremental-game_$(VERSION)_linux_amd64.tar.gz incremental-game
	tar -czvf incremental-game_$(VERSION)_windows_amd64.tar.gz incremental-game.exe

archive: all
	tar -czvf incremental-game_$(VERSION)_linux_amd64.tar.gz incremental-game
	7z a -mx9 incremental-game_$(VERSION)_windows_amd64.7z incremental-game.exe

# include .env file
archive-env: all
	tar -czvf incremental-game_$(VERSION)_linux_amd64.tar.gz incremental-game .env
	7z a -mx9 incremental-game_$(VERSION)_windows_amd64.7z incremental-game.exe .env

# create a 7z archive for windows split into 10mb parts (for sending through discord)
win-discord-env: windows
	7z -mx9 -v10m a incremental-game_$(VERSION)_windows_amd64.7z incremental-game.exe .env

clean:
	rm incremental-game*