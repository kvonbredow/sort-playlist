package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/kvonbredow/sort-playlist/download-pl"
)

type handler struct {
	client spotify.Client
	user   string
}

func (h *handler) GetTracks(ctx context.Context, req *pb.PlaylistRequest) (*pb.PlaylistResponse, error) {
	id := spotify.ID(req.Id)
	playlistTracks, err := h.client.GetPlaylistTracks(h.user, id)
	if err != nil {
		msg := fmt.Sprintf("couldn't get playlist tracks: %v", err)
		return nil, errors.New(msg)
	}

	ids := make([]string, len(playlistTracks.Tracks))
	for k, pt := range playlistTracks.Tracks {
		ids[k] = string(pt.Track.ID)
	}
	resp := &pb.PlaylistResponse{
		Ids: ids,
	}
	return resp, nil
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

	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDownloadPlaylistServer(s, h)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
