# Time Guess Game (Flask Web App)

Author: Jace Claassen

## Overview

Time Guess Game is a simple web-based game built with Python and Flask where users try to estimate how much time has passed. Players log in with a username, select a target time interval, start and stop a timer, and then see how close their guess was to the selected duration.

Game results are stored in a CSV file, and the best scores are displayed on a leaderboard showing the most accurate guesses.

---

## Features

- Web-based interface using Flask and HTML templates
- Simple username-based login using Flask sessions
- Time interval selection with predefined options
- Start/stop timing gameplay mechanic
- Accuracy calculation based on distance from target time
- Persistent score storage using a CSV file
- Leaderboard displaying the top 5 closest guesses
- Server-side validation and session handling

---

## Requirements

- Python 3.x
- Flask

Install Flask if needed:

pip install flask

No database setup is required. Scores are stored in a CSV file.

---

## Project Structure

- app.py — Main Flask application
- templates/
  - login.html — User login page
  - select_time.html — Time selection page
  - game.html — Gameplay page
  - leaderboard.html — Leaderboard display
- game_data.csv — Stores usernames, chosen times, and accuracy scores

---

## How to Run

1. Open a terminal and navigate to the project directory.
2. Start the Flask application:

   python3 app.py

3. Open a web browser and go to:

   http://127.0.0.1:5000

---

## How to Play

1. Enter a username to log in.
2. Select a target time interval.
3. Start the timer and stop it when you believe the selected amount of time has passed.
4. Submit your guess.
5. View your result and see how close you were.
6. Check the leaderboard to see the top scores.

---

## Game Rules

- Time options range from:
  - 10 seconds to 1 minute (10-second increments)
  - 1 minute to 10 minutes (30-second increments)
- Your score is calculated as the absolute difference between the chosen time and the actual elapsed time.
- Smaller distances indicate better performance.
- The leaderboard shows the top 5 closest guesses.

---

## Design Notes

- Flask sessions are used to track logged-in users and selected game times.
- CSV storage is used instead of a database for simplicity and transparency.
- The CSV file is automatically initialized with headers if missing or malformed.
- Leaderboard data is sorted server-side to ensure accurate ranking.
- The project emphasizes core web concepts: routing, sessions, templates, and persistence.
