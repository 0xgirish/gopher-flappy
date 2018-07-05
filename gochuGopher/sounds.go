package main

import "gopher/gochuGopher/player"

type sounds struct {
	swoosh *player.Player
	hit    *player.Player
	die    *player.Player
}

func newSounds(s, h, d, c string) *sounds {
	swoosh := player.NewPlayer(s, c)
	hit := player.NewPlayer(h, c)
	die := player.NewPlayer(d, c)
	return &sounds{
		swoosh: swoosh,
		hit:    hit,
		die:    die,
	}
}
