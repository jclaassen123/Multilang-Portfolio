package main

import (
	"fmt"
	"image/color"
	"log"

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

	mouseWasDown bool
	quit         bool
}

func NewGame() *Game {
	sprites, err := loadCardSprites("assets/cards.png")
	if err != nil {
		log.Fatal(err)
	}

	g := &Game{sprites: sprites}
	g.startNewMatch()
	return g
}

func (g *Game) startNewMatch() {
	g.deck = NewDeck()
	g.deck.Shuffle()
	g.state = StatePlaying
	g.round = RoundColor
	g.firstCard = Card{}
	g.secondCard = Card{}
	g.missedCard = nil
	g.attempt = nil
	g.message = "Get 4 correct in a row before the deck runs out."
	g.setupRoundButtons()
}

func (g *Game) setupRoundButtons() {
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

func (g *Game) setupPlayAgainButtons() {
	g.question = "Play again?"
	g.buttons = []Button{
		{X: 360, Y: 642, W: 140, H: 58, Label: "Yes", Fill: color.RGBA{46, 134, 78, 255}, Action: g.startNewMatch},
		{X: 520, Y: 642, W: 140, H: 58, Label: "No", Fill: color.RGBA{121, 44, 44, 255}, Action: func() { g.quit = true }},
	}
}

func (g *Game) drawCard() (Card, bool) {
	card, ok := g.deck.Draw()
	if !ok {
		g.state = StateLost
		g.message = "Out of cards. You lose."
		g.setupPlayAgainButtons()
		return Card{}, false
	}
	g.attempt = append(g.attempt, card)
	if len(g.attempt) > 4 {
		g.attempt = g.attempt[len(g.attempt)-4:]
	}
	return card, true
}

func (g *Game) resetToRoundOne(wrongCard Card) {
	g.round = RoundColor
	g.firstCard = Card{}
	g.secondCard = Card{}
	g.missedCard = &wrongCard
	g.attempt = nil
	g.message = fmt.Sprintf("Wrong: %s of %s. Back to Round 1.", wrongCard.Name, wrongCard.Suit)
	g.setupRoundButtons()
}

func (g *Game) guessColor(guessRed bool) {
	card, ok := g.drawCard()
	if !ok {
		return
	}
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

func (g *Game) guessHigherLower(guessHigher bool) {
	card, ok := g.drawCard()
	if !ok {
		return
	}

	correct := false
	if guessHigher {
		correct = card.Value > g.firstCard.Value
	} else {
		correct = card.Value < g.firstCard.Value
	}

	if correct {
		g.missedCard = nil
		g.secondCard = card
		g.round = RoundInsideOutside
		g.message = fmt.Sprintf("Correct: %s of %s.", card.Name, card.Suit)
		g.setupRoundButtons()
		return
	}
	g.resetToRoundOne(card)
}

func (g *Game) guessInsideOutside(guessInside bool) {
	card, ok := g.drawCard()
	if !ok {
		return
	}

	low := min(g.firstCard.Value, g.secondCard.Value)
	high := max(g.firstCard.Value, g.secondCard.Value)
	isInside := card.Value > low && card.Value < high
	isOutside := card.Value < low || card.Value > high
	correct := (guessInside && isInside) || (!guessInside && isOutside)

	if correct {
		g.missedCard = nil
		g.round = RoundSuit
		g.message = fmt.Sprintf("Correct: %s of %s.", card.Name, card.Suit)
		g.setupRoundButtons()
		return
	}
	g.resetToRoundOne(card)
}

func (g *Game) guessSuit(suit string) {
	card, ok := g.drawCard()
	if !ok {
		return
	}

	if card.Suit == suit {
		g.missedCard = nil
		g.state = StateWon
		g.message = fmt.Sprintf("You won with %s of %s!", card.Name, card.Suit)
		g.setupPlayAgainButtons()
		return
	}
	g.resetToRoundOne(card)
}

func (g *Game) Update() error {
	if g.quit {
		return ebiten.Termination
	}

	isDown := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if isDown && !g.mouseWasDown {
		mx, my := ebiten.CursorPosition()
		for _, b := range g.buttons {
			if mx >= b.X && mx <= b.X+b.W && my >= b.Y && my <= b.Y+b.H {
				b.Action()
				break
			}
		}
	}
	g.mouseWasDown = isDown
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{14, 82, 42, 255})

	ebitenutil.DrawRect(screen, 14, 14, 996, 692, color.RGBA{20, 98, 49, 255})
	ebitenutil.DrawRect(screen, 18, 18, 988, 82, color.RGBA{10, 48, 24, 255})
	ebitenutil.DrawRect(screen, 18, 108, 260, 450, color.RGBA{17, 68, 35, 255})
	ebitenutil.DrawRect(screen, 286, 108, 720, 450, color.RGBA{17, 68, 35, 255})
	ebitenutil.DrawRect(screen, 18, 566, 988, 140, color.RGBA{10, 48, 24, 255})

	drawCenteredScaledText(screen, "RIDE THE BUS", 18, 18, 988, 3.0, 18)

	g.drawDeckStack(screen, 82, 176)
	drawCenteredScaledText(screen, fmt.Sprintf("Cards Left: %d", len(g.deck.Cards)), 18, 462, 260, 2.0, 0)
	if g.state == StatePlaying {
		drawCenteredScaledText(screen, fmt.Sprintf("Round: %d / 4", g.round), 18, 500, 260, 2.0, 0)
	}

	g.drawHistory(screen, 326, 290)
	drawCenteredScaledText(screen, g.message, 286, 520, 720, 2.0, 0)

	drawCenteredQuestion(screen, g.question, 18, 566, 988)

	for _, b := range g.buttons {
		g.drawButton(screen, b)
	}
}

func (g *Game) drawDeckStack(screen *ebiten.Image, x, y int) {
	if g.sprites.Back == nil {
		return
	}

	count := len(g.deck.Cards)
	layers := count / 5
	if layers < 1 && count > 0 {
		layers = 1
	}
	if layers > 10 {
		layers = 10
	}

	scale := 6.0
	for i := 0; i < layers; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(x+i*3), float64(y-i*2))
		screen.DrawImage(g.sprites.Back, op)
	}
	ebitenutil.DebugPrintAt(screen, "Deck", x+44, y+130)
}

