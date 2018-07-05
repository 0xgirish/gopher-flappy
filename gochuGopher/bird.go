package main

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

var gravity = 0.2

const (
	jumpSpeed   = 5
	decayFactor = 0.9
)

type bird struct {
	mu       sync.RWMutex
	time     int
	textures []*sdl.Texture
	x, y     int32
	w, h     int32
	speed    float64
	dead     bool
	score    int
}

func newBird(r *sdl.Renderer) (*bird, error) {
	var textures []*sdl.Texture
	for i := 1; i <= 4; i++ {
		texture, err := img.LoadTexture(r, fmt.Sprintf("res/img/frame-%d.png", i))
		if err != nil {
			return nil, fmt.Errorf("could not load frame-image: %v", err)
		}
		textures = append(textures, texture)
	}

	return &bird{textures: textures, x: 10, y: 300, speed: 0, w: 50, h: 43}, nil
}

func (b *bird) update() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.time++
	b.y -= int32(b.speed)

	if b.y < 0 {
		b.speed = -b.speed * decayFactor
	} else {
		b.speed += gravity
	}
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	birdPosition := &sdl.Rect{X: b.x, Y: (600 - b.y) - b.h/2, W: b.w, H: b.h}
	frame := b.time / 10 % len(b.textures)
	if err := r.Copy(b.textures[frame], nil, birdPosition); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}
	return nil
}

func (b *bird) restart() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.y = 300
	b.speed = 0
	b.dead = false
}

func (b *bird) Destroy() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, tx := range b.textures {
		tx.Destroy()
	}
}

func (b *bird) jump() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.speed = -jumpSpeed
}

func (b *bird) isDead() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.dead
}

func (b *bird) touch(s *scene) {
	b.mu.Lock()
	defer b.mu.Unlock()
	isTouch := false
	ps := s.pipes
	for _, p := range ps.pipes {
		p.mu.RLock()

		if p.x > b.x+b.w-4 {
			p.mu.RUnlock()
			continue
		}

		if p.x+p.w-4 < b.x {
			p.mu.RUnlock()
			continue
		}

		if !p.inverted && p.h < b.y-b.h/2+10 {
			p.mu.RUnlock()
			continue
		}

		if p.inverted && (600-p.h+10) > b.y+b.h/2 {
			p.mu.RUnlock()
			continue
		}
		isTouch = true
		s.sounds.hit.Play()
		p.mu.RUnlock()
		break
	}

	if isTouch {
		b.dead = true
	}
}

func (b *bird) sound() {

}
