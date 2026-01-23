// ---------------------------
// Section References
// ---------------------------

// Grab main sections of the app
const loginSection = document.getElementById('login-section');
const timeSelection = document.getElementById('time-selection');
const gameSection = document.getElementById('game-section');
const leaderboardSection = document.getElementById('leaderboard-section');

// Global variables
let username = '';
let chosenTime = 0;
let counting = false;
let startTime = 0;

// Dropdown for time selection
const chosenTimeSelect = document.getElementById('chosen-time');


// ---------------------------
// Utility Functions
// ---------------------------

/**
 * Format seconds into a human-readable string.
 * <60s: show "xx.xx sec"
 * >=60s: show "m:ss min"
 */
function formatTimeDisplay(seconds) {
    if (seconds < 60) {
        return seconds + ' sec';
    } else {
        let mins = Math.floor(seconds / 60);
        let secs = seconds % 60;
        return mins + ':' + (secs < 10 ? '0' + secs : secs) + ' min';
    }
}

/**
 * Format target time for leaderboard / game display
 * - Seconds < 60: show 2 decimal places
 * - Seconds >= 60: round to nearest second
 */
function formatTargetTime(seconds) {
    seconds = parseFloat(seconds);
    if (seconds < 60) {
        return seconds.toFixed(2) + ' sec';
    } else {
        let mins = Math.floor(seconds / 60);
        let secs = Math.round(seconds % 60);
        return mins + ':' + (secs < 10 ? '0' : '') + secs + ' min';
    }
}

/**
 * Format distance away for leaderboard display
 * - Seconds < 60: show 2 decimal places
 * - Seconds >= 60: round to nearest second
 */
function formatDistance(seconds) {
    seconds = parseFloat(seconds);
    if (seconds < 60) {
        return seconds.toFixed(2) + ' sec';
    } else {
        let mins = Math.floor(seconds / 60);
        let secs = Math.round(seconds % 60);
        return mins + ':' + (secs < 10 ? '0' : '') + secs + ' min';
    }
}


// ---------------------------
// Initialize Time Options
// ---------------------------

// 10s-60s in 10s increments
for (let t = 10; t <= 60; t += 10) {
    let opt = document.createElement('option');
    opt.value = t;
    opt.text = formatTimeDisplay(t);
    chosenTimeSelect.add(opt);
}

// 90s-600s in 30s increments
for (let t = 90; t <= 600; t += 30) {
    let opt = document.createElement('option');
    opt.value = t;
    opt.text = formatTimeDisplay(t);
    chosenTimeSelect.add(opt);
}


// ---------------------------
// Event Functions
// ---------------------------

/**
 * Transition from login section to time selection
 */
function startTimeSelection() {
    username = document.getElementById('username').value.trim();
    if (username === '') {
        alert('Enter a username');
        return;
    }
    loginSection.classList.add('hidden');
    timeSelection.classList.remove('hidden');
}

/**
 * Start the game
 * - Set target time
 * - Transition sections
 * - Attach start/stop button logic
 */
function startGame() {
    chosenTime = parseFloat(chosenTimeSelect.value);
    document.getElementById('target-time').textContent = formatTargetTime(chosenTime);

    timeSelection.classList.add('hidden');
    gameSection.classList.remove('hidden');

    const btn = document.getElementById('start-btn');
    const status = document.getElementById('status');

    btn.disabled = false;
    btn.textContent = "Start";
    status.textContent = "Click start to begin";

    btn.onclick = function () {
        if (!counting) {
            // Start counting
            counting = true;
            startTime = Date.now();
            status.textContent = "Counting...";
            btn.textContent = "Stop";
        } else {
            // Stop counting
            counting = false;
            let elapsed = ((Date.now() - startTime) / 1000).toFixed(2);
            status.textContent = `You stopped at: ${elapsed} sec`;
            btn.disabled = true;

            let distance = Math.abs(chosenTime - elapsed);

            // Send score to backend and show leaderboard
            setTimeout(() => {
                fetch('/submit_score', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username, chosen_time: chosenTime, distance })
                }).then(() => showLeaderboard());
            }, 1000);
        }
    }
}

/**
 * Show leaderboard and fetch top 5 scores
 */
function showLeaderboard() {
    gameSection.classList.add('hidden');
    leaderboardSection.classList.remove('hidden');

    fetch('/leaderboard')
        .then(res => res.json())
        .then(data => {
            const table = document.getElementById('leaderboard-table');
            table.innerHTML = `<tr>
                <th>Username</th>
                <th>Target Time</th>
                <th>Distance Away</th>
            </tr>`;

            // Display only top 5
            data.slice(0, 5).forEach(row => {
                const tr = document.createElement('tr');
                tr.innerHTML = `<td>${row.username}</td>
                                <td>${formatTargetTime(row.chosen_time)}</td>
                                <td>${formatDistance(row.distance)}</td>`;
                table.appendChild(tr);
            });
        });
}

/**
 * Finish the game and reset to login
 */
function finishGame() {
    leaderboardSection.classList.add('hidden');
    loginSection.classList.remove('hidden');

    document.getElementById('username').value = '';
    chosenTimeSelect.selectedIndex = 0;
    username = '';
    chosenTime = 0;
}

/**
 * Restart game from leaderboard to time selection
 */
function restartGame() {
    leaderboardSection.classList.add('hidden');
    timeSelection.classList.remove('hidden');
}
