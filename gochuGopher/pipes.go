package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type pipes struct {
	mu      sync.RWMutex
	texture *sdl.Texture
	speed   int32
	pipes   []*pipe
}

func newPipes(r *sdl.Renderer) (*pipes, error) {
	tx, err := img.LoadTexture(r, "res/img/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("could not load pipe image: %v", err)
	}

	ps := &pipes{texture: tx, speed: 2}

	go func() {
		for {
			ps.mu.Lock()
			ps.pipes = append(ps.pipes, newPipe())
			ps.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return ps, nil
}

func (ps *pipes) paint(r *sdl.Renderer) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		if err := p.paint(r, ps.texture); err != nil {
			return err
		}
	}
	return nil
}

func (ps *pipes) update() {

	ps.mu.Lock()
	defer ps.mu.Unlock()

	var rem []*pipe
	for _, p := range ps.pipes {
		p.mu.Lock()
		p.x -= ps.speed
		if p.x+p.w > 0 {
			rem = append(rem, p)
		}
		p.mu.Unlock()
	}
	ps.pipes = rem
}
func (ps *pipes) restart() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.pipes = nil
}
func (ps *pipes) Destroy() {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.texture.Destroy()
}

type pipe struct {
	mu       sync.RWMutex
	x        int32
	h        int32
	w        int32
	inverted bool
}

func newPipe() *pipe {
	return &pipe{
		x:        800,
		h:        100 + int32(rand.Intn(270)),
		w:        50,
		inverted: rand.Float32() > 0.5,
	}
}

func (p *pipe) paint(r *sdl.Renderer, tx *sdl.Texture) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	rect := &sdl.Rect{X: p.x, Y: 600 - p.h, W: p.w, H: p.h}
	flip := sdl.FLIP_NONE
	if p.inverted {
		rect.Y = 0
		flip = sdl.FLIP_VERTICAL
	}

	if err := r.CopyEx(tx, nil, rect, 0, nil, flip); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}
	return nil
}
