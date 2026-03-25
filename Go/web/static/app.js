const size = 4;
const boardElement = document.getElementById("board");
const scoreElement = document.getElementById("score");
const bestScoreElement = document.getElementById("best-score");
const newGameButton = document.getElementById("new-game-button");
const infoButton = document.getElementById("info-button");
const exitButton = document.getElementById("exit-button");
const overlay = document.getElementById("overlay");
const overlayTitle = document.getElementById("overlay-title");
const overlayMessage = document.getElementById("overlay-message");
const overlayButton = document.getElementById("overlay-button");
const infoModal = document.getElementById("info-modal");
const closeModalButton = document.getElementById("close-modal-button");

let board = [];
let score = 0;
let bestScore = Number(localStorage.getItem("best-2048-score") || 0);
let gameWon = false;
let isAnimating = false;
const moveAnimationMs = 110;

bestScoreElement.textContent = bestScore;

/**
 * Resets all gameplay state and starts a fresh 2048 session.
 */
function startGame() {
  board = Array.from({ length: size }, () => Array(size).fill(0));
  score = 0;
  gameWon = false;
  isAnimating = false;
  updateScore();
  hideOverlay();
  addRandomTile();
  addRandomTile();
  renderBoard();
}

/**
 * Updates the visible score and persists the best score in local storage.
 */
function updateScore() {
  scoreElement.textContent = score;
  if (score > bestScore) {
    bestScore = score;
    localStorage.setItem("best-2048-score", String(bestScore));
  }
  bestScoreElement.textContent = bestScore;
}

/**
 * Places a new tile into a random empty position on the board.
 */
function addRandomTile() {
  const emptyCells = [];

  for (let row = 0; row < size; row += 1) {
    for (let col = 0; col < size; col += 1) {
      if (board[row][col] === 0) {
        emptyCells.push({ row, col });
      }
    }
  }

  if (emptyCells.length === 0) {
    return;
  }

  const { row, col } = emptyCells[Math.floor(Math.random() * emptyCells.length)];
  // A 2 appears most of the time to match standard 2048 behavior.
  board[row][col] = Math.random() < 0.9 ? 2 : 4;
}

/**
 * Rebuilds the visual grid from the current board state.
 */
function renderBoard() {
  boardElement.innerHTML = "";

  board.forEach((row) => {
    row.forEach((value) => {
      const cell = document.createElement("div");
      cell.className = "cell";

      if (value !== 0) {
        const tile = document.createElement("div");
        tile.className = `tile value-${value}${value > 2048 ? " super" : ""}`;
        if (value >= 1024) {
          tile.classList.add("shrink-text");
        }
        tile.textContent = value;
        cell.appendChild(tile);
      }

      boardElement.appendChild(cell);
    });
  });
}

/**
 * Creates a deep copy of the board so move calculations do not mutate the
 * current state until the animation has completed.
 *
 * @param {number[][]} sourceBoard The board to copy.
 * @returns {number[][]} A cloned 4x4 board.
 */
function cloneBoard(sourceBoard) {
  return sourceBoard.map((row) => [...row]);
}

/**
 * Compresses one row or column toward the move direction and merges matching
 * adjacent values exactly once per move.
 *
 * @param {number[]} line The row or column being processed.
 * @returns {{line: number[], lineScore: number}} The merged line and score gain.
 */
function slideAndMerge(line) {
  const compact = line.filter((value) => value !== 0);
  const merged = [];
  let lineScore = 0;

  for (let i = 0; i < compact.length; i += 1) {
    if (compact[i] !== 0 && compact[i] === compact[i + 1]) {
      const value = compact[i] * 2;
      merged.push(value);
      lineScore += value;
      i += 1;
    } else {
      merged.push(compact[i]);
    }
  }

  while (merged.length < size) {
    merged.push(0);
  }

  return { line: merged, lineScore };
}

/**
 * Computes the result of a move without immediately updating the rendered
 * board, which lets the UI play a short directional transition first.
 *
 * @param {number[][]} currentBoard The current game board.
 * @param {"left"|"right"|"up"|"down"} direction The requested move direction.
 * @returns {{moved: boolean, nextBoard: number[][], gained: number, won: boolean}}
 * The updated board, score gain, and movement metadata.
 */
function calculateMove(currentBoard, direction) {
  const nextBoard = cloneBoard(currentBoard);
  let gained = 0;
  let won = false;

  for (let index = 0; index < size; index += 1) {
    let line = [];

    if (direction === "left" || direction === "right") {
      line = [...nextBoard[index]];
      if (direction === "right") {
        line.reverse();
      }
    } else {
      for (let row = 0; row < size; row += 1) {
        line.push(nextBoard[row][index]);
      }
      if (direction === "down") {
        line.reverse();
      }
    }

    const result = slideAndMerge(line);
    let updatedLine = result.line;
    gained += result.lineScore;
    won = won || updatedLine.includes(2048);

    if (direction === "right" || direction === "down") {
      updatedLine = [...updatedLine].reverse();
    }

    if (direction === "left" || direction === "right") {
      nextBoard[index] = updatedLine;
    } else {
      for (let row = 0; row < size; row += 1) {
        nextBoard[row][index] = updatedLine[row];
      }
    }
  }

  const moved = nextBoard.some((row, rowIndex) =>
    row.some((value, colIndex) => value !== currentBoard[rowIndex][colIndex])
  );

  return { moved, nextBoard, gained, won };
}

