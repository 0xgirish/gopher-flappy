package main

import "gopher/player"

type sounds struct {
	wing   *player.Player
	swoosh *player.Player
	hit    *player.Player
	die    *player.Player
}

func newSounds(w, s, h, d, c string) *sounds {
	wing := player.NewPlayer(w, c)
	swoosh := player.NewPlayer(s, c)
	hit := player.NewPlayer(h, c)
	die := player.NewPlayer(d, c)
	return &sounds{
		wing:   wing,
		swoosh: swoosh,
		hit:    hit,
		die:    die,
	}
}
