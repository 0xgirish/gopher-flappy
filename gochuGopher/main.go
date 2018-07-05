package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	if err := run(); err != nil {
		log.Panic(err)
	}
}

func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("could not initialize SDL: %v", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("could not initialize ttf: %v", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer w.Destroy()

	w.SetTitle("Gochu Gopher")

	r.SetDrawColor(255, 255, 255, 255)

	if err := drawTitle(r, "Gochu Gopher"); err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}

	time.Sleep(time.Second)

	s, err := newScene(r)

	swoosh := "res/sounds/sfx_swooshing.ogg"
	hit := "res/sounds/sfx_hit.ogg"
	die := "res/sounds/sfx_die.ogg"
	s.sounds = newSounds(swoosh, hit, die, "ogg123")

	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer s.destroy()

	events := make(chan sdl.Event)

	go func() {
		for {
			events <- sdl.WaitEvent()
		}
	}()

	return <-s.run(events, r, w)
}

func drawTitle(r *sdl.Renderer, title string) error {
	r.Clear()
	tf, err := ttf.OpenFont("res/fonts/Flappy.ttf", 10)

	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	defer tf.Close()

	color := sdl.Color{R: 241, G: 146, B: 17, A: 255}
	surf, err := tf.RenderUTF8Solid(title, color)
	if err != nil {
		return fmt.Errorf("could not render font: %v", err)
	}
	defer surf.Free()

	tx, err := r.CreateTextureFromSurface(surf)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer tx.Destroy()

	if err := r.Copy(tx, nil, &sdl.Rect{X: 320, Y: 270, W: 160, H: 60}); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}
	r.Present()

	return nil
}
