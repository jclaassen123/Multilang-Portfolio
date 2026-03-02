/**
 * Initializes leaderboard functionality once the DOM is fully loaded.
 *
 * Attaches a change event listener to the board size selector and
 * loads the leaderboard for the default board size on page load.
 */
document.addEventListener("DOMContentLoaded", () => {
    const sizeSelect = document.getElementById("board-size");

    // Reload leaderboard when board size is changed
    sizeSelect.addEventListener("change", loadLeaderboard);

    // Load the leaderboard on initial page load
    loadLeaderboard();
});

/**
 * Fetches and renders the leaderboard for the selected board size.
 *
 * Retrieves the top scores from the backend `/api/leaderboard` endpoint
 * and updates the leaderboard table in the HTML.
 */
async function loadLeaderboard() {
    const size = document.getElementById("board-size").value;
    const [rows, cols] = size.split("x"); // Expect format like "4x4"

    try {
        // Fetch leaderboard data from backend
        const res = await fetch(`/api/leaderboard?rows=${rows}&cols=${cols}`);
        const data = await res.json();

        const tbody = document.querySelector("#leaderboard-table tbody");
        tbody.innerHTML = "";

        // Display a message if no scores exist for this board size
        if (data.length === 0) {
            tbody.innerHTML = `
                <tr>
                    <td colspan="4" class="no-results">No scores yet</td>
                </tr>
            `;
            return;
        }

        // Render each leaderboard entry as a table row
        data.forEach((entry, index) => {
            const row = `
                <tr>
                    <td>${getRankDisplay(index)}</td>
                    <td>${entry.username}</td>
                    <td>${entry.guesses}</td>
                    <td>${entry.completedAt}</td>
                </tr>
            `;
            tbody.innerHTML += row;
        });
    } catch (error) {
        console.error("Error loading leaderboard:", error);
    }
}

/**
 * Converts a leaderboard index into a display string with optional medal emojis.
 *
 * @param {number} index Zero-based index of the leaderboard entry
 * @returns {string} Rank display string, e.g., "🥇 1", "🥈 2", "🥉 3", "4"
 */
function getRankDisplay(index) {
    if (index === 0) return "🥇 1";
    if (index === 1) return "🥈 2";
    if (index === 2) return "🥉 3";
    return index + 1;
}