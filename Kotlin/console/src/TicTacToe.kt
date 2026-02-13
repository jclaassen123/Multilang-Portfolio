import java.util.ArrayDeque

private const val RED = "\u001B[31m"
private const val GREEN = "\u001B[32m"
private const val RESET = "\u001B[0m"

/**
 * Main entry point for the game.
 * Allows repeated play until the user chooses to exit.
 */
fun main() {
    do {
        playGame()
    } while (askPlayAgain())
}

/**
 * Runs a single game of Queue Tic-Tac-Toe.
 * Manages game loop, board state, player turns, and queue mechanics.
 */
fun playGame() {
    printInstructions()

    // Initialize board and player queues
    val board = Array(3) { Array<Char?>(3) { null } }
    val xQueue = ArrayDeque<Pair<Int, Int>>()
    val oQueue = ArrayDeque<Pair<Int, Int>>()
    var currentPlayer = 'X'
    var winningLine: List<Pair<Int, Int>>? = null

    while (true) {
        val queue = if (currentPlayer == 'X') xQueue else oQueue
        var highlightCell: Pair<Int, Int>? = null

        // Highlight the oldest marker if player has 3 markers on board
        if (queue.size == 3) {
            highlightCell = queue.first()
            println("Player $currentPlayer will remove their oldest marker at " +
                    "${highlightCell.first} ${('A' + highlightCell.second)}")
        }

        // Print the current board with optional highlights
        printBoard(board, highlightCell, winningLine)

        println("Player $currentPlayer's turn")
        println("Enter your move (example: 1B, B1, 2c, c2)")
        print("> ")

        // Read input and normalize (remove spaces, uppercase)
        val inputRaw = readLine()?.trim()?.replace(" ", "")?.uppercase()
        if (inputRaw == null || inputRaw.length != 2) {
            println("${RED}Invalid input. Example: 1B or B1$RESET")
            continue
        }

        // Separate digits and letters from input
        val digits = inputRaw.filter { it.isDigit() }
        val letters = inputRaw.filter { it.isLetter() }

        if (digits.length != 1 || letters.length != 1) {
            println("${RED}Invalid input. Must contain one number and one letter (0-2, A-C).$RESET")
            continue
        }

        // Parse row (0-2)
        val row = digits.toIntOrNull()
        // Parse column (A-C)
        val col = when (letters[0]) {
            'A' -> 0
            'B' -> 1
            'C' -> 2
            else -> null
        }

        if (row == null || row !in 0..2 || col == null) {
            println("${RED}Invalid row or column. Rows 0-2, columns A-C.$RESET")
            continue
        }

        // Check if selected square is empty
        if (board[row][col] != null) {
            println("${RED}That square is already occupied. Try another.$RESET")
            continue
        }

        // Remove oldest marker if player already has 3
        if (queue.size == 3) {
            val (oldRow, oldCol) = queue.removeFirst()
            board[oldRow][oldCol] = null
        }

        // Place new marker on the board and add to queue
        board[row][col] = currentPlayer
        queue.addLast(Pair(row, col))

        // Check for win condition
        val winningCells = checkWin(board, currentPlayer)
        if (winningCells != null) {
            winningLine = winningCells
            printBoard(board, null, winningLine)
            println("${GREEN}Player $currentPlayer wins the game!$RESET")
            break
        }

        // Switch to the next player
        currentPlayer = if (currentPlayer == 'X') 'O' else 'X'
    }
}

/**
 * Prints instructions for how to play the game.
 */
fun printInstructions() {
    println("======================================")
    println("        QUEUE TIC-TAC-TOE")
    println("======================================")
    println("How to play:")
    println("• 3×3 board, players take turns placing X or O")
    println("• Rows: 0, 1, 2 | Columns: A, B, C")
    println("• Input can be row+column or column+row (example: 1B or B1)")
    println("• No spaces needed, case-insensitive")
    println()
    println("Queue Rule:")
    println("• Max 3 markers per player on board")
    println("• Placing a 4th marker removes the oldest")
    println("• Marker about to be removed is shown RED on the board")
    println()
    println("Win: 3 in a row, column, or diagonal (highlighted GREEN)")
    println("======================================\n")
}

/**
 * Prints the current board state.
 *
 * @param board The 3x3 board array
 * @param highlightCell Optional cell to highlight in RED (about-to-be-removed marker)
 * @param winningLine Optional list of cells to highlight in GREEN (winning line)
 */
fun printBoard(
    board: Array<Array<Char?>>,
    highlightCell: Pair<Int, Int>? = null,
    winningLine: List<Pair<Int, Int>>? = null
) {
    println("      A   B   C")
    println("    ┌───┬───┬───┐")
    for (i in 0..2) {
        print("  $i │")
        for (j in 0..2) {
            val cell = board[i][j]?.toString() ?: " "
            when {
                winningLine != null && winningLine.contains(i to j) -> print(" $GREEN$cell$RESET │")
                highlightCell != null && highlightCell.first == i && highlightCell.second == j -> print(" $RED$cell$RESET │")
                else -> print(" $cell │")
            }
        }
        println()
        if (i < 2) println("    ├───┼───┼───┤")
    }
    println("    └───┴───┴───┘")
    println()
}

/**
 * Checks if the specified player has a winning line.
 *
 * @param board The current board state
 * @param player The player character ('X' or 'O')
 * @return List of winning cells if the player wins, otherwise null
 */
fun checkWin(board: Array<Array<Char?>>, player: Char): List<Pair<Int, Int>>? {
    // Check rows
    for (i in 0..2) if ((0..2).all { board[i][it] == player }) return (0..2).map { i to it }
    // Check columns
    for (i in 0..2) if ((0..2).all { board[it][i] == player }) return (0..2).map { it to i }
    // Check diagonal top-left -> bottom-right
    if ((0..2).all { board[it][it] == player }) return (0..2).map { it to it }
    // Check diagonal top-right -> bottom-left
    if ((0..2).all { board[it][2 - it] == player }) return (0..2).map { it to 2 - it }
    return null
}

/**
 * Asks the user whether they want to play again.
 *
 * @return true if the user wants to play again, false otherwise
 */
fun askPlayAgain(): Boolean {
    println("Do you want to play again? (Y/N)")
    while (true) {
        print("> ")
        val input = readLine()?.trim()?.uppercase()
        when (input) {
            "Y", "YES" -> return true
            "N", "NO" -> return false
            else -> println("${RED}Please enter Y or N$RESET")
        }
    }
}
