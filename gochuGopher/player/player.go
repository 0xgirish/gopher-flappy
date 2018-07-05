package player

import (
	"fmt"
	"os/exec"
)

// Player to play music by exec.Command
type Player struct {
	command string
	file    string
}

// set os command to play sound file
func (p *Player) SetCommand(c string) {
	p.command = c
}

// New player get new player
func NewPlayer(f, c string) *Player {
	return &Player{
		command: c,
		file:    f,
	}
}

// set sound file
func (p *Player) SetSound(s string) {
	p.file = s
}

// play sound file
func (p *Player) Play() error {
	cmd := exec.Command(p.command, p.file)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not play: %v", err)
	}
	return nil
}
