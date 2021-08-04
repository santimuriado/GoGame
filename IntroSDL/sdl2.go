/* PAINTS A GENERIC PICTURE TO TEST SDL2 */

package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const WIDTH, HEIGHT, PIXELFORMAT int = 800, 600, 4

type color struct {
	r, g, b byte
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

func paintPixels(pixels []byte) {

	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			setPixel(x, y, color{byte(x % 255), byte(y % 255), 0}, pixels)
		}
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
}

func main() {

	window, err := sdl.CreateWindow(
		"Testing SDL2",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(WIDTH),
		int32(HEIGHT),
		sdl.WINDOW_SHOWN,
	)
	errorCheck(err)

	//With a GPU the renderer is ACCELERATED if not it's SOFTWARE.
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	errorCheck(err)

	tex, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		int32(WIDTH),
		int32(HEIGHT),
	)
	errorCheck(err)

	//PIXELFORMAT is 4 for each primary color and alpha.
	pixels := make([]byte, WIDTH*HEIGHT*PIXELFORMAT)

	paintPixels(pixels)
	extras(tex, renderer, pixels)

	defer destroy(tex, renderer, window)

	sdl.Delay(2000)

}
