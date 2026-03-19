package main

import (
	"fmt"
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 1024
	screenHeight = 720
)

type Round int

const (
	RoundColor Round = iota + 1
	RoundHighLow
	RoundInsideOutside
	RoundSuit
)

type GameState int

const (
	StatePlaying GameState = iota
	StateWon
	StateLost
)

type Button struct {
	X, Y, W, H int
	Label      string
	Fill       color.RGBA
	Action     func()
}

type Game struct {
	deck       Deck
	round      Round
	state      GameState
	firstCard  Card
	secondCard Card
	missedCard *Card
	attempt    []Card
	question   string
	message    string

	sprites CardSprites
	buttons []Button

	mouseWasDown     bool
	quit             bool
	showInstructions bool
}

// NewGame initializes the game state and loads the card sprite assets.
func NewGame() *Game {
	// Load shared card art once before any rendering happens.
	sprites, err := loadCardSprites("assets/cards.png")
	if err != nil {
		log.Fatal(err)
	}

	// Start from a fresh match instead of exposing partially initialized state.
	g := &Game{sprites: sprites}
	g.startNewMatch()
	return g
}

// startNewMatch resets the game to the first round with a fresh shuffled deck.
func (g *Game) startNewMatch() {
	// Rebuild and reshuffle the deck so each match is independent.
	g.deck = NewDeck()
	g.deck.Shuffle()

	// Clear all progress markers from the previous run.
	g.state = StatePlaying
	g.round = RoundColor
	g.firstCard = Card{}
	g.secondCard = Card{}
	g.missedCard = nil
	g.attempt = nil
	g.message = "Get 4 correct in a row before the deck runs out."
	g.setupRoundButtons()
}

// setupRoundButtons configures the prompt and actions for the current round.
func (g *Game) setupRoundButtons() {
	// Each round swaps in a different prompt and button set.
	switch g.round {
	case RoundColor:
		g.question = "Round 1: Is the next card red or black?"
		g.buttons = []Button{
			{X: 250, Y: 642, W: 220, H: 58, Label: "Red", Fill: color.RGBA{186, 46, 46, 255}, Action: func() { g.guessColor(true) }},
			{X: 550, Y: 642, W: 220, H: 58, Label: "Black", Fill: color.RGBA{36, 36, 36, 255}, Action: func() { g.guessColor(false) }},
		}
	case RoundHighLow:
		g.question = fmt.Sprintf("Round 2: Higher or lower than %s of %s?", g.firstCard.Name, g.firstCard.Suit)
		g.buttons = []Button{
			{X: 250, Y: 642, W: 220, H: 58, Label: "Higher", Fill: color.RGBA{57, 110, 198, 255}, Action: func() { g.guessHigherLower(true) }},
			{X: 550, Y: 642, W: 220, H: 58, Label: "Lower", Fill: color.RGBA{57, 110, 198, 255}, Action: func() { g.guessHigherLower(false) }},
		}
	case RoundInsideOutside:
		g.question = fmt.Sprintf("Round 3: Is the next card inside or outside %s and %s?", g.firstCard.Name, g.secondCard.Name)
		g.buttons = []Button{
			{X: 250, Y: 642, W: 220, H: 58, Label: "Inside", Fill: color.RGBA{92, 92, 168, 255}, Action: func() { g.guessInsideOutside(true) }},
			{X: 550, Y: 642, W: 220, H: 58, Label: "Outside", Fill: color.RGBA{92, 92, 168, 255}, Action: func() { g.guessInsideOutside(false) }},
		}
	case RoundSuit:
		g.question = "Round 4: Choose the suit"
		g.buttons = []Button{
			{X: 90, Y: 642, W: 200, H: 58, Label: "Hearts", Fill: color.RGBA{186, 46, 46, 255}, Action: func() { g.guessSuit("Hearts") }},
			{X: 310, Y: 642, W: 200, H: 58, Label: "Diamonds", Fill: color.RGBA{214, 67, 67, 255}, Action: func() { g.guessSuit("Diamonds") }},
			{X: 530, Y: 642, W: 200, H: 58, Label: "Clubs", Fill: color.RGBA{30, 30, 30, 255}, Action: func() { g.guessSuit("Clubs") }},
			{X: 750, Y: 642, W: 180, H: 58, Label: "Spades", Fill: color.RGBA{30, 30, 30, 255}, Action: func() { g.guessSuit("Spades") }},
		}
	}
}

// setupPlayAgainButtons replaces the round actions with play-again controls.
func (g *Game) setupPlayAgainButtons() {
	g.question = "Play again?"

	// Once the game ends, the only valid actions are restart or quit.
	g.buttons = []Button{
		{X: 360, Y: 642, W: 140, H: 58, Label: "Yes", Fill: color.RGBA{46, 134, 78, 255}, Action: g.startNewMatch},
		{X: 520, Y: 642, W: 140, H: 58, Label: "No", Fill: color.RGBA{121, 44, 44, 255}, Action: func() { g.quit = true }},
	}
}

