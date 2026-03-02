package com.example.memorygame.controller

import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping

/**
 * Controller for rendering static pages of the memory game.
 *
 * This controller serves the HTML pages for the game interface and leaderboard.
 * It does not handle any game logic or API calls.
 */
@Controller
class PageController {

    /**
     * Renders the main memory game page.
     *
     * @return the name of the HTML template to render ("game.html")
     */
    @GetMapping("/")
    fun game(): String {
        return "game"
    }

    /**
     * Renders the leaderboard page showing top scores.
     *
     * @return the name of the HTML template to render ("leaderboard.html")
     */
    @GetMapping("/leaderboard")
    fun leaderboard(): String {
        return "leaderboard"
    }
}