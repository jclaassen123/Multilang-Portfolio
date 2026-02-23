package com.example.memorygame.controller

import com.example.memorygame.game.GameBoard
import com.example.memorygame.game.FlipResult
import org.springframework.web.bind.annotation.*
import java.util.*

/**
 * REST controller that manages the memory game sessions.
 *
 * Provides endpoints to start a new game and flip tiles.
 */
@RestController
@RequestMapping("/api")
class GameController {

    // Stores active game sessions with a unique session ID
    private val sessions = mutableMapOf<String, GameBoard>()

    /**
     * Response returned when a new game is started.
     * @property sessionId unique ID for this game session
     * @property tiles list of tiles with their IDs and image IDs
     * @property rows number of rows on the board
     * @property cols number of columns on the board
     */
    data class StartResponse(val sessionId: String, val tiles: List<Map<String, Any>>, val rows: Int, val cols: Int)

    /**
     * Response returned when a tile is flipped.
     * @property match indicates a successful match
     * @property miss indicates a failed match
     * @property complete indicates the game is complete
     * @property tile1Id ID of the first flipped tile
     * @property tile2Id ID of the second flipped tile
     */
    data class FlipResponse(
        val match: Boolean = false,
        val miss: Boolean = false,
        val complete: Boolean = false,
        val tile1Id: Int = -1,
        val tile2Id: Int = -1
    )

    /**
     * Starts a new memory game session with the given number of rows and columns.
     * @param rows number of rows for the game board
     * @param cols number of columns for the game board
     * @return a StartResponse containing session details and the initial tile list
     */
    @PostMapping("/start")
    fun startGame(@RequestParam rows: Int, @RequestParam cols: Int): StartResponse {
        val board = GameBoard(rows, cols, 10)
        val sessionId = UUID.randomUUID().toString()
        sessions[sessionId] = board

        val tiles = board.tiles.map { mapOf("id" to it.id, "imageId" to it.imageId) }
        return StartResponse(sessionId, tiles, rows, cols)
    }

    /**
     * Handles flipping a tile in a game session.
     * @param sessionId the unique ID of the game session
     * @param tileId the ID of the tile to flip
     * @param firstTileId optional ID of the first tile flipped in the current turn
     * @return a FlipResponse indicating the result of the flip
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
            is FlipResult.Match -> {
                val complete = board.tiles.all { it.isMatched }
                FlipResponse(
                    match = true,
                    complete = complete,
                    tile1Id = result.t1.id,
                    tile2Id = result.t2.id
                )
            }
            is FlipResult.Miss -> FlipResponse(miss = true, tile1Id = result.t1.id, tile2Id = result.t2.id)
            is FlipResult.Complete -> FlipResponse(
                match = true,
                complete = true,
                tile1Id = result.t1.id,
                tile2Id = result.t2.id
            )
            else -> FlipResponse()
        }
    }
}
