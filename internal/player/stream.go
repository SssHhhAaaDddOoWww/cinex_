package player

import (
	"fmt"
	"net/http"
	"time"

	"github.com/anacrolix/torrent"
)

func Start(magnet string) error {

	cfg := torrent.NewDefaultClientConfig()
	cfg.Seed = false // we're streaming, not seeding
	client, err := torrent.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create torrent client: %w", err)
	}
	defer client.Close()

	fmt.Println("connecting to peers...")
	t, err := client.AddMagnet(magnet)
	if err != nil {
		return fmt.Errorf("failed to add magnet: %w", err)
	}
	<-t.GotInfo()

	file := largestFile(t)
	if file == nil {
		return fmt.Errorf("no files found in torrent")
	}

	fmt.Printf("found: %s (%.2f MB)\n", file.DisplayPath(), float64(file.Length())/(1<<20))

	t.DownloadAll()

	go serveFile(file)

	time.Sleep(1 * time.Second)

	streamURL := "http://localhost:8888/stream"
	fmt.Printf("opening player: %s\n", streamURL)
	return OpenPlayer(streamURL)
}

func largestFile(t *torrent.Torrent) *torrent.File {
	var largest *torrent.File
	for _, f := range t.Files() {
		f := f
		if largest == nil || f.Length() > largest.Length() {
			largest = f
		}
	}
	return largest
}

func serveFile(file *torrent.File) {
	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		reader := file.NewReader()
		// readahead tells anacrolix to prioritize downloading
		// the next chunk so playback stays smooth
		reader.SetReadahead(file.Length() / 100)
		http.ServeContent(w, r, file.DisplayPath(), time.Now(), reader)
	})

	fmt.Println("stream server running on :8888")
	http.ListenAndServe(":8888", nil)
}
