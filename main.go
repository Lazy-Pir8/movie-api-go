package main

// all the imports
import (
	"encoding/json"
	//"log"
	// "fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func containsGenre(moviesGenres string, targetGenre string) bool {
	genres := strings.Split(moviesGenres, ",")
	for _, g := range genres {
		if strings.EqualFold(strings.TrimSpace(g), targetGenre) {
			return true
		}
	}
	return false
}

// data structures
type Movie struct {
	Title    string        `json:"Title"`
	Year     string        `json:"Year"`
	Plot     string        `json:"Plot"`
	Country  string        `json:"Country"`
	Awards   string        `json:"Awards"`
	Director string        `json:"Director"`
	Ratings  []interface{} `json:"Ratings"`
}

type Episode struct {
	Title      string `json:"Title"`
	Released   string `json:"Released"`
	Season     string `json:"Season"`
	Episode    string `json:"Episode"`
	ImdbRating string `json:"imdbRating"`
	Plot       string `json:"Plot"`
}

type OMDbError struct {
	Response string `json:"Response"`
	Error    string `json:"Error"`
}

var titlesToTry = []string{
	"Inception", "The Dark Knight", "Forrest Gump", "Pulp Fiction", "The Matrix",
	"Gladiator", "The Godfather", "Fight Club", "Interstellar", "Avengers: Endgame",
	"The Shawshank Redemption", "Titanic", "Avatar", "Joker", "The Lion King",
	"The Avengers", "Jurassic Park", "Back to the Future", "The Departed", "Skyfall",
	"Mad Max: Fury Road", "Braveheart", "Whiplash", "Goodfellas", "The Prestige",
	// ...add more
}

type ShortMovieInfo struct {
	Title      string `json:"Title"`
	Year       string `json:"Year"`
	ImdbRating string `json:"imdbRating"`

	Genre string `json:"Genre"`
}

type RecMovieInfo struct {
	Title      string `json:"Title"`
	Year       string `json:"Year"`
	ImdbRating string `json:imdbRating`
	Genre      string `json:"Genre"`
	Director   string `json:"Director"`
	Actors     string `json:"Actors"`
}

func containsAny(a, b string) bool {
	aList := strings.Split(a, ",")
	bList := strings.Split(b, ",")
	for _, av := range aList {
		for _, bv := range bList {
			if strings.EqualFold(strings.TrimSpace(av), strings.TrimSpace(bv)) {
				return true
			}
		}
	}
	return false
}

func main() {
	_ = godotenv.Load()

	router := gin.Default()
	// ignore this
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// this is where my code starts
	router.GET("/api/movie", func(c *gin.Context) {
		title := c.Query("title")
		if title == "" {
			c.JSON(400, gin.H{"error": "Missing 'title' query parameter"})
			return
		}

		apiKey := os.Getenv("OMDB_API_KEY")
		url := "http://www.omdbapi.com/?apikey=" + apiKey + "&t=" + title

		resp, err := http.Get(url)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to reach OMDb API"})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to read OMDB response"})
			return
		}

		var omdbErr OMDbError
		json.Unmarshal(body, &omdbErr)
		if omdbErr.Response == "False" {
			c.JSON(404, gin.H{"error": omdbErr.Error})
			return
		}
		var movie Movie
		err = json.Unmarshal(body, &movie)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to decode OMDb response"})
			return
		}

		c.JSON(200, movie)

	})

	router.GET("/api/episode", func(c *gin.Context) {
		seriesTitle := c.Query("series_title")
		season := c.Query("season")
		episodeNum := c.Query("episode_number")

		if seriesTitle == "" || season == "" || episodeNum == "" {
			c.JSON(400, gin.H{"error": "Missing series_title, season, or episode_number"})
			return
		}

		apiKey := os.Getenv("OMDB_API_KEY")
		encodedTitle := strings.ReplaceAll(seriesTitle, " ", "+")
		url := "http://www.omdbapi.com/?apikey=" + apiKey + "&t=" +
			encodedTitle + "&Season=" + season + "&Episode=" + episodeNum

		// log.Println("OMDb URL:", url)

		resp, err := http.Get(url)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to reach OMDb API"})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to reach OMDb API"})
			return
		}
		// log.Println("OMDb RAW:", string(body))

		var omdbErr OMDbError
		json.Unmarshal(body, &omdbErr)
		if omdbErr.Response == "False" {
			c.JSON(404, gin.H{"error": omdbErr.Error})
			return
		}

		var episode Episode
		err = json.Unmarshal(body, &episode)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to decode OMDb response"})
			return
		}
		// fmt.Println(string(body))
		c.JSON(200, episode)

	})
	// check for movies or movie?
	router.GET("/api/movies/genre", func(c *gin.Context) {
		genre := c.Query("genre")
		if genre == "" {
			c.JSON(400, gin.H{"error": "Missing genre query parameter"})
			return
		}
		apiKey := os.Getenv("OMDB_API_KEY")

		var results []ShortMovieInfo

		for _, title := range titlesToTry {
			url := "http://www.omdbapi.com/?apikey=" + apiKey + "&t=" + title

			resp, err := http.Get(url)
			if err != nil {
				continue
			}
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			var omdbErr OMDbError
			json.Unmarshal(body, &omdbErr)
			if omdbErr.Response == "False" {
				continue
			}

			var movie ShortMovieInfo
			_ = json.Unmarshal(body, &movie)

			if containsGenre(movie.Genre, genre) {
				results = append(results, movie)
			}
		}
		sort.Slice(results, func(i, j int) bool {
			ratingI, _ := strconv.ParseFloat(results[i].ImdbRating, 64)
			ratingJ, _ := strconv.ParseFloat(results[j].ImdbRating, 64)
			return ratingI > ratingJ
		})

		if len(results) > 15 {
			results = results[:15]
		}

		c.JSON(200, results)
	})
	// fix the episode handler (done)
	router.GET("/api/recommendations", func(c *gin.Context) {
		favorite := c.Query("favorite_movie")
		if favorite == "" {
			c.JSON(400, gin.H{"error": "Missing favorite_movie query parameter"})
			return
		}
		apiKey := os.Getenv("OMDB_API_KEY")

		url := "http://www.omdbapi.com/?apikey=" + apiKey + "&t=" + favorite
		resp, err := http.Get(url)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to reach OMDb API"})
			return
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var omdbErr OMDbError
		json.Unmarshal(body, &omdbErr)
		if omdbErr.Response == "False" {
			c.JSON(404, gin.H{"error": omdbErr.Error})
			return
		}
		var base RecMovieInfo
		json.Unmarshal(body, &base)

		var genreMatches, directorMatches, actorMatches []RecMovieInfo

		for _, title := range titlesToTry {
			if strings.EqualFold(title, base.Title) {
				continue
			}
			url := "http://www.omdbapi.com/?apikey=" + apiKey + "&t=" + title
			resp, err := http.Get(url)
			if err != nil {
				continue
			}
			b, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			var omdbErr OMDbError
			json.Unmarshal(b, &omdbErr)
			if omdbErr.Response == "False" {
				continue
			}
			var m RecMovieInfo
			json.Unmarshal(b, &m)

			if containsAny(base.Genre, m.Genre) {
				genreMatches = append(genreMatches, m)

			}
			if containsAny(base.Director, m.Director) {
				directorMatches = append(directorMatches, m)
				continue
			}
			if containsAny(base.Actors, m.Actors) {
				actorMatches = append(actorMatches, m)
			}
		}

		final := append([]RecMovieInfo{}, genreMatches...)
		if len(final) < 20 {
			n := 20 - len(final)
			if n > len(directorMatches) {
				n = len(directorMatches)
			}
			final = append(final, directorMatches[:n]...)
		}
		if len(final) < 20 {
			n := 20 - len(final)
			if n > len(actorMatches) {
				n = len(actorMatches)

			}
			final = append(final, actorMatches[:n]...)
		}

		sort.Slice(final, func(i, j int) bool {
			ri, _ := strconv.ParseFloat(final[i].ImdbRating, 64)
			rj, _ := strconv.ParseFloat(final[j].ImdbRating, 64)
			return ri > rj
		})
		if len(final) > 20 {
			final = final[:20]
		}
		c.JSON(200, final)

	})

	// fmt.Println("API KEY:", apiKey)

	router.Run(":8080")

}
