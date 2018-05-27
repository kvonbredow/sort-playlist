package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

type handler struct {
	client spotify.Client
	user   string
}

func (h *handler) getTracks(playlist spotify.ID) ([]spotify.ID, error) {
	playlistTracks, err := h.client.GetPlaylistTracks(h.user, playlist)
	if err != nil {
		msg := fmt.Sprintf("couldn't get playlist tracks: %v", err)
		return nil, errors.New(msg)
	}

	ids := make([]spotify.ID, len(playlistTracks.Tracks))
	for k, pt := range playlistTracks.Tracks {
		ids[k] = pt.Track.ID
	}
	return ids, nil
}

func main() {
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}

	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	h := &handler{
		client: spotify.Authenticator{}.NewClient(token),
		user:   "kvb24",
	}

	ids, idErr := h.getTracks("71f1MdrRKbbPXhYgQGg0j0")
	if idErr != nil {
		log.Fatalf("%v", idErr)
	}
	for _, x := range ids {
		fmt.Println(x)
	}
}
