package com.example.memorygame.game

/**
 * Represents a single tile on the memory game board.
 *
 * @property id unique identifier for the tile
 * @property imageId identifier for the image associated with this tile; used for matching pairs
 * @property isMatched indicates whether this tile has already been matched
 */
data class Tile(
    val id: Int,
    val imageId: Int,
    var isMatched: Boolean = false // Tracks if the tile has been successfully matched
)
