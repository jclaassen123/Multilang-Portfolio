package com.example.memorygame.repository

import com.example.memorygame.game.GameResult
import org.springframework.data.jpa.repository.JpaRepository

/**
 * Repository interface for accessing [GameResult] entities from the database.
 *
 * Extends [JpaRepository] to provide standard CRUD operations.
 * Includes a custom query method to retrieve the top 10 scores for a given board size.
 */
interface GameResultRepository : JpaRepository<GameResult, Long> {

    /**
     * Finds the top 10 game results for a specific board size, ordered by the fewest guesses.
     *
     * @param rows Number of rows in the game board.
     * @param cols Number of columns in the game board.
     * @return List of up to 10 [GameResult] objects with the fewest guesses for the given board size.
     */
    fun findTop10ByRowsAndColsOrderByGuessesAsc(
        rows: Int,
        cols: Int
    ): List<GameResult>
}