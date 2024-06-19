package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/kkdai/youtube/v2"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Uso: %s <URL da playlist> <diretório de saída>", os.Args[0])
	}

	playlistURL := os.Args[1]
	dir := os.Args[2]

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			log.Fatalf("Erro ao criar o diretório %s: %v", dir, err)
		}
	}

	err := downloadPlaylist(playlistURL, dir)
	if err != nil {
		log.Fatalf("Erro ao baixar a playlist: %v", err)
	}

	fmt.Println("Conversão concluída!")
}

func downloadPlaylist(url, dir string) error {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 60 * time.Second,
	}

	client := youtube.Client{
		HTTPClient: &http.Client{
			Transport: transport,
			Timeout:   180 * time.Second,
		},
	}

	playlist, err := client.GetPlaylist(url)
	if err != nil {
		return fmt.Errorf("erro ao obter a playlist: %w", err)
	}

	ctx := context.Background()

	semaphore := make(chan struct{}, 5)
	var wg sync.WaitGroup

	for _, entry := range playlist.Videos {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(entry *youtube.PlaylistEntry) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			video, err := client.GetVideo(entry.ID)
			if err != nil {
				log.Printf("Erro ao obter o vídeo %s: %v", entry.Title, err)
				return
			}

			err = downloadAndConvert(ctx, &client, video, dir)
			if err != nil {
				log.Printf("Erro ao baixar e converter o vídeo %s: %v", entry.Title, err)
			}
		}(entry)
	}

	wg.Wait()
	return nil
}

func downloadAndConvert(ctx context.Context, client *youtube.Client, video *youtube.Video, dir string) error {
	formats := video.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		return fmt.Errorf("erro ao obter o stream: %w", err)
	}
	defer stream.Close()

	audioBytes, err := ioutil.ReadAll(stream)
	if err != nil {
		return fmt.Errorf("erro ao ler o stream de áudio: %w", err)
	}

	mp3FilePath := filepath.Join(dir, sanitizeFilename(video.Title)+".mp3")
	err = ioutil.WriteFile(mp3FilePath, audioBytes, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar o arquivo MP3: %w", err)
	}

	fmt.Printf("Baixado e convertido %s\n", video.Title)
	return nil
}

func sanitizeFilename(name string) string {
	name = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(name, "")
	name = strings.ReplaceAll(name, " ", "_")
	if len(name) > 100 {
		name = name[:100]
	}
	return name
}
