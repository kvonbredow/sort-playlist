package main

import (
	"context"
	"flag"
	"log"
	"os"

	dp "github.com/kvonbredow/sort-playlist/download-pl"
	si "github.com/kvonbredow/sort-playlist/get-songinfo"
	"google.golang.org/grpc"
)

func main() {
	pl := flag.String("id", "1ervNoTdYL9SDWXJMcQNzJ", "Playlist ID")
	flag.Parse()

	// Build Clients
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	var dp_cc, si_cc *grpc.ClientConn
	var err error

	dp_cc, err = grpc.Dial(host+":5000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	dp_client := dp.NewDownloadPlaylistClient(dp_cc)

	si_cc, err = grpc.Dial(host+":5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	si_client := si.NewSongInfoClient(si_cc)

	// Pull Info
	track_req := &dp.PlaylistRequest{
		Id: *pl,
	}
	track_resp, track_err := dp_client.GetTracks(context.Background(), track_req)
	if track_err != nil {
		log.Fatalf("%v", track_err)
	}

	info_req := &si.InfoRequest{
		Ids: track_resp.GetIds(),
	}
	info_resp, info_err := si_client.GetInfo(context.Background(), info_req)
	if info_err != nil {
		log.Fatalf("%v", info_err)
	}
	for _, f := range info_resp.GetAfs() {
		log.Println(f)
	}
}
