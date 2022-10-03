package api

import (
	"encoding/json"
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
	var reqForm ListPicturesURLFormRequest
	if err := ctx.ShouldBindQuery(&reqForm); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if reqForm.From.After(reqForm.To) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(fromLaterThanToResponse()))
		return
	}

	fromDate, err := parseDate(reqForm.From, ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	toDate, err := parseDate(reqForm.To, ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	days := make([]string, 0)
	for d := fromDate; !d.After(toDate); d = d.AddDate(0, 0, 1) {
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

	select {
	case s, ok := <-errors:
		if ok && s == fmt.Sprint(http.StatusTooManyRequests) {
			log.Printf("apiRateLimitExceeded with %v error code", s)
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, errorResponse(apiRateLimitExceededResponse()))
			return
		} else {
			log.Printf("errors channel closed.")
		}
	default:
		log.Printf("no value ready on errors channel, moving on.")
	}

	ctx.JSON(http.StatusOK, gin.H{"urls": URLs})
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

	if resp.Body != nil {
		defer resp.Body.Close()
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

// parseDate parses reqTime to "2006-01-02" time format.
func parseDate(reqTime time.Time, ctx *gin.Context) (time.Time, error) {
	date, err := time.Parse("2006-01-02", fmt.Sprint(reqTime.Format("2006-01-02")))
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}
