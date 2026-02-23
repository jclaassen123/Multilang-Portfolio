package com.example.memorygame.game

/**
 * Represents the result of flipping tiles in the memory game.
 * This is a sealed class, so all possible flip outcomes are defined here.
 */
sealed class FlipResult {

    /**
     * Indicates that the first tile has been flipped and no match attempt has occurred yet.
     * @property tile the first tile that was flipped
     */
    data class First(val tile: Tile) : FlipResult()
    // First flip scenario: player has flipped one tile but hasn't flipped a second one yet.

    /**
     * Indicates that two tiles were flipped and they match.
     * @property t1 the first tile flipped
     * @property t2 the second tile flipped that matches the first
     */
    data class Match(val t1: Tile, val t2: Tile) : FlipResult()
    // Successful match scenario: both flipped tiles are the same.

    /**
     * Indicates that two tiles were flipped but do not match.
     * @property t1 the first tile flipped
     * @property t2 the second tile flipped that does not match the first
     */
    data class Miss(val t1: Tile, val t2: Tile) : FlipResult()
    // Failed match scenario: tiles are different, so the player will need to try again.

    /**
     * Indicates that the game has been completed successfully.
     * @property t1 the last flipped tile
     * @property t2 the second-to-last flipped tile
     * @property guesses the total number of guess attempts made in the game
     */
    data class Complete(val t1: Tile, val t2: Tile, val guesses: Int) : FlipResult()
    // Game completion scenario: all tiles are matched, includes total guesses for scoring.
}