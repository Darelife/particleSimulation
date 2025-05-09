package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidthPx  = 800
	screenHeightPx = 600
	particleCount  = 2000
	particleRadius = 5
	gravityY       = 0.098
	gravityX       = -0.005
)

type Particle struct {
	X, Y       float64
	VelX, VelY float64
	Color      [3]uint8
}

type Game struct {
	particles    []Particle
	showMenu     bool
	numParticles int
	boostButton  struct {
		x, y, width, height int
	}
	upBoostButton struct {
		x, y, width, height int
	}
	particlesToAdd int
	frameCount     int
}

func (g *Game) Update() error {
	g.frameCount++
	mouseX, mouseY := ebiten.CursorPosition()
	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	if ebiten.IsKeyPressed(ebiten.KeyM) && ebiten.IsKeyPressed(ebiten.KeyControl) {
		g.showMenu = !g.showMenu
	}

	if g.particlesToAdd > 0 && g.frameCount > 0 {
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))

		startX := float64(particleRadius * 2)
		startY := float64(particleRadius * 2)

		// angle := math.Pi / 36
		angle := rng.Float64() * (2 * math.Pi)
		baseAngle := 0.0
		speed := 10.00

		newParticle := Particle{
			X:    startX + rng.Float64()*particleRadius,
			Y:    startY + rng.Float64()*particleRadius,
			VelX: speed * math.Cos(baseAngle+angle),
			VelY: speed*math.Sin(baseAngle+angle) + rng.Float64()*0.5,
			Color: [3]uint8{
				uint8(50 + rng.Intn(50)),
				uint8(50 + rng.Intn(50)),
				uint8(200 + rng.Intn(55)),
			},
		}

		g.particles = append(g.particles, newParticle)
		g.particlesToAdd--
	}

	if mousePressed {
		if mouseX >= g.boostButton.x && mouseX <= g.boostButton.x+g.boostButton.width &&
			mouseY >= g.boostButton.y && mouseY <= g.boostButton.y+g.boostButton.height {
			if len(g.particles)+g.particlesToAdd < 1000 {
				g.particlesToAdd += 20
			}
		}

		if mouseX >= g.upBoostButton.x && mouseX <= g.upBoostButton.x+g.upBoostButton.width &&
			mouseY >= g.upBoostButton.y && mouseY <= g.upBoostButton.y+g.upBoostButton.height {
			for i := range g.particles {
				g.particles[i].VelY -= 10.0
			}
		}

		if g.showMenu {
			menuX := 20
			menuY := 100
			buttonHeight := 30
			buttonWidth := 200

			if mouseX >= menuX && mouseX <= menuX+buttonWidth {
				if mouseY >= menuY && mouseY <= menuY+buttonHeight {
					g.particlesToAdd += 100
					g.numParticles = 100
				} else if mouseY >= menuY+buttonHeight+10 && mouseY <= menuY+buttonHeight*2+10 {
					g.particlesToAdd += 300
					g.numParticles = 300
				} else if mouseY >= menuY+buttonHeight*2+20 && mouseY <= menuY+buttonHeight*3+20 {
					g.particlesToAdd += 500
					g.numParticles = 500
				}
			}
		}
	}

	for i := range g.particles {
		p := &g.particles[i]

		p.VelY += gravityY
		p.VelX += gravityX

		p.X += p.VelX
		p.Y += p.VelY

		if p.X < particleRadius {
			p.X = particleRadius
			p.VelX = -p.VelX * 0.3
		} else if p.X > screenWidthPx-particleRadius {
			p.X = screenWidthPx - particleRadius
			p.VelX = -p.VelX * 0.3
		}

		if p.Y < particleRadius {
			p.Y = particleRadius
			p.VelY = -p.VelY * 0.3
		} else if p.Y > screenHeightPx-particleRadius {
			p.Y = screenHeightPx - particleRadius
			p.VelY = -p.VelY * 0.3
		}
	}

	const collisionIterations = 3
	for iter := 0; iter < collisionIterations; iter++ {
		for i := range g.particles {
			for j := i + 1; j < len(g.particles); j++ {
				p1 := &g.particles[i]
				p2 := &g.particles[j]

				dx := p1.X - p2.X
				dy := p1.Y - p2.Y
				distSquared := dx*dx + dy*dy

				minDist := 2.0 * float64(particleRadius)
				if distSquared < minDist*minDist {
					dist := math.Sqrt(distSquared)

					nx := dx / dist
					ny := dy / dist

					overlap := minDist - dist

					moveX := nx * overlap * 0.5
					moveY := ny * overlap * 0.5

					p1.X += moveX
					p1.Y += moveY
					p2.X -= moveX
					p2.Y -= moveY

					v1n := p1.VelX*nx + p1.VelY*ny
					v2n := p2.VelX*nx + p2.VelY*ny

					const restitution = 0.9
					v1nNew := (v1n*(1-restitution) + v2n*(1+restitution)) / 2
					v2nNew := (v2n*(1-restitution) + v1n*(1+restitution)) / 2

					dvx1 := (v1nNew - v1n) * nx
					dvy1 := (v1nNew - v1n) * ny
					dvx2 := (v2nNew - v2n) * nx
					dvy2 := (v2nNew - v2n) * ny

					p1.VelX += dvx1
					p1.VelY += dvy1
					p2.VelX += dvx2
					p2.VelY += dvy2
				}
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.particles {
		vector.DrawFilledCircle(
			screen,
			float32(p.X), float32(p.Y),
			particleRadius,
			createRGBA(p.Color[0], p.Color[1], p.Color[2], 255),
			false,
		)
	}

	g.boostButton.x = screenWidthPx/2 - 120
	g.boostButton.y = 20
	g.boostButton.width = 100
	g.boostButton.height = 40

	g.upBoostButton.x = screenWidthPx/2 + 20
	g.upBoostButton.y = 20
	g.upBoostButton.width = 100
	g.upBoostButton.height = 40

	vector.DrawFilledRect(
		screen,
		float32(g.boostButton.x),
		float32(g.boostButton.y),
		float32(g.boostButton.width),
		float32(g.boostButton.height),
		createRGBA(0, 150, 255, 255),
		false,
	)

	vector.DrawFilledRect(
		screen,
		float32(g.upBoostButton.x),
		float32(g.upBoostButton.y),
		float32(g.upBoostButton.width),
		float32(g.upBoostButton.height),
		createRGBA(255, 100, 100, 255),
		false,
	)

	ebitenutil.DebugPrintAt(screen, "ADD WATER", g.boostButton.x+15, g.boostButton.y+15)
	ebitenutil.DebugPrintAt(screen, "BOOST UP", g.upBoostButton.x+20, g.upBoostButton.y+15)

	if g.showMenu {
		menuX := 20
		menuY := 100
		buttonHeight := 30
		buttonWidth := 200

		vector.DrawFilledRect(
			screen,
			float32(menuX-10),
			float32(menuY-10),
			float32(buttonWidth+20),
			float32(buttonHeight*3+40),
			createRGBA(200, 200, 200, 200),
			false,
		)

		vector.DrawFilledRect(
			screen,
			float32(menuX),
			float32(menuY),
			float32(buttonWidth),
			float32(buttonHeight),
			createRGBA(100, 100, 255, 255),
			false,
		)
		ebitenutil.DebugPrintAt(screen, "Add 100 Particles", menuX+10, menuY+10)

		vector.DrawFilledRect(
			screen,
			float32(menuX),
			float32(menuY+buttonHeight+10),
			float32(buttonWidth),
			float32(buttonHeight),
			createRGBA(100, 100, 255, 255),
			false,
		)
		ebitenutil.DebugPrintAt(screen, "Add 300 Particles", menuX+10, menuY+buttonHeight+10+10)

		vector.DrawFilledRect(
			screen,
			float32(menuX),
			float32(menuY+buttonHeight*2+20),
			float32(buttonWidth),
			float32(buttonHeight),
			createRGBA(100, 100, 255, 255),
			false,
		)
		ebitenutil.DebugPrintAt(screen, "Add 500 Particles", menuX+10, menuY+buttonHeight*2+20+10)
	}

	ebitenutil.DebugPrintAt(screen, "FPS: "+fmt.Sprintf("%.2f", ebiten.ActualFPS()), 10, 10)
	ebitenutil.DebugPrintAt(screen, "Press Ctrl+M to toggle menu", 10, 25)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Particles: %d", len(g.particles)), 10, 40)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidthPx, screenHeightPx
}

func createRGBA(r, g, b, a uint8) color.RGBA {
	return color.RGBA{R: r, G: g, B: b, A: a}
}

// func initializeParticles(count int) []Particle {
// 	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	particles := make([]Particle, count)

// 	startX := float64(particleRadius * 2)
// 	startY := float64(particleRadius * 2)

// 	for i := range particles {
// 		angle := rng.Float64() * (2 * math.Pi)

// 		speed := 2.0 + rng.Float64()*3.0

// 		particles[i] = Particle{
// 			X:    startX,
// 			Y:    startY,
// 			VelX: speed * math.Cos(angle),
// 			VelY: speed * math.Sin(angle),
// 			Color: [3]uint8{
// 				uint8(rng.Intn(256)),
// 				uint8(rng.Intn(256)),
// 				uint8(rng.Intn(256)),
// 			},
// 		}
// 	}

// 	return particles
// }

func main() {
	game := &Game{
		particles:      []Particle{},
		numParticles:   particleCount,
		particlesToAdd: particleCount,
	}

	game.boostButton.x = screenWidthPx/2 - 120
	game.boostButton.y = 20
	game.boostButton.width = 100
	game.boostButton.height = 40

	game.upBoostButton.x = screenWidthPx/2 + 20
	game.upBoostButton.y = 20
	game.upBoostButton.width = 100
	game.upBoostButton.height = 40

	ebiten.SetWindowSize(screenWidthPx, screenHeightPx)
	ebiten.SetWindowTitle("Particle Simulator")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
