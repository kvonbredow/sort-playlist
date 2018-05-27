package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

type handler struct {
	client spotify.Client
}

func (h *handler) getInfo(ids []spotify.ID) ([]*spotify.AudioFeatures, error) {
	return h.client.GetAudioFeatures(ids...)
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
	}

	fp, fpErr := os.Open("../ids.txt")
	if fpErr != nil {
		log.Fatalf("couldn't open file: %v", err)
	}
	defer fp.Close()

	var ids []spotify.ID
	sc := bufio.NewScanner(fp)
	for sc.Scan() {
		ids = append(ids, spotify.ID(sc.Text()))
	}
	if sc.Err() != nil {
		log.Fatalf("couldn't read file: %v", sc.Err())
	}

	afs, infoErr := h.getInfo(ids)
	if infoErr != nil {
		log.Fatalf("couldn't get audio features: %v", infoErr)
	}

	for _, x := range afs {
		fmt.Println(x.Key)
	}
}
