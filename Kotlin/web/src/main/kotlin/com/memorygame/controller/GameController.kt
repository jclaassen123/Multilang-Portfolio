package com.example.memorygame.controller

import com.example.memorygame.game.GameBoard
import com.example.memorygame.game.FlipResult
import com.example.memorygame.game.GameResult
import com.example.memorygame.repository.GameResultRepository
import org.springframework.web.bind.annotation.*
import java.time.format.DateTimeFormatter
import java.util.*

/**
 * REST controller that handles memory game operations such as starting a game,
 * flipping tiles, submitting a username for leaderboard, and fetching leaderboard results.
 */
@RestController
@RequestMapping("/api")
class GameController(
    private val gameResultRepository: GameResultRepository
) {

    // Map to track active game sessions by session ID
    private val sessions = mutableMapOf<String, GameBoard>()

    /**
     * Response returned when a new game is started.
     *
     * @property sessionId Unique ID for the game session.
     * @property tiles List of tiles with their IDs and image IDs.
     * @property rows Number of rows in the game board.
     * @property cols Number of columns in the game board.
     */
    data class StartResponse(
        val sessionId: String,
        val tiles: List<Map<String, Any>>,
        val rows: Int,
        val cols: Int
    )

    /**
     * Response returned when a tile is flipped.
     *
     * @property match True if the flip resulted in a match.
     * @property miss True if the flip did not result in a match.
     * @property complete True if the game is completed after this flip.
     * @property tile1Id ID of the first flipped tile in the current flip operation.
     * @property tile2Id ID of the second flipped tile in the current flip operation.
     */
    data class FlipResponse(
        val match: Boolean = false,
        val miss: Boolean = false,
        val complete: Boolean = false,
        val tile1Id: Int = -1,
        val tile2Id: Int = -1
    )

    /**
     * Starts a new memory game.
     *
     * @param rows Number of rows for the game board.
     * @param cols Number of columns for the game board.
     * @return StartResponse containing session ID and initial board state.
     */
    @PostMapping("/start")
    fun startGame(
        @RequestParam rows: Int,
        @RequestParam cols: Int
    ): StartResponse {
        // Initialize a new game board with 10 unique images
        val board = GameBoard(rows, cols, 10)
        val sessionId = UUID.randomUUID().toString()
        sessions[sessionId] = board

        // Convert tiles to a lightweight representation for the client
        val tiles = board.tiles.map {
            mapOf("id" to it.id, "imageId" to it.imageId)
        }

        return StartResponse(sessionId, tiles, rows, cols)
    }

    /**
     * Handles a tile flip in an existing game session.
     *
     * @param sessionId ID of the game session.
     * @param tileId ID of the tile to flip.
     * @param firstTileId Optional ID of the first tile if this is the second flip.
     * @return FlipResponse indicating match, miss, or completion status.
     * @throws IllegalArgumentException if the session does not exist.
     */
    @PostMapping("/flip/{sessionId}/{tileId}")
    fun flip(
        @PathVariable sessionId: String,
        @PathVariable tileId: Int,
        @RequestParam(required = false) firstTileId: Int?
    ): FlipResponse {

        val board = sessions[sessionId] ?: error("Game not found")
        val result = board.flip(tileId, firstTileId)

        return when (result) {
            is FlipResult.Match -> FlipResponse(
                match = true,
                tile1Id = result.t1.id,
                tile2Id = result.t2.id
            )
            is FlipResult.Miss -> FlipResponse(
                miss = true,
                tile1Id = result.t1.id,
                tile2Id = result.t2.id
            )
            is FlipResult.Complete -> FlipResponse(
                match = true,
                complete = true,
                tile1Id = result.t1.id,
                tile2Id = result.t2.id
            )
            else -> FlipResponse()
        }
    }

    /**
     * Retrieves the top 10 leaderboard entries for a given board size.
     *
     * @param rows Number of rows in the board.
     * @param cols Number of columns in the board.
     * @return List of maps containing username, guesses, and completion date.
     */
    @GetMapping("/leaderboard")
    fun leaderboard(
        @RequestParam rows: Int,
        @RequestParam cols: Int
    ): List<Map<String, Any>> {

        val formatter = DateTimeFormatter.ofPattern("MM-dd-yy")

        return gameResultRepository
            .findTop10ByRowsAndColsOrderByGuessesAsc(rows, cols)
            .map {
                mapOf(
                    "guesses" to it.guesses,
                    "completedAt" to it.completedAt.format(formatter),
                    "username" to it.username
                )
            }
    }

    /**
     * Submits a username for the current game session and saves the result to the leaderboard.
     *
     * @param sessionId ID of the game session.
     * @param username Player's username.
     * @return "ok" if successfully saved.
     * @throws IllegalArgumentException if the session does not exist.
     */
    @PostMapping("/submit-username/{sessionId}")
    fun submitUsername(
        @PathVariable sessionId: String,
        @RequestParam username: String
    ): String {
        val board = sessions[sessionId] ?: error("Game not found")

        // Save the completed game result
        val lastResult = GameResult(
            rows = board.rows,
            cols = board.cols,
            guesses = board.guesses,
            username = username
        )
        gameResultRepository.save(lastResult)

        return "ok"
    }
}