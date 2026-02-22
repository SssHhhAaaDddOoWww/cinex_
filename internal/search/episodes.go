package search

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/SssHhhAaaDddOoWww/cinex_/internal/types"
	"github.com/joho/godotenv"
)

func GetEpisodes(seriesID int, seasonNumber int) ([]types.Episode, error) {
	godotenv.Load()
	apiKey := os.Getenv("TMDB_Key")

	if apiKey == "" {
		return nil, fmt.Errorf("key not set")
	}

	endPoint := fmt.Sprintf(
		"https://api.themoviedb.org/3/tv/%d/season/%d",
		seriesID, seasonNumber,
	)

	req, err := http.NewRequest("GET", endPoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Episodes []types.Episode `json:"episodes"`
	}
	json.Unmarshal(body, &result)

	return result.Episodes, nil
}
