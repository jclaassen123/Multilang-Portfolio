// Memory game JavaScript

// Track the first and second tiles clicked, lock board to prevent extra clicks, session ID and guess count
let firstTile = null;
let secondTile = null;
let lockBoard = false;
let sessionId = null;
let guesses = 0;

/**
 * Starts a new game with the selected board size.
 * @param size board size in format 'rowsxcols', e.g., '4x4'
 */
async function startGame(size) {
    const [rows, cols] = size.split("x").map(Number);

    // Call backend to start new game
    const res = await fetch(`/api/start?rows=${rows}&cols=${cols}`, { method: "POST" });
    const data = await res.json();

    // Initialize game state
    sessionId = data.sessionId;
    guesses = 0;
    lockBoard = false;
    firstTile = null;
    secondTile = null;
    document.getElementById("guesses").innerText = guesses;

    const board = document.getElementById("board");
    board.style.gridTemplateColumns = `repeat(${cols}, 100px)`;
    board.innerHTML = "";

    // Create tile elements
    data.tiles.forEach(tile => {
        const div = document.createElement("div");
        div.className = "tile";
        div.dataset.tileId = tile.id;
        div.dataset.imageId = tile.imageId;
        div.dataset.flipped = "false";

        div.addEventListener("click", () => flipTile(div));
        board.appendChild(div);
    });
}

/**
 * Handles a tile flip when clicked.
 * @param element the tile element being flipped
 */
async function flipTile(element) {
    if (lockBoard || element.dataset.flipped === "true") return;

    // Prevent selecting the same tile twice
    if (firstTile && element === firstTile) return;

    // Show tile image
    element.classList.add("flipped", "selected");
    element.style.backgroundImage = `url('/images/tile${element.dataset.imageId}.png')`;

    if (!firstTile) {
        firstTile = element;
        return;
    }

    // Second tile clicked, lock board
    secondTile = element;
    lockBoard = true;

    // Call backend to check match
    const res = await fetch(
        `/api/flip/${sessionId}/${element.dataset.tileId}?firstTileId=${firstTile.dataset.tileId}`,
        { method: "POST" }
    );
    const result = await res.json();

    guesses++;
    document.getElementById("guesses").innerText = guesses;

    if (result.match && result.complete) {
        // Matched and game complete
        [firstTile, secondTile].forEach(el => {
            el.dataset.flipped = "true";
            el.classList.remove("selected");
            el.classList.add("matched");
        });
        resetSelection();
        showWinPopup();
    } else if (result.match) {
        // Matched but game continues
        [firstTile, secondTile].forEach(el => {
            el.dataset.flipped = "true";
            el.classList.remove("selected");
            el.classList.add("matched");
        });
        resetSelection();
    } else if (result.miss) {
        // Not a match, flip back after delay
        setTimeout(() => {
            [firstTile, secondTile].forEach(el => {
                el.classList.remove("flipped", "selected");
                el.style.backgroundImage = "";
                el.dataset.flipped = "false";
            });
            resetSelection();
        }, 800);
    }
}

/**
 * Resets the currently selected tiles and unlocks the board.
 */
function resetSelection() {
    firstTile = null;
    secondTile = null;
    lockBoard = false;
}

/**
 * Resets the board and starts a new game with the currently selected size.
 */
function resetBoard() {
    const sizeSelect = document.getElementById("board-size");
    startGame(sizeSelect.value);
}

/**
 * Displays the win popup for 2 seconds.
 */
function showWinPopup() {
    const popup = document.getElementById("win-popup");
    popup.classList.add("show");
    setTimeout(() => popup.classList.remove("show"), 2000);
}

/**
 * Shows the instructions popup.
 */
function showInstructions() {
    const popup = document.getElementById("instructions-popup");
    popup.classList.add("show");
}

/**
 * Hides the instructions popup.
 */
function hideInstructions() {
    const popup = document.getElementById("instructions-popup");
    popup.classList.remove("show");
}