func (g *Game) drawHistory(screen *ebiten.Image, x, y int) {
	drawCenteredScaledText(screen, "RECENT CARDS", 286, 170, 720, 2.0, 0)
	if len(g.attempt) == 0 {
		if g.missedCard != nil {
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

func (g *Game) drawButton(screen *ebiten.Image, b Button) {
	ebitenutil.DrawRect(screen, float64(b.X-2), float64(b.Y-2), float64(b.W+4), float64(b.H+4), color.RGBA{0, 0, 0, 180})
	ebitenutil.DrawRect(screen, float64(b.X), float64(b.Y), float64(b.W), float64(b.H), b.Fill)
	ebitenutil.DebugPrintAt(screen, b.Label, b.X+20, b.Y+20)
}

func drawScaledDebugText(screen *ebiten.Image, s string, x, y int, scale float64) {
	w := len(s) * 6
	if w <= 0 {
		return
	}
	tmp := ebiten.NewImage(w, 16)
	ebitenutil.DebugPrintAt(tmp, s, 0, 0)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(tmp, op)
}

func drawCenteredScaledText(screen *ebiten.Image, s string, x, y, w int, scale float64, yOffset int) {
	textW := int(float64(len(s)*6) * scale)
	drawX := x + (w-textW)/2
	drawScaledDebugText(screen, s, drawX, y+yOffset, scale)
}

func drawCenteredQuestion(screen *ebiten.Image, question string, x, y, w int) {
	lines := wrapQuestion(question, 56)
	startY := y + 8
	if len(lines) == 2 {
		startY = y
	}
	for i, line := range lines {
		lineW := int(float64(len(line)*6) * 2.0)
		lineX := x + (w-lineW)/2
		drawScaledDebugText(screen, line, lineX, startY+(i*34), 2.0)
	}
}

func wrapQuestion(s string, limit int) []string {
	if len(s) <= limit {
		return []string{s}
	}
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
	return []string{s[:split], s[split+1:]}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Ride the Bus")
	if err := ebiten.RunGame(NewGame()); err != nil && err != ebiten.Termination {
		log.Fatal(err)
	}
}
