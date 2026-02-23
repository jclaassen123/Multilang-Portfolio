package com.example.memorygame.controller

import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping

/**
 * Controller for rendering static pages related to the memory game.
 */
@Controller
class PageController {

    /**
     * Renders the main game page.
     *
     * @return the name of the HTML template to render (game.html)
     */
    @GetMapping("/")
    fun game(): String {
        return "game"
    }
}