// drawCard draws the next card, updating loss state if the deck is empty.
func (g *Game) drawCard() (Card, bool) {
	card, ok := g.deck.Draw()
	if !ok {
		// Running out of cards ends the match immediately.
		g.state = StateLost
		g.message = "Out of cards. You lose."
		g.setupPlayAgainButtons()
		return Card{}, false
	}

	// Keep only the last four revealed cards for the history display.
	g.attempt = append(g.attempt, card)
	if len(g.attempt) > 4 {
		g.attempt = g.attempt[len(g.attempt)-4:]
	}
	return card, true
}

// resetToRoundOne sends the player back to round one after an incorrect guess.
func (g *Game) resetToRoundOne(wrongCard Card) {
	// Wrong guesses wipe round progress but preserve the missed card for feedback.
	g.round = RoundColor
	g.firstCard = Card{}
	g.secondCard = Card{}
	g.missedCard = &wrongCard
	g.attempt = nil
	g.message = fmt.Sprintf("Wrong: %s of %s. Back to Round 1.", wrongCard.Name, wrongCard.Suit)
	g.setupRoundButtons()
}

// guessColor resolves the red-or-black guess for round one.
func (g *Game) guessColor(guessRed bool) {
	card, ok := g.drawCard()
	if !ok {
		return
	}

	// A correct first guess advances the player and stores the reference card.
	if isRed(card) == guessRed {
		g.missedCard = nil
		g.firstCard = card
		g.round = RoundHighLow
		g.message = fmt.Sprintf("Correct: %s of %s.", card.Name, card.Suit)
		g.setupRoundButtons()
		return
	}
	g.resetToRoundOne(card)
}

// guessHigherLower resolves the higher-or-lower guess for round two.
func (g *Game) guessHigherLower(guessHigher bool) {
	card, ok := g.drawCard()
	if !ok {
		return
	}

	// Compare against the first revealed card from round one.
	correct := false
	if guessHigher {
		correct = card.Value > g.firstCard.Value
	} else {
		correct = card.Value < g.firstCard.Value
	}

	if correct {
		// The second card becomes the second anchor for round three.
		g.missedCard = nil
		g.secondCard = card
		g.round = RoundInsideOutside
		g.message = fmt.Sprintf("Correct: %s of %s.", card.Name, card.Suit)
		g.setupRoundButtons()
		return
	}
	g.resetToRoundOne(card)
}

// guessInsideOutside resolves the inside-or-outside guess for round three.
func (g *Game) guessInsideOutside(guessInside bool) {
	card, ok := g.drawCard()
	if !ok {
		return
	}

	// Normalize the anchor values so card order does not matter.
	low := min(g.firstCard.Value, g.secondCard.Value)
	high := max(g.firstCard.Value, g.secondCard.Value)
	isInside := card.Value > low && card.Value < high
	isOutside := card.Value < low || card.Value > high
	correct := (guessInside && isInside) || (!guessInside && isOutside)

	if correct {
		// A successful third round unlocks the suit guess to win.
		g.missedCard = nil
		g.round = RoundSuit
		g.message = fmt.Sprintf("Correct: %s of %s.", card.Name, card.Suit)
		g.setupRoundButtons()
		return
	}
	g.resetToRoundOne(card)
}

// guessSuit resolves the suit selection for the final round.
func (g *Game) guessSuit(suit string) {
	card, ok := g.drawCard()
	if !ok {
		return
	}

	// Matching the exact suit wins the full four-round sequence.
	if card.Suit == suit {
		g.missedCard = nil
		g.state = StateWon
		g.message = fmt.Sprintf("You won with %s of %s!", card.Name, card.Suit)
		g.setupPlayAgainButtons()
		return
	}
	g.resetToRoundOne(card)
}

// Update processes mouse input and triggers button actions for the current frame.
func (g *Game) Update() error {
	if g.quit {
		return ebiten.Termination
	}

	// Trigger clicks on the press edge so holding the mouse does not repeat actions.
	isDown := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if isDown && !g.mouseWasDown {
		mx, my := ebiten.CursorPosition()

		// Utility controls stay available regardless of game state.
		if g.handleButtonClick(g.utilityButtons(), mx, my) {
			g.mouseWasDown = isDown
			return nil
		}

		if g.showInstructions {
			g.handleButtonClick(g.instructionButtons(), mx, my)
			g.mouseWasDown = isDown
			return nil
		}

		// Walk the active buttons and fire the first one that contains the cursor.
		g.handleButtonClick(g.buttons, mx, my)
	}
	g.mouseWasDown = isDown
	return nil
}

