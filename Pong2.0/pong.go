/* Version 2.0 */

/* FRAME RATE INDEPENDENCE. DONE.
   SCOREBOARD. DONE.
   MAKING THE BALL FASTER WITH EVERY FRAME. DONE.
   IMPROVE AI SO IT MAKES MISTAKES SOMETIMES. DONE.
*/

package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const WIDTH, HEIGHT, PIXELFORMAT int = 800, 600, 4
const HALFWIDTH, HALFHEIGHT int = WIDTH / 2, HEIGHT / 2
const XVELOCITY, YVELOCITY float32 = 600, 600

const PADDLEWIDTH, PADDLEHEIGHT int = 20, 100
const PADDLESPEED int = 500

//Change ball radius at your own will but big values may cause unexpected behavior.
const BALLRADIUS int = 20

/*THE HIGHER IT IS THE EASIEST IT BECOMES.
  RECOMMENDED VALUE IS 3.5. */
const DIFFICULTYFACTOR float32 = 3.5

const MAXIMUMSCORE int = 3

type color struct {
	r, g, b byte
}

type position struct {
	x, y float32
}

type ball struct {
	position
	radius float32
	vx     float32
	vy     float32
	color  color
}

type paddle struct {
	position
	width    float32
	height   float32
	velocity float32
	points   int
	color    color
}

//The 1's make up the number.
var numbers = [][]byte{
	//0.
	{
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
	},
	//1.
	{
		1, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1,
	},
	//2.
	{
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
	},
	//3.
	{
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
	}}

type gameState int

const (
	beginning gameState = iota
	playing
)

var state = beginning

//Paint pixels for numbers in the scoreboard.
func drawNumbers(pos position, color color, size int, number int, pixels []byte) {

	initialX := int(pos.x) - (size*3)/2
	initialY := int(pos.y) - (size*5)/2

	for i, v := range numbers[number] {
		if v == 1 {
			for y := initialY; y < initialY+size; y++ {
				for x := initialX; x < initialX+size; x++ {
					setPixel(x, y, color, pixels)
				}
			}
		}
		initialX += size
		if (i+1)%3 == 0 {
			initialY += size
			initialX -= size * 3
		}
	}

}

//Draw ball.Method explained in Pong 1.0.
func (ball *ball) draw(pixels []byte) {

	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x+x), int(ball.y+y), ball.color, pixels)
			}
		}
	}
}

//Returns position struct with the middle of the screen.
func getCentre() position {
	return position{float32(HALFWIDTH), float32(HALFHEIGHT)}
}

//Centres both paddles in the Y axis.
func centrePaddles(leftPaddle *paddle, rightPaddle *paddle) {

	leftPaddle.y = float32(HALFHEIGHT)
	rightPaddle.y = float32(HALFHEIGHT)
}

//Resets ball and paddles position. Also balls original velocity.
func reset(winningPaddle *paddle, losingPaddle *paddle, ball *ball) {
	winningPaddle.points++
	ball.position = getCentre()
	ball.vx = XVELOCITY
	ball.vy = YVELOCITY
	centrePaddles(winningPaddle, losingPaddle)
	state = beginning
}

//Changes ball velocity depending on the direction it's going. Used to speed up the game.
func (ball *ball) changeVelocity() {

	if ball.vy < 0 {
		ball.vy -= 0.1
	}
	if ball.vy > 0 {
		ball.vy += 0.1
	}
	if ball.vx < 0 {
		ball.vx -= 0.1
	}
	if ball.vx > 0 {
		ball.vx += 0.1
	}
}

//Changes ball direction when it hits something.
func (ball *ball) update(leftPaddle *paddle, rightPaddle *paddle, duration float32) {

	ball.x += ball.vx * duration
	ball.y += ball.vy * duration

	//Radius of the ball has to be taken into account so that it doesn't clip the limits or the paddles.
	//Changes ball velocity whenever the ball touches the ceiling or the floor.
	if ball.y-ball.radius < 0 || ball.y+ball.radius > float32(HEIGHT) {
		ball.vy = -ball.vy
	}

	//Changes ball velocity when it hits either paddle.
	if ball.x-ball.radius < leftPaddle.x+leftPaddle.width/2 {
		if ball.y > leftPaddle.y-leftPaddle.height/2 && ball.y < leftPaddle.y+leftPaddle.height/2 {
			ball.vx = -ball.vx
			ball.x = leftPaddle.x + leftPaddle.width/2.0 + ball.radius
		}
	}
	if ball.x+ball.radius > rightPaddle.x-rightPaddle.width/2 {
		if ball.y > rightPaddle.y-rightPaddle.height/2 && ball.y < rightPaddle.y+rightPaddle.height/2 {
			ball.vx = -ball.vx
			ball.x = rightPaddle.x - rightPaddle.width/2.0 - ball.radius

		}
	}

	//Resets paddles and ball position if someone scores.
	if ball.x < 0 {
		reset(rightPaddle, leftPaddle, ball)

	} else if ball.x > float32(WIDTH) {
		reset(leftPaddle, rightPaddle, ball)
	}

	//Speeds up the ball every frame.
	ball.changeVelocity()

}

//Lerp stands for linear interpolation.
//Used to get the value of X coordinate when drawing the scoreboard numbers.
func lerp(a float32, b float32, percentage float32) float32 {
	return a + percentage*(b-a)
}

//Drawing of the paddle.
func (paddle *paddle) draw(pixels []byte) {

	initialX := int(paddle.x - paddle.width/2)
	initialY := int(paddle.y - paddle.height/2)

	for y := 0; y < int(paddle.height); y++ {
		for x := 0; x < int(paddle.width); x++ {
			setPixel(initialX+x, initialY+y, paddle.color, pixels)
		}
	}

	xValue := lerp(paddle.x, getCentre().x, 0.2)
	drawNumbers(position{xValue, 35}, paddle.color, 10, paddle.points, pixels)

}

