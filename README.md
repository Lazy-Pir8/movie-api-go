# Movie Info API (Go + Gin)

A simple RESTful API server in Go that fetches movie and TV show data from the OMDb API.  
Perfect for assignments, demo projects, or as a backend learning exercise!

## Features

- **/api/movie**: Get details about a movie by title.
- **/api/episode**: Get specific TV episode info (by series, season, episode).
- **/api/movies/genre**: List top 15 movies for a genre, sorted by IMDb rating.
- **/api/recommendations**: Get smart movie recommendations based on a favorite movie (by genre/director/actors).

## How To Run

1. **Clone the repository:**
2. **Install dependencies:**


3. **Create a `.env` file with your OMDb API key:**
4. **Run the server:**
The API will listen on `http://localhost:8080`.

## API Endpoints

### 1. Movie Details  
`GET /api/movie?title=Inception`
- **Returns**: Title, Year, Plot, Country, Awards, Director, Ratings.

### 2. Episode Details  
`GET /api/episode?series_title=Breaking Bad&season=1&episode_number=1`
- **Returns**: Title, Released Date, Season, Episode, IMDb Rating, Plot.

### 3. Genre Movies  
`GET /api/movies/genre?genre=Action`
- **Returns**: Top 15 movies in that genre frm a pool, sorted by IMDb rating.

### 4. Recommendations  
`GET /api/recommendations?favorite_movie=Inception`
- **Returns**: Up to 20 recommended movies, prioritized by shared genre, director, and actors.


---
