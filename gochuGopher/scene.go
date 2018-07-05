package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	su      sync.Mutex
	bg      *sdl.Texture
	bird    *bird
	pipes   *pipes
	sounds  *sounds
	isStart bool
}

// get newscene
func newScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "res/img/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}

	bird, err := newBird(r)
	if err != nil {
		return nil, err
	}

	pipes, err := newPipes(r)
	if err != nil {
		return nil, err
	}

	return &scene{bg: bg, bird: bird, pipes: pipes}, nil
}

func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer, w *sdl.Window) chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			case <-tick:
				s.update()
				if !s.isStart {
					gravity = 0.0
				} else {
					gravity = 0.2
				}
				if s.bird.isDead() {
					s.sounds.die.Play()
					time.Sleep(time.Second)
					drawTitle(r, "Game Over")
					time.Sleep(time.Second)
					drawTitle(r, fmt.Sprintf("SCORE : %d", s.bird.score/25))
					time.Sleep(time.Second)
					s.restart()
				}
				if err := s.paint(r); err != nil {
					errc <- err
				}
				w.SetTitle(fmt.Sprintf("Gochu Score: %d", s.bird.score/25))
			}
		}
	}()

	return errc
}

func (s *scene) restart() {
	s.bird.score = 0
	s.isStart = false
	s.bird.restart()
	s.pipes.restart()
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch e := event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.KeyboardEvent:
		if e.GetType() == sdl.KEYDOWN && e.Keysym.Scancode == sdl.SCANCODE_SPACE && !s.bird.again {
			s.isStart = true
			go s.sounds.swoosh.Play()
			s.bird.jump()
		}
	}
	return false
}

func (s *scene) update() {
	s.bird.score++
	s.bird.update()
	if !s.bird.again {
		s.pipes.update()
		s.bird.touch(s)
	}
}

// paint newScene to r (Renderer)
func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()
	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}
	if err := s.pipes.paint(r); err != nil {
		return err
	}
	if err := s.bird.paint(r); err != nil {
		return err
	}
	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.Destroy()
	s.pipes.Destroy()
}
