# ebiten-asteroids

A recreation of the classic Atari game Asteroids with Go and Ebiten.

To keep with the original Asteroids, the graphics are all generated from lines drawn on the screen instead of by loading assets.

## ToDo

- [ ] Title screen
- [ ] Asteroids and enemy ships
- [ ] Collision detection
- [ ] Shooting
- [ ] Score
- [ ] Player lives
- [ ] Sound
- [ ] Game over condition + screen

## Build & run

Built with Go 1.16 and [Ebiten v2](https://github.com/hajimehoshi/ebiten) on MacOS.

     go build
    ./ebiten-asteroids

The `game.json` file is for controlling certain settings, like window size and fullscreen toggle.

## Keymaps

Arrows keys or WASD : Move

'p' : debug info toggle
