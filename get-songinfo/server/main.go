package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/kvonbredow/sort-playlist/get-songinfo"
)

type handler struct {
	client spotify.Client
}

func (h *handler) GetInfo(ctx context.Context, req *pb.InfoRequest) (*pb.InfoResponse, error) {
	ids := make([]spotify.ID, len(req.Ids))
	for k, id := range req.Ids {
		ids[k] = spotify.ID(id)
	}
	afs, err := h.client.GetAudioFeatures(ids...)
	if err != nil {
		return nil, err
	}

	features := make([]*pb.Features, len(afs))
	for k, af := range afs {
		features[k] = &pb.Features{
			Acousticness:     af.Acousticness,
			AnalysisURL:      af.AnalysisURL,
			Danceability:     af.Danceability,
			Duration:         int64(af.Duration),
			Energy:           af.Energy,
			Instrumentalness: af.Instrumentalness,
			Key:              int64(af.Key),
			Liveness:         af.Liveness,
			Loudness:         af.Loudness,
			Mode:             int64(af.Mode),
			Speechiness:      af.Speechiness,
			Tempo:            af.Tempo,
			TimeSignature:    int64(af.TimeSignature),
			TrackURL:         af.TrackURL,
		}
	}

	resp := &pb.InfoResponse{
		Afs: features,
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
	}

	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSongInfoServer(s, h)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
