"""
Time Guess Game Flask App
Author: Jace Claassen
Description:
    A simple web game where users try to guess a target time.
    Users log in with a unique username, select a time interval, start/stop
    the counter, and see how close they were. Scores are saved to CSV and
    the top scores are displayed on a leaderboard.

    How to Run:
    Terminal: python3 app.py
    Browser: http://127.0.0.1:5000
"""

from flask import Flask, render_template, request, redirect, url_for, session, jsonify
import csv
import os

# ---------------------------
# App Initialization
# ---------------------------
app = Flask(__name__)
app.secret_key = 'supersecretkey'  # Needed for session handling

CSV_FILE = 'game_data.csv'

# ---------------------------
# CSV Initialization
# ---------------------------
def init_csv():
    """
    Initialize the CSV file with headers if it doesn't exist
    or if headers are missing/malformed.
    """
    needs_header = False

    if not os.path.exists(CSV_FILE):
        needs_header = True
    else:
        with open(CSV_FILE, 'r', newline='') as file:
            first_line = file.readline().strip()
            if first_line != 'username,chosen_time,distance':
                needs_header = True

    if needs_header:
        with open(CSV_FILE, 'w', newline='') as file:
            writer = csv.writer(file)
            writer.writerow(['username', 'chosen_time', 'distance'])

init_csv()

# ---------------------------
# Routes
# ---------------------------

@app.route("/", methods=['GET', 'POST'])
def login():
    """
    Handle user login.
    GET: display login page.
    POST: store username in session and redirect to time selection.
    """
    if request.method == 'POST':
        username = request.form.get('username')
        if username:
            session['username'] = username
            return redirect(url_for('select_time'))
    return render_template('login.html')


@app.route("/select_time", methods=['GET', 'POST'])
def select_time():
    """
    Let the user select the game length.
    Valid time options:
        - 10s to 1 min in 10s increments
        - 1 min to 10 min in 30s increments
    """
    if 'username' not in session:
        return redirect(url_for('login'))

    times_seconds = list(range(10, 61, 10)) + list(range(90, 601, 30))

    if request.method == 'POST':
        chosen_time = float(request.form.get('time'))
        session['chosen_time'] = chosen_time
        return redirect(url_for('game'))

    return render_template('select_time.html', times=times_seconds)


@app.route("/game", methods=['GET', 'POST'])
def game():
    """
    Game page where user starts/stops the timer and guesses the duration.
    POST: calculate distance from target, save to CSV, redirect to leaderboard.
    """
    if 'username' not in session or 'chosen_time' not in session:
        return redirect(url_for('login'))

    if request.method == 'POST':
        actual_time = float(request.form.get('actual_time'))
        chosen_time = session['chosen_time']
        distance = abs(chosen_time - actual_time)

        # Append result to CSV
        with open(CSV_FILE, 'a', newline='') as file:
            writer = csv.writer(file)
            writer.writerow([session['username'], chosen_time, round(distance, 2)])

        return redirect(url_for('leaderboard'))

    return render_template('game.html', chosen_time=session['chosen_time'])


@app.route("/leaderboard")
def leaderboard():
    """
    Display the top 5 scores from the CSV file sorted by distance (closest guesses first).
    """
    results = []
    if os.path.exists(CSV_FILE):
        with open(CSV_FILE, 'r') as file:
            reader = csv.DictReader(file)
            for row in reader:
                try:
                    row['distance'] = float(row['distance'])
                    row['chosen_time'] = float(row['chosen_time'])
                    results.append(row)
                except (KeyError, ValueError):
                    continue  # Skip malformed rows

    results = sorted(results, key=lambda x: x['distance'])[:5]
    return render_template('leaderboard.html', results=results)


# ---------------------------
# Run App
# ---------------------------
if __name__ == "__main__":
    app.run(debug=True)
