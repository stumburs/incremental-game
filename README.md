# incremental-game

This is a simple COOP incremental game about programming. It uses Firebase's realtime database as a _server_ and _storage_. It was created as a fun experiment between me and my friend.

## Quick Start

To build the game you need Go 1.24.1.

```sh
go build src/incremental-game.go
```

### Database

To run the game you also need to create a Firebase project and set up the database. I cannot help you with this as I'm using my friend's Firebase project.

However, you do need a .env file with the following variables in the root directory of the project:

```env
APIKey = "YOUR REST API KEY"
DatabaseURL = "YOUR DATABASE URL"
```

## Features

Just click the "Compile" button to compile some code and get "Compile Points" which you can use to buy upgrades.
![Screenshot from 2025-04-15 00-47-29](https://github.com/user-attachments/assets/d543f9b9-2610-4517-aa88-108e3eb5b38a)