// handleButtonClick triggers the first button containing the given mouse position.
func (g *Game) handleButtonClick(buttons []Button, mx, my int) bool {
	for _, b := range buttons {
		if mx >= b.X && mx <= b.X+b.W && my >= b.Y && my <= b.Y+b.H {
			b.Action()
			return true
		}
	}
	return false
}

// Draw renders the full game scene for the current state.
func (g *Game) Draw(screen *ebiten.Image) {
	// Paint the felt table background first, then layer framed panels on top.
	screen.Fill(color.RGBA{14, 82, 42, 255})

	ebitenutil.DrawRect(screen, 14, 14, 996, 692, color.RGBA{20, 98, 49, 255})
	ebitenutil.DrawRect(screen, 18, 18, 988, 82, color.RGBA{10, 48, 24, 255})
	ebitenutil.DrawRect(screen, 18, 108, 260, 450, color.RGBA{17, 68, 35, 255})
	ebitenutil.DrawRect(screen, 286, 108, 720, 450, color.RGBA{17, 68, 35, 255})
	ebitenutil.DrawRect(screen, 18, 566, 988, 140, color.RGBA{10, 48, 24, 255})

	// The left panel shows deck state, while the center shows recent cards and messages.
	drawCenteredScaledText(screen, "RIDE THE BUS", 18, 18, 988, 3.0, 18)

	g.drawDeckStack(screen, 82, 176)
	drawCenteredScaledText(screen, fmt.Sprintf("Cards Left: %d", len(g.deck.Cards)), 18, 462, 260, 2.0, 0)
	if g.state == StatePlaying {
		drawCenteredScaledText(screen, fmt.Sprintf("Round: %d / 4", g.round), 18, 500, 260, 2.0, 0)
	}

	g.drawHistory(screen, 326, 290)
	drawCenteredScaledText(screen, g.message, 286, 520, 720, 2.0, 0)

	// The footer always contains the current prompt and available actions.
	drawCenteredQuestion(screen, g.question, 18, 566, 988)

	for _, b := range g.buttons {
		g.drawButton(screen, b)
	}
	for _, b := range g.utilityButtons() {
		g.drawButton(screen, b)
	}

	if g.showInstructions {
		g.drawInstructionsModal(screen)
	}
}

// utilityButtons returns the always-available top-bar controls.
func (g *Game) utilityButtons() []Button {
	return []Button{
		{X: 728, Y: 34, W: 120, H: 42, Label: "How to Play", Fill: color.RGBA{191, 152, 64, 255}, Action: func() { g.showInstructions = true }},
		{X: 866, Y: 34, W: 100, H: 42, Label: "Quit", Fill: color.RGBA{121, 44, 44, 255}, Action: func() { g.quit = true }},
	}
}

// instructionButtons returns the controls shown inside the instructions modal.
func (g *Game) instructionButtons() []Button {
	return []Button{
		{X: 442, Y: 526, W: 140, H: 40, Label: "Close", Fill: color.RGBA{57, 110, 198, 255}, Action: func() { g.showInstructions = false }},
	}
}

// drawInstructionsModal renders the how-to-play overlay and its close button.
func (g *Game) drawInstructionsModal(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, 0, 0, screenWidth, screenHeight, color.RGBA{0, 0, 0, 170})
	ebitenutil.DrawRect(screen, 198, 118, 628, 458, color.RGBA{241, 231, 202, 255})
	ebitenutil.DrawRect(screen, 206, 126, 612, 442, color.RGBA{32, 74, 45, 255})

	ebitenutil.DebugPrintAt(screen, "HOW TO PLAY", 467, 158)

	paragraphs := []string{
		"The goal is to guess correctly 4 times in a row before the deck runs out of cards.",
		"1. Guess whether the first card is red or black.",
		"2. Guess whether the next card is higher or lower.",
		"3. Guess whether the third card is inside or outside the values of the first two cards.",
		"4. Guess the exact suit of the final card.",
		"A wrong guess sends you back to Round 1. You can open these instructions or quit at any time.",
	}

	y := 220
	for _, paragraph := range paragraphs {
		for _, line := range wrapText(paragraph, 72) {
			ebitenutil.DebugPrintAt(screen, line, 224, y)
			y += 24
		}
		y += 10
	}

	for _, b := range g.instructionButtons() {
		g.drawButton(screen, b)
	}
}

// drawDeckStack renders a stylized stack for the remaining deck.
func (g *Game) drawDeckStack(screen *ebiten.Image, x, y int) {
	if g.sprites.Back == nil {
		return
	}

	count := len(g.deck.Cards)

	// Convert remaining cards into a capped visual stack depth.
	layers := count / 5
	if layers < 1 && count > 0 {
		layers = 1
	}
	if layers > 10 {
		layers = 10
	}

	scale := 6.0
	for i := 0; i < layers; i++ {
		// Offset each layer slightly so the deck looks stacked instead of flat.
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(x+i*3), float64(y-i*2))
		screen.DrawImage(g.sprites.Back, op)
	}
	ebitenutil.DebugPrintAt(screen, "Deck", x+44, y+130)
}

