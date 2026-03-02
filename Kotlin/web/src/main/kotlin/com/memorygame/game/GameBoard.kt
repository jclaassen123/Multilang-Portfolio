package com.example.memorygame.game

/**
 * Represents the game board for the memory tile-matching game.
 *
 * The board consists of a grid of tiles, each with an `imageId`. The goal
 * is to match pairs of tiles with the same `imageId`. This class handles
 * tile generation, flipping logic, and tracking of guesses.
 *
 * @property rows number of rows on the board
 * @property cols number of columns on the board
 * @property tiles the list of tiles on the board
 * @property guesses total number of guess attempts made by the player
 */
class GameBoard(val rows: Int, val cols: Int, images: Int) {

    /** List of all tiles on the board, initialized during board creation */
    val tiles: List<Tile>

    /** Total number of guesses made by the player */
    var guesses = 0
        private set

    init {
        // Ensure the board has an even number of tiles to form pairs
        require(rows * cols % 2 == 0) { "Board must have even number of tiles" }
        tiles = generateTiles()
    }

    /**
     * Generates a shuffled list of tiles with matching pairs.
     *
     * Each tile is assigned an `imageId` that represents its pair. The list
     * of tiles is then shuffled to randomize their positions on the board.
     *
     * @return a shuffled list of Tile objects
     */
    private fun generateTiles(): List<Tile> {
        val pairs = (rows * cols) / 2
        val values = (1..pairs).flatMap { listOf(it, it) }.shuffled()
        return values.mapIndexed { index, value -> Tile(id = index, imageId = value) }
    }

    /**
     * Handles flipping a tile on the board.
     *
     * This method determines if a flip is the first or second in a turn,
     * checks for matches, updates matched state, increments guess count,
     * and returns the appropriate [FlipResult].
     *
     * @param tileId the ID of the tile being flipped
     * @param firstTileId the ID of the first tile flipped in the current turn,
     *                    or null if this is the first flip
     * @return the result of the flip as a [FlipResult]
     */
    fun flip(tileId: Int, firstTileId: Int?): FlipResult {
        val tile = tiles[tileId]

        // Ignore flipping tiles that have already been matched
        if (tile.isMatched) return FlipResult.First(tile)

        // This is the first tile flipped in the current turn
        if (firstTileId == null) {
            return FlipResult.First(tile)
        }

        // This is the second tile flipped; check for a match
        val firstTile = tiles[firstTileId]
        guesses++

        return if (firstTile.imageId == tile.imageId) {
            // Mark both tiles as matched
            firstTile.isMatched = true
            tile.isMatched = true

            // Check if all tiles have been matched
            val complete = tiles.all { it.isMatched }
            if (complete) FlipResult.Complete(firstTile, tile, guesses)
            else FlipResult.Match(firstTile, tile)
        } else {
            // Tiles do not match
            FlipResult.Miss(firstTile, tile)
        }
    }
}