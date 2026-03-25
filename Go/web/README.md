# 2048 Go Web Application

A clean, browser-based 2048 game built with Go and Gin. The application serves a responsive single-page interface with keyboard controls, score tracking, an in-app rules dialog, and a server-backed exit action.

## Why Gin

This project uses `gin` because it is the lightest and most direct fit for a small web game:

- Fast to set up for a single-page application
- Minimal boilerplate compared with larger MVC-style frameworks
- Easy static file serving and simple route handling

## Features

- Playable 2048 board with arrow key and `WASD` controls
- Score and best-score tracking
- Responsive layout sized for laptop screens
- Info dialog explaining the rules
- Exit button that requests graceful server shutdown
- Lightweight tile movement animation before board updates

## Project Structure

```text
.
├── main.go
├── go.mod
├── go.sum
├── static/
│   ├── index.html
│   ├── styles.css
│   └── app.js
└── README.md
```

## Requirements

- Go `1.26.0` or newer

## Running the Application

From the project folder:

```bash
go run .
```

Then open:

```text
http://localhost:8080
```

## Controls

- `Up Arrow` or `W`: move up
- `Down Arrow` or `S`: move down
- `Left Arrow` or `A`: move left
- `Right Arrow` or `D`: move right

## How to Play

1. Move all tiles in one direction using the keyboard.
2. Matching tiles merge when they collide.
3. Each merge creates a tile with double the value.
4. A new tile appears after each valid move.
5. Reach `2048` to win.
6. The game ends when the board is full and no merges remain.

## Development Notes

- Static frontend files are embedded into the Go binary using `embed`.
- Best score is stored in browser `localStorage`.
- The exit button calls the `/exit` route, which shuts down the Go server gracefully.

## Build Check

To verify the project compiles:

```bash
go build ./...
```
