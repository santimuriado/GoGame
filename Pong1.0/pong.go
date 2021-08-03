/* Version 1.0 */

/*Games usually work with an EVENT LOOP.
  for {
	   recieve input(mouse/keyboard etc)
	   update everything (physics, ai)
	   draw
   }
*/

/* THINGS TO IMPROVE:
FRAME RATE INDEPENDENCE.
SCOREBOARD IMPLEMENTATION.
IMPROVE AI.
*/

package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const WIDTH, HEIGHT, PIXELFORMAT int = 800, 600, 4
const HALFWIDTH, HALFHEIGHT int = WIDTH / 2, HEIGHT / 2

type color struct {
	r, g, b byte
}

//In game development it's better to use float32 as it's smaller.
type position struct {
	x, y float32
}

type ball struct {
	position
	radius int
	vx     float32
	vy     float32
	color  color
}

type paddle struct {
	position
	width  int
	height int
	color  color
}

/* To draw the ball it iterates over a square and then it checks if the pixel is inside the
   radius. */
func (ball *ball) draw(pixels []byte) {

	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			//Using squared it's heavier on the CPU.
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x)+x, int(ball.y)+y, ball.color, pixels)
			}
		}
	}
}

func (ball *ball) update(leftPaddle *paddle, rightPaddle *paddle) {

	ball.x += ball.vx
	ball.y += ball.vy

	if int(ball.y)-ball.radius < 0 || int(ball.y)+ball.radius > HEIGHT {
		ball.vy = -ball.vy
	}
	if ball.x < 0 || int(ball.x) > WIDTH {
		ball.x = float32(HALFWIDTH)
		ball.y = float32(HALFHEIGHT)
	}

	if int(ball.x)-ball.radius < int(leftPaddle.x)+rightPaddle.width/2 {
		if int(ball.y) > int(leftPaddle.y)-leftPaddle.height/2 && int(ball.y) < int(leftPaddle.y)+leftPaddle.height/2 {
			ball.vx = -ball.vx
		}
	}
	if int(ball.x)+ball.radius > int(rightPaddle.x)-rightPaddle.width/2 {
		if int(ball.y) > int(rightPaddle.y)-rightPaddle.height/2 && int(ball.y) < int(rightPaddle.y)+rightPaddle.height/2 {
			ball.vx = -ball.vx
		}
	}

}

//Always draw from top left to bottom right.
func (paddle *paddle) draw(pixels []byte) {

	initialX := int(paddle.x) - paddle.width/2
	initialY := int(paddle.y) - paddle.height/2

	for y := 0; y < paddle.height; y++ {
		for x := 0; x < paddle.width; x++ {
			setPixel(initialX+x, initialY+y, paddle.color, pixels)
		}
	}

}

//Input for arrow keys.
func (paddle *paddle) update(keyState []uint8) {

	if keyState[sdl.SCANCODE_UP] != 0 {
		if int(paddle.y) > paddle.height/2 {
			paddle.y -= 5
		}
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		if int(paddle.y) < HEIGHT-paddle.height/2 {
			paddle.y += 5
		}
	}
}

//Simple "AI" to test.
func (paddle *paddle) aiUpdate(ball *ball) {
	paddle.y = ball.y
}

func clean(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c color, pixels []byte) {

	index := (y*WIDTH + x) * PIXELFORMAT

	if index < len(pixels)-PIXELFORMAT && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}

}

func errorCheck(err error) {

	if err != nil {
		fmt.Println(err)
		return
	}
}

func extras(tex *sdl.Texture, renderer *sdl.Renderer, pixels []byte) {

	tex.Update(nil, pixels, WIDTH*PIXELFORMAT)
	renderer.Copy(tex, nil, nil)
	renderer.Present()
}

func destroy(tex *sdl.Texture, renderer *sdl.Renderer, window *sdl.Window) {
	window.Destroy()
	renderer.Destroy()
	tex.Destroy()
	sdl.Quit()
}

func main() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	errorCheck(err)

	window, err := sdl.CreateWindow("PONG 1.0", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(WIDTH), int32(HEIGHT), sdl.WINDOW_SHOWN)
	errorCheck(err)

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	errorCheck(err)

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(WIDTH), int32(HEIGHT))
	errorCheck(err)

	pixels := make([]byte, WIDTH*HEIGHT*PIXELFORMAT)

	defer destroy(tex, renderer, window)

	player1 := paddle{position{50, 100}, 20, 100, color{255, 0, 0}}
	player2 := paddle{position{float32(WIDTH) - 50, 100}, 20, 100, color{0, 255, 0}}
	ball := ball{position{float32(HALFWIDTH), float32(HALFHEIGHT)}, 20, 5, 5, color{0, 0, 255}}

	keyState := sdl.GetKeyboardState()

	//EVENT LOOP
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		clean(pixels)

		player1.draw(pixels)
		player1.update(keyState)
		player2.aiUpdate(&ball)
		player2.draw(pixels)

		ball.draw(pixels)
		ball.update(&player1, &player2)

		extras(tex, renderer, pixels)
		sdl.Delay(16)
	}

}