// drawHistory renders the recent cards or the missed card in the history panel.
func (g *Game) drawHistory(screen *ebiten.Image, x, y int) {
	drawCenteredScaledText(screen, "RECENT CARDS", 286, 170, 720, 2.0, 0)
	if len(g.attempt) == 0 {
		if g.missedCard != nil {
			// After a miss, keep showing the failing card until the next successful draw.
			key := g.missedCard.Name + "_of_" + g.missedCard.Suit
			img, ok := g.sprites.Faces[key]
			if ok && img != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(6, 6)
				op.GeoM.Translate(584, 260)
				screen.DrawImage(img, op)
				return
			}
		}
		drawCenteredScaledText(screen, "NO CARDS DRAWN YET", 286, 300, 720, 2.0, 0)
		return
	}

	// Lay out the recent cards in draw order from left to right.
	for i, c := range g.attempt {
		key := c.Name + "_of_" + c.Suit
		img, ok := g.sprites.Faces[key]
		if !ok || img == nil {
			continue
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(6, 6)
		op.GeoM.Translate(float64(x+i*118), float64(y))
		screen.DrawImage(img, op)
	}
}

// drawButton draws a single interactive button.
func (g *Game) drawButton(screen *ebiten.Image, b Button) {
	// Draw a simple shadow first so the button stands out from the felt background.
	ebitenutil.DrawRect(screen, float64(b.X-2), float64(b.Y-2), float64(b.W+4), float64(b.H+4), color.RGBA{0, 0, 0, 180})
	ebitenutil.DrawRect(screen, float64(b.X), float64(b.Y), float64(b.W), float64(b.H), b.Fill)
	ebitenutil.DebugPrintAt(screen, b.Label, b.X+20, b.Y+20)
}

// drawScaledDebugText draws bitmap debug text with a scale transform applied.
func drawScaledDebugText(screen *ebiten.Image, s string, x, y int, scale float64) {
	w := len(s) * 6
	if w <= 0 {
		return
	}

	// Render into a temporary image because DebugPrint itself does not support scaling.
	tmp := ebiten.NewImage(w, 16)
	ebitenutil.DebugPrintAt(tmp, s, 0, 0)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(tmp, op)
}

// drawCenteredScaledText centers scaled text within a horizontal region.
func drawCenteredScaledText(screen *ebiten.Image, s string, x, y, w int, scale float64, yOffset int) {
	// Approximate text width from the debug font's fixed character size.
	textW := int(float64(len(s)*6) * scale)
	drawX := x + (w-textW)/2
	drawScaledDebugText(screen, s, drawX, y+yOffset, scale)
}

// drawCenteredQuestion wraps and centers the current question in the footer area.
func drawCenteredQuestion(screen *ebiten.Image, question string, x, y, w int) {
	// Wrap long prompts into two lines so they fit the footer panel cleanly.
	lines := wrapQuestion(question, 56)
	startY := y + 8
	if len(lines) == 2 {
		startY = y
	}
	for i, line := range lines {
		// Center each line independently because their lengths can differ.
		lineW := int(float64(len(line)*6) * 2.0)
		lineX := x + (w-lineW)/2
		drawScaledDebugText(screen, line, lineX, startY+(i*34), 2.0)
	}
}

// wrapQuestion splits a question into at most two lines near the given limit.
func wrapQuestion(s string, limit int) []string {
	if len(s) <= limit {
		return []string{s}
	}

	// Search backward for whitespace so the wrap point stays between words.
	split := -1
	for i := limit; i >= 0; i-- {
		if s[i] == ' ' {
			split = i
			break
		}
	}
	if split == -1 {
		return []string{s}
	}

	// Return at most two lines; longer text is left to the caller to avoid overcomplicating layout.
	return []string{s[:split], s[split+1:]}
}

// wrapText breaks a paragraph into lines that fit within the given character limit.
func wrapText(s string, limit int) []string {
	if len(s) <= limit {
		return []string{s}
	}

	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}

	lines := make([]string, 0, len(words))
	current := words[0]
	for _, word := range words[1:] {
		if len(current)+1+len(word) <= limit {
			current += " " + word
			continue
		}
		lines = append(lines, current)
		current = word
	}
	lines = append(lines, current)
	return lines
}

// Layout reports the logical screen size used by Ebiten.
func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

// main configures the game window and starts the Ebiten loop.
func main() {
	// Configure the native window before handing control to Ebiten's game loop.
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ride the Bus")
	if err := ebiten.RunGame(NewGame()); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
