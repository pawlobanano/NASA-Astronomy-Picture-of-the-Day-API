package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// NASAAPODJSONResponse struct to map 'url' values from the NASA APOD API JSON reponse.
type NASAAPODJSONResponse struct {
	MediaType string `json:"media_type"`
	URL       string `json:"url"`
}

// ListPicturesURLFormRequest struct to map 'from' and 'to' values from the client's JSON request.
type ListPicturesURLFormRequest struct {
	From time.Time `form:"from" time_format:"2006-01-02" binding:"required"`
	To   time.Time `form:"to" time_format:"2006-01-02" binding:"required"`
}

// listPicturesURL returns JSON object with array of pictures URLs.
func (server *Server) listPicturesURL(ctx *gin.Context) {
	if UserIPLimiter[ctx.ClientIP()] {
		ctx.AbortWithStatusJSON(http.StatusTooManyRequests, errorResponse(apiRateLimitExceededResponse()))
		return
	}
	UserIPLimiter[ctx.ClientIP()] = true

	var reqForm ListPicturesURLFormRequest
	if err := ctx.ShouldBindQuery(&reqForm); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if reqForm.From.After(reqForm.To) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(errors.New("'from' should be earlier than 'to'")))
		return
	}

	startDate, err := time.Parse("2006-01-02", fmt.Sprint(reqForm.From.Format("2006-01-02")))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	endDate, err := time.Parse("2006-01-02", fmt.Sprint(reqForm.To.Format("2006-01-02")))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	days := make([]string, 0)
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		days = append(days, fmt.Sprint(d.Format("2006-01-02")))
	}

	var numJobs = len(days)
	jobs := make(chan string, numJobs)
	errors := make(chan string, numJobs)
	results := make(chan string, numJobs)
	for w := 1; w <= server.config.ConcurrentRequests; w++ {
		go server.worker(w, jobs, errors, results)
	}

	for _, day := range days {
		jobs <- day
	}
	close(jobs)

	URLs := make([]string, 0)
	for r := 1; r <= numJobs; r++ {
		result := <-results
		if result == "" {
			continue
		}
		URLs = append(URLs, result)
	}

	if <-errors == fmt.Sprint(http.StatusTooManyRequests) {
		ctx.AbortWithStatusJSON(http.StatusTooManyRequests, errorResponse(apiRateLimitExceededResponse()))
		UserIPLimiter[ctx.ClientIP()] = false
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"urls": URLs})
	UserIPLimiter[ctx.ClientIP()] = false
}

// worker is a goroutine taking 'w' (worker id), 'jobs', 'errors', 'results' channels.
func (server *Server) worker(w int, jobs <-chan string, errors chan<- string, results chan<- string) {
	for j := range jobs {
		log.Printf("worker %v started job - day %v", w, j)

		result, err := getAPIResponse(w, server.config.NASAAPIKey, j)
		if err != nil {
			errors <- err.Error()
			results <- ""
		} else {
			results <- result.URL
		}

		log.Printf("worker %v finished job - day %v", w, j)
	}
}

// getAPIResponse returns NASAAPODJSONResponse struct or error.
func getAPIResponse(id int, NASAAPIKey string, day string) (NASAAPODJSONResponse, error) {
	log.Printf("worker %v performed HTTP GET request - day %v", id, day)

	resp, err := http.Get(
		fmt.Sprintf("%s?api_key=%s&date=%s",
			NASA_APOD_API_URL,
			NASAAPIKey,
			day))
	if err != nil {
		return NASAAPODJSONResponse{}, err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return NASAAPODJSONResponse{}, fmt.Errorf("%v", http.StatusTooManyRequests)
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result NASAAPODJSONResponse
	JSONError := json.Unmarshal(respData, &result)
	if JSONError != nil {
		log.Fatal(JSONError)
	}

	if result.MediaType != "image" {
		return NASAAPODJSONResponse{}, fmt.Errorf("no image for the %s", day)
	}

	return result, nil
}
