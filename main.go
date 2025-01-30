package main

import (
	"bytes"
	_ "embed"
	"log"
	"strconv"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth = 640
	screenHeight = 480
	ballSpeed = 3
	paddleSpeed = 6
)

type Object struct {
	X, Y, W, H int
}

type Paddle struct {
	Object
}

type Ball struct {
	Object
	VX, VY int // X and Y velocity
}

type Game struct {
	paddle Paddle
	ball Ball
	score int
	highScore int
}

//go:embed PressStart2P-Regular.ttf
var pressStart2P []byte
var pressStart2PFaceSource *text.GoTextFaceSource

func main() {
	ebiten.SetWindowTitle("Simple Pong Game in Ebiten")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	paddle := Paddle{
		Object: Object{
			X: 600,
			Y: 200,
			W: 15,
			H: 100,
		},
	}

	ball := Ball{
		Object: Object{
			X: 0,
			Y: 0,
			W: 15,
			H: 15,
		},
		VX: ballSpeed,
		VY: ballSpeed,
	}

	g := &Game{
		paddle: paddle,
		ball: ball,
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Paddle
	vector.DrawFilledRect(screen,
		float32(g.paddle.X), float32(g.paddle.Y),
		float32(g.paddle.W), float32(g.paddle.H),
		color.White, false,
	)

	// Ball
	vector.DrawFilledRect(screen,
		float32(g.ball.X), float32(g.ball.Y),
		float32(g.ball.W), float32(g.ball.H),
		color.White, false,
	)

	// Text Options
	s, err := text.NewGoTextFaceSource(bytes.NewReader(pressStart2P))
	if err != nil {
		log.Fatal(err)
	}
	pressStart2PFaceSource = s

	scoreTextOptions := &text.DrawOptions{}
	scoreTextOptions.GeoM.Translate(10, 10)
	scoreTextOptions.ColorScale.Scale(1, 1, 1, 1)
	scoreTextOptions.LineSpacing = 1.5

	highScoreTextOptions := &text.DrawOptions{}
	highScoreTextOptions.GeoM.Translate(10, 30)
	highScoreTextOptions.ColorScale.Scale(1, 1, 1, 1)
	highScoreTextOptions.LineSpacing = 1.5

	textFace := &text.GoTextFace{
		Source: pressStart2PFaceSource,
		Size:  13,
	}

	scoreStr := "Score: " + strconv.Itoa(g.score)
	text.Draw(screen , scoreStr, textFace, scoreTextOptions)

	HighScoreStr := "High Score: " + strconv.Itoa(g.highScore)
	text.Draw(screen, HighScoreStr, textFace, highScoreTextOptions)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Update() error {
	g.paddle.MoveOnKeyPress()
	g.ball.Move()
	g.CollideWithWall()
	g.CollideWithPaddle()
	return nil
}

func (p *Paddle) MoveOnKeyPress() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		p.Y -= paddleSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		p.Y += paddleSpeed
	}
}

func (b *Ball) Move() {
	b.X += b.VX
	b.Y += b.VY
}

func (g *Game) Reset() {
	g.ball.X = 0
	g.ball.Y = 0

	g.score = 0
}

func (g *Game) CollideWithWall() {
	// If the ball hits the right wall the game is over
	if g.ball.X > screenWidth {
		g.Reset()
	} else if g.ball.X <= 0 {
		g.ball.VX = ballSpeed
	} else if g.ball.Y <= 0 {
		g.ball.VY = ballSpeed
	} else if g.ball.Y >= screenHeight {
		g.ball.VY = -ballSpeed
	}
}

func (g *Game) CollideWithPaddle() {
	if g.ball.X >= g.paddle.X && g.ball.Y >= g.paddle.Y && g.ball.Y <= g.paddle.Y + g.paddle.H {
		g.ball.VX = -g.ball.VX
		g.score++

		if g.score > g.highScore {
			g.highScore = g.score
		}
	}
}