//Updates paddle position according to user input.
func (paddle *paddle) update(keyState []uint8, duration float32, controllerAxis int16) {

	//Moves the paddle up when pressing up arrow key.
	if keyState[sdl.SCANCODE_UP] != 0 {
		if paddle.y > paddle.height/2 {
			paddle.y -= paddle.velocity * duration
		}
	}

	//Moves the paddle down when pressing down arrow key.
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		if paddle.y < float32(HEIGHT)-paddle.height/2 {
			paddle.y += paddle.velocity * duration
		}
	}
	//Recieves input from controller.
	if math.Abs(float64(controllerAxis)) > 1500 {
		percentage := float32(controllerAxis) / 32767.00
		paddle.y += paddle.velocity * percentage * duration
	}

}

//How the opposite paddle moves.
//DIFFICULTYFACTOR changes how much the paddle moves in a random manner.
//Basically it makes the AI make mistakes.
func (paddle *paddle) aiUpdate(ball *ball, duration float32) {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	if ball.vy > 0 && paddle.y < float32(HEIGHT)-paddle.height/2 {
		paddle.y += paddle.velocity * duration * (float32(r1.Intn(10)) / float32(DIFFICULTYFACTOR))
	}
	if ball.vy < 0 && paddle.y > paddle.height/2 {
		paddle.y -= paddle.velocity * duration * (float32(r1.Intn(10)) / float32(DIFFICULTYFACTOR))
	}
}

//Cleans pixels that shouldn't be painted anymore.
func cleanPixels(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

//Paints pixels.
func setPixel(x, y int, c color, pixels []byte) {

	index := (y*WIDTH + x) * PIXELFORMAT

	if index < len(pixels)-PIXELFORMAT && index >= 0 {
		pixels[int(index)] = c.r
		pixels[int(index+1)] = c.g
		pixels[int(index+2)] = c.b
	}

}

func errorCheck(err error) {

	if err != nil {
		fmt.Println(err)
		return
	}
}

//Bundled functions for code cleanness.
func extras(tex *sdl.Texture, renderer *sdl.Renderer, pixels []byte) {

	tex.Update(nil, pixels, int(WIDTH*PIXELFORMAT))
	renderer.Copy(tex, nil, nil)
	renderer.Present()
}

//Bundled destroy functions for code cleanness.
func destroy(tex *sdl.Texture, renderer *sdl.Renderer, window *sdl.Window) {
	window.Destroy()
	renderer.Destroy()
	tex.Destroy()
	sdl.Quit()
}

//Handles the amount of controllers connected.
func controllers(controllerHandler []*sdl.GameController) {

	for i := 0; i < sdl.NumJoysticks(); i++ {
		controllerHandler = append(controllerHandler, sdl.GameControllerOpen(i))
		defer controllerHandler[i].Close()
	}
}

func leftAnalog(controllerHandler []*sdl.GameController, controllerAxis *int16) {
	for _, controller := range controllerHandler {
		if controller != nil {
			*controllerAxis = controller.Axis(sdl.CONTROLLER_AXIS_LEFTY)
		}
	}
}

//Handles the event loop to close the program.
func closeWindow() bool {

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			return true
		}
	}
	return false
}

func main() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	errorCheck(err)

	window, err := sdl.CreateWindow("PONG 2.0", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(WIDTH), int32(HEIGHT), sdl.WINDOW_SHOWN)
	errorCheck(err)

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	errorCheck(err)

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(WIDTH), int32(HEIGHT))
	errorCheck(err)

	defer destroy(tex, renderer, window)

	pixels := make([]byte, int(WIDTH*HEIGHT*PIXELFORMAT))

	var controllerHandler []*sdl.GameController
	controllers(controllerHandler)

	//CREATES PLAYERS AND BALL.
	player1 := paddle{position{50, float32(HALFHEIGHT)}, float32(PADDLEWIDTH), float32(PADDLEHEIGHT),
		float32(PADDLESPEED), 0, color{255, 0, 0}}
	player2 := paddle{position{float32(WIDTH) - 50, float32(HALFHEIGHT)}, float32(PADDLEWIDTH), float32(PADDLEHEIGHT),
		float32(PADDLESPEED), 0, color{0, 255, 0}}
	ball := ball{position{float32(HALFWIDTH), float32(HALFHEIGHT)}, float32(BALLRADIUS), XVELOCITY, YVELOCITY,
		color{0, 0, 255}}

	keyState := sdl.GetKeyboardState()

	var frameStart time.Time
	var duration float32
	var controllerAxis int16

	//EVENT LOOP
	for {
		//Time struct with current time.
		frameStart = time.Now()

		if closeWindow() {
			return
		}
		leftAnalog(controllerHandler, &controllerAxis)

		if state == playing {
			player1.update(keyState, duration, controllerAxis)
			player2.aiUpdate(&ball, duration)
			ball.update(&player1, &player2, duration)

		} else if state == beginning {
			//Waits for user to press spacebar.
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				//Resets scoreboard when someone gets to the MAXIMUMSCORE.
				if player1.points == MAXIMUMSCORE || player2.points == MAXIMUMSCORE {
					player1.points = 0
					player2.points = 0
				}
				state = playing
			}
		}
		cleanPixels(pixels)

		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		extras(tex, renderer, pixels)

		//Frame rate independence.
		duration = float32(time.Since(frameStart).Seconds())
		if duration < 0.005 {
			sdl.Delay(5 - uint32(duration/1000.0))
			duration = float32(time.Since(frameStart).Seconds())
		}
	}

}