/**
 * Adds a directional CSS class so tiles appear to slide before the state update.
 *
 * @param {"left"|"right"|"up"|"down"} direction The requested move direction.
 */
function animateBoard(direction) {
  boardElement.classList.add(`motion-${direction}`);
}

/**
 * Removes the temporary movement class after the transition finishes.
 *
 * @param {"left"|"right"|"up"|"down"} direction The completed move direction.
 */
function clearBoardAnimation(direction) {
  boardElement.classList.remove(`motion-${direction}`);
}

/**
 * Runs one full move cycle: calculate, animate, update score, spawn a new tile,
 * and show win or loss overlays when appropriate.
 *
 * @param {"left"|"right"|"up"|"down"} direction The requested move direction.
 */
function move(direction) {
  if (isAnimating) {
    return;
  }

  const previousBoard = cloneBoard(board);
  const result = calculateMove(previousBoard, direction);

  if (!result.moved) {
    return;
  }

  isAnimating = true;
  animateBoard(direction);

  window.setTimeout(() => {
    clearBoardAnimation(direction);
    // The board updates after the brief slide effect so the motion reads first.
    board = result.nextBoard;
    score += result.gained;
    updateScore();
    addRandomTile();
    renderBoard();

    if (result.won && !gameWon) {
      gameWon = true;
      showOverlay("You Win!", "You reached 2048. Keep playing or start a fresh game.");
    } else if (!hasMoves()) {
      showOverlay("Game Over", "No more moves are available. Start a new game and try again.");
    }

    isAnimating = false;
  }, moveAnimationMs);
}

/**
 * Determines whether at least one legal move remains on the board.
 *
 * @returns {boolean} True when an empty cell or valid merge still exists.
 */
function hasMoves() {
  for (let row = 0; row < size; row += 1) {
    for (let col = 0; col < size; col += 1) {
      const current = board[row][col];
      if (current === 0) {
        return true;
      }
      if (col < size - 1 && current === board[row][col + 1]) {
        return true;
      }
      if (row < size - 1 && current === board[row + 1][col]) {
        return true;
      }
    }
  }

  return false;
}

/**
 * Displays the end-state overlay with the supplied title and message.
 *
 * @param {string} title The overlay headline.
 * @param {string} message The supporting message shown below the title.
 */
function showOverlay(title, message) {
  overlayTitle.textContent = title;
  overlayMessage.textContent = message;
  overlay.classList.remove("hidden");
}

/**
 * Hides the overlay so the player can continue interacting with the board.
 */
function hideOverlay() {
  overlay.classList.add("hidden");
}

/**
 * Opens the modal that explains the controls and objective of the game.
 */
function openInfoModal() {
  infoModal.classList.remove("hidden");
}

/**
 * Closes the rules modal.
 */
function closeInfoModal() {
  infoModal.classList.add("hidden");
}

/**
 * Sends a shutdown request to the Go server and reports the result to the user.
 *
 * @returns {Promise<void>} Resolves once the request attempt has completed.
 */
async function exitGame() {
  try {
    await fetch("/exit", {
      method: "POST",
      headers: { "Content-Type": "application/json" }
    });
    showOverlay("Application Closed", "The Go server has been asked to shut down. You can close this tab.");
  } catch (error) {
    showOverlay("Exit Failed", "The app could not contact the server to exit cleanly.");
  }
}

document.addEventListener("keydown", (event) => {
  if (isAnimating) {
    return;
  }

  if (!infoModal.classList.contains("hidden")) {
    if (event.key === "Escape") {
      closeInfoModal();
    }
    return;
  }

  const key = event.key.toLowerCase();
  const directionMap = {
    arrowup: "up",
    w: "up",
    arrowdown: "down",
    s: "down",
    arrowleft: "left",
    a: "left",
    arrowright: "right",
    d: "right"
  };

  const direction = directionMap[key];
  if (!direction) {
    return;
  }

  // Prevent the browser from scrolling when the player uses arrow keys.
  event.preventDefault();
  move(direction);
});

newGameButton.addEventListener("click", startGame);
overlayButton.addEventListener("click", startGame);
infoButton.addEventListener("click", openInfoModal);
closeModalButton.addEventListener("click", closeInfoModal);
exitButton.addEventListener("click", exitGame);

infoModal.addEventListener("click", (event) => {
  if (event.target === infoModal) {
    closeInfoModal();
  }
});

startGame();
