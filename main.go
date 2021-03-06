package main

import (
	"image/color"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	bgColor = color.RGBA{0, 0, 0, 0}
	lifeColor = color.RGBA{60, 200, 40, 255}
}

// World represents the game state.
type World struct {
	area   []color.RGBA
	width  int
	height int
}

// NewWorld creates a new world.
func NewWorld(width, height int, maxInitLiveCells int) *World {
	w := &World{
		area:   make([]color.RGBA, width*height),
		width:  width,
		height: height,
	}
	w.init(maxInitLiveCells)
	return w
}

// init inits world with a random state.
func (w *World) init(maxLiveCells int) {
	// for y := 0; y < w.height; y++ {
	// 	for x := 0; x < w.width; x++ {
	// 		w.area[y*w.width+x] = bgColor
	// 	}
	// }
	for i := 0; i < maxLiveCells; i++ {
		x := rand.Intn(w.width)
		y := rand.Intn(w.height)
		w.area[y*w.width+x] = lifeColor
	}
}

var stepTimeStamp = time.Now()

// Update game state by one tick.
func (w *World) Update() {
	if time.Since(stepTimeStamp) < 10*time.Millisecond {
		return
	}
	stepTimeStamp = time.Now()
	width := w.width
	height := w.height
	next := make([]color.RGBA, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pop, emergingColor := neighbourCount(w.area, width, height, x, y)
			switch {
			case pop < 2:
				// rule 1. Any live cell with fewer than two live neighbours
				// dies, as if caused by under-population.
				next[y*width+x] = bgColor

			case (pop == 2 || pop == 3) && (w.area[y*width+x] != bgColor):
				// rule 2. Any live cell with two or three live neighbours
				// lives on to the next generation.
				next[y*width+x] = emergingColor

			case pop > 3:
				// rule 3. Any live cell with more than three live neighbours
				// dies, as if by over-population.
				next[y*width+x] = bgColor

			case pop == 3:
				// rule 4. Any dead cell with exactly three live neighbours
				// becomes a live cell, as if by reproduction.
				next[y*width+x] = emergingColor
			}
		}
	}

	if currentWrench.status == fresh {
		currentWrench.Lock()
		currentWrench.status = running
		currentWrench.Unlock()
	}
	if currentWrench.progress >= screenWidth && currentWrench.textPointer >= len(currentWrench.text) {
		currentWrench.Lock()
		currentWrench.status = stopped
		currentWrench.progress = 0
		currentWrench.textPointer = 0
		currentWrench.Unlock()
	}
	if currentWrench.status == running {
		if currentWrench.progress >= screenWidth {
			currentWrench.Lock()
			currentWrench.progress = 0
			currentWrench.textPointer += 1
			currentWrench.Unlock()
		}
		currentWrench.Lock()
		for h := 0; h < currentWrench.boxHeight; h++ {
			for w := 0; w < currentWrench.boxWidth; w++ {
				if (currentWrench.y+h)*width+currentWrench.progress+w < len(next) {
					next[(currentWrench.y+h)*width+currentWrench.progress+w] = currentWrench.color
				}
			}
		}
		//log.Printf("%d -> %d \n", currentWrench.y, currentWrench.y*width+currentWrench.progress)
		currentWrench.progress += currentWrench.stepSize
		currentWrench.Unlock()
	}

	w.area = next
}

// Draw paints current game state.
func (w *World) Draw(pix []byte) {
	for i, v := range w.area {
		pix[4*i] = v.R
		pix[4*i+1] = v.G
		pix[4*i+2] = v.B
		pix[4*i+3] = v.A
	}
}

const (
	screenWidth  = 640 // 640 //320
	screenHeight = 360 // 480 //240
)

type Game struct {
	world  *World
	pixels []byte
}

func (g *Game) Update() error {
	g.world.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, screenWidth*screenHeight*4)
	}
	g.world.Draw(g.pixels)
	screen.ReplacePixels(g.pixels)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

var bgColor color.RGBA
var lifeColor color.RGBA

type status int

const (
	fresh   status = 0
	running status = 1
	stopped status = 2
)

type wrench struct {
	sync.Mutex
	y           int
	progress    int
	status      status
	stepSize    int
	boxHeight   int
	boxWidth    int
	text        string
	textPointer int
	color       color.RGBA
}

//var chaosChan = make(chan wrench)
var currentWrench = wrench{
	status:    fresh,
	stepSize:  1,
	boxHeight: 15,
	boxWidth:  1,
	text:      "", //"It's alive",
	color:     color.RGBA{250, 50, 0, 255},
}

func main() {
	go func() {
		for {
			//newWrench := wrench{y: rand.Intn(screenHeight)}
			//chaosChan <- newWrench
			if currentWrench.status != running {
				currentWrench.Lock()
				currentWrench.y = rand.Intn(screenHeight)
				currentWrench.boxHeight = rand.Intn(25) + 1
				currentWrench.boxWidth = rand.Intn(10) + 1
				currentWrench.status = fresh
				currentWrench.Unlock()
				log.Printf("???? %d\n", currentWrench.y)
			}
			<-time.After(5 * time.Second)
		}
	}()
	g := &Game{
		//world: NewWorld(screenWidth, screenHeight, int((screenWidth*screenHeight)/10)),
		world: NewWorld(screenWidth, screenHeight, int(1000)),
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Game of Life (Ebiten Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
