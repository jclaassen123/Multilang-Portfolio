import java.awt.*
import javax.swing.*
import java.util.*

/**
 * Queue-based Tic Tac Toe Game.
 *
 * Each player (X and O) can have at most 3 active marks on the board.
 * When placing a 4th mark, the oldest mark for that player is removed.
 * The oldest mark that would be removed is displayed in red.
 *
 * Winning condition: 3 in a row horizontally, vertically, or diagonally.
 * Winning marks are displayed in green. Players can choose to play again after a win.
 */
class TicTacToe : JFrame() {

    private val boardSize = 3
    private val buttons = Array(boardSize) { arrayOfNulls<JButton>(boardSize) }

    private val xMoves: Queue<JButton> = LinkedList()
    private val oMoves: Queue<JButton> = LinkedList()
    private var currentPlayer = "X"

    /**
     * Initializes the GUI window, title, grid, and rules button.
     */
    init {
        title = "Queue Tic Tac Toe"
        defaultCloseOperation = EXIT_ON_CLOSE
        layout = BorderLayout()

        // Title label
        val titleLabel = JLabel("Queue Tic Tac Toe", SwingConstants.CENTER)
        titleLabel.font = Font("Arial", Font.BOLD, 24)
        add(titleLabel, BorderLayout.NORTH)

        // Grid for Tic Tac Toe
        val grid = JPanel(GridLayout(boardSize, boardSize, 5, 5))
        for (i in 0 until boardSize) {
            for (j in 0 until boardSize) {
                val btn = JButton("")
                btn.font = Font("Arial", Font.BOLD, 32)
                btn.addActionListener { handleMove(btn) }
                buttons[i][j] = btn
                grid.add(btn)
            }
        }
        add(grid, BorderLayout.CENTER)

        // Rules button
        val rulesButton = JButton("View Rules")
        rulesButton.addActionListener { showRules() }
        add(rulesButton, BorderLayout.SOUTH)

        setSize(400, 450)
        setLocationRelativeTo(null)
        isVisible = true
    }

    /**
     * Handles a player's move when a button is clicked.
     *
     * @param button The JButton clicked by the player.
     */
    private fun handleMove(button: JButton) {
        if (button.text.isNotEmpty()) return

        val currentQueue = if (currentPlayer == "X") xMoves else oMoves

        // Remove oldest mark if queue already has 3
        if (currentQueue.size == 3) {
            val oldest = currentQueue.poll()
            oldest?.text = ""
            oldest?.foreground = Color.BLACK
        }

        // Place new move
        button.text = currentPlayer
        button.foreground = Color.BLACK
        currentQueue.add(button)

        // Highlight the oldest mark if applicable
        highlightOldest()

        // Check for win
        val winningButtons = getWinningButtons(currentPlayer)
        if (winningButtons.isNotEmpty()) {
            // Highlight winning row/column/diagonal in green
            winningButtons.forEach { it.foreground = Color.GREEN }

            // Ask user to play again
            val option = JOptionPane.showConfirmDialog(
                this,
                "Player $currentPlayer wins! Play again?",
                "Game Over",
                JOptionPane.YES_NO_OPTION
            )

            if (option == JOptionPane.YES_OPTION) {
                resetBoard()
            } else {
                dispose()
            }
            return
        }

        // Switch player
        currentPlayer = if (currentPlayer == "X") "O" else "X"

        // Update highlight for next turn
        highlightOldest()
    }

    /**
     * Highlights the oldest mark of the current player in red,
     * only if it is the next to be removed (queue is full).
     */
    private fun highlightOldest() {
        // Reset all marks to black
        xMoves.forEach { it.foreground = Color.BLACK }
        oMoves.forEach { it.foreground = Color.BLACK }

        // Highlight oldest if current player has 3 marks
        val queue = if (currentPlayer == "X") xMoves else oMoves
        if (queue.size == 3) {
            queue.peek()?.foreground = Color.RED
        }
    }

    /**
     * Determines which buttons form a winning row, column, or diagonal.
     *
     * @param player The current player ("X" or "O")
     * @return A list of buttons forming the winning combination, empty if none.
     */
    private fun getWinningButtons(player: String): List<JButton> {
        val winning = mutableListOf<JButton>()

        // Check rows
        for (i in 0 until boardSize) {
            if ((0 until boardSize).all { j -> buttons[i][j]?.text == player }) {
                for (j in 0 until boardSize) buttons[i][j]?.let { winning.add(it) }
                return winning
            }
        }

        // Check columns
        for (j in 0 until boardSize) {
            if ((0 until boardSize).all { i -> buttons[i][j]?.text == player }) {
                for (i in 0 until boardSize) buttons[i][j]?.let { winning.add(it) }
                return winning
            }
        }

        // Check main diagonal
        if ((0 until boardSize).all { i -> buttons[i][i]?.text == player }) {
            for (i in 0 until boardSize) buttons[i][i]?.let { winning.add(it) }
            return winning
        }

        // Check anti-diagonal
        if ((0 until boardSize).all { i -> buttons[i][boardSize - 1 - i]?.text == player }) {
            for (i in 0 until boardSize) buttons[i][boardSize - 1 - i]?.let { winning.add(it) }
            return winning
        }

        return emptyList()
    }

    /**
     * Displays the game rules in a popup.
     */
    private fun showRules() {
        JOptionPane.showMessageDialog(
            this,
            """
            Queue Tic Tac Toe Rules:

            1. Players take turns, X always goes first.
            2. Each player can have at most 3 active marks on the board.
            3. If a player places a 4th mark, their oldest mark is removed.
            4. The oldest mark that would be removed is shown in red.
            5. Winning condition: 3 marks in a row horizontally, vertically, or diagonally.
            6. Winning marks are highlighted in green. After a win, you can choose to play again.
            """.trimIndent()
        )
    }

    /**
     * Resets the board and queues for a new game.
     */
    private fun resetBoard() {
        for (i in 0 until boardSize) {
            for (j in 0 until boardSize) {
                buttons[i][j]?.text = ""
                buttons[i][j]?.foreground = Color.BLACK
            }
        }
        xMoves.clear()
        oMoves.clear()
        currentPlayer = "X"
    }
}

/**
 * Main entry point to launch the Tic Tac Toe GUI.
 */
fun main() {
    SwingUtilities.invokeLater { TicTacToe() }
}
