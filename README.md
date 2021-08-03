# GamesGo

GoGame was started with the idea to get acquainted with Golang and make something fun in the process.
Games currently run on SDL2 which is a library written in C and has bindings available for Go but also other
languages as well.

# Directories

It contains 3 different folders.

+ IntroSDL has some SDL2 implementation just to test the library and build something simple. An image will appear
and after some seconds disappear.

+ Pong1.0 has the first iteration of Pong. The AI is pretty simple and doesn't make mistakes so it's impossible to win.
It also doesn't have a scoreboard.

+ Pong2.0 has the final version of Pong. It is as close to the original Pong. In this iteration it's possible to win as the AI is improved
and can and will make mistakes from time to time. It also has a scoreboard and a few bugs have been solved from previous iteration.

# Requirements

Below is the command to install the required packages in Ubuntu. To install in other Linux distributions, Windows or macOS
go to https://github.com/veandco/go-sdl2 where everything is explained.
Pong was built with the latest Golang release that at this time is 1.16.5 but should work with 1.13 forward.

On Ubuntu 14.04 and above, type:

    apt install libsdl2{,-image,-mixer,-ttf,-gfx}-dev

Might have to use sudo for the command to work correctly.

# Run the Program

+ Clone the repo:

      git clone https://github.com/santimuriado/GamesGo.git
    
+ Generate the go.sum with:

      go mod tidy
      
+ Run the game:

      go run pong.go
    
