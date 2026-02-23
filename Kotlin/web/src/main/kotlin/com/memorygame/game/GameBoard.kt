package com.example.memorygame.game

/**
 * Represents the game board for the memory tile matching game.
 *
 * @property rows number of rows on the board
 * @property cols number of columns on the board
 * @property tiles the list of tiles on the board
 * @property guesses the total number of guess attempts made by the player
 */
class GameBoard(val rows: Int, val cols: Int, images: Int) {

    val tiles: List<Tile>

    var guesses = 0
        private set

    init {
        require(rows * cols % 2 == 0) { "Board must have even number of tiles" }
        tiles = generateTiles()
    }

    /**
     * Generates a shuffled list of tiles with matching pairs.
     *
     * Each tile is assigned an `imageId` representing its matching pair.
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
     * @param tileId the ID of the tile being flipped
     * @param firstTileId the ID of the first tile flipped in the current turn, or null if this is the first flip
     * @return the result of the flip as a FlipResult
     */
    fun flip(tileId: Int, firstTileId: Int?): FlipResult {
        val tile = tiles[tileId]

        // Ignore already matched tiles
        if (tile.isMatched) return FlipResult.First(tile)

        // First tile of a pair
        if (firstTileId == null) {
            return FlipResult.First(tile)
        }

        // Second tile of a pair
        val firstTile = tiles[firstTileId]
        guesses++

        return if (firstTile.imageId == tile.imageId) {
            firstTile.isMatched = true
            tile.isMatched = true

            val complete = tiles.all { it.isMatched }
            if (complete) FlipResult.Complete(firstTile, tile, guesses)
            else FlipResult.Match(firstTile, tile)
        } else {
            FlipResult.Miss(firstTile, tile)
        }
    }
}
