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

// NASAAPODJSONResponse struct to map URL values from the NASA APOD API JSON reponse.
type NASAAPODJSONResponse struct {
	URL string `json:"url"`
}

type ListPicturesURLFormRequest struct {
	From time.Time `form:"from" time_format:"2006-01-02" binding:"required"`
	To   time.Time `form:"to" time_format:"2006-01-02" binding:"required"`
}

// listPicturesURL returns JSON array with pictures URL.
func (server *Server) listPicturesURL(ctx *gin.Context) {
	var reqForm ListPicturesURLFormRequest
	if err := ctx.ShouldBindQuery(&reqForm); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if reqForm.From.After(reqForm.To) {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("'from' should be earlier than 'to'")))
		return
	}

	response, err := http.Get(
		fmt.Sprintf("%s?api_key=%s&start_date=%s&end_date=%s",
			NASA_APOD_API_URL,
			server.config.NASAAPIKey,
			reqForm.From.Format("2006-01-02"),
			reqForm.To.Format("2006-01-02")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if response.StatusCode == http.StatusTooManyRequests {
		ctx.JSON(http.StatusTooManyRequests, errorResponse(apiRateLimitExceededResponse()))
		return
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result []NASAAPODJSONResponse
	JSONError := json.Unmarshal(responseData, &result)
	if JSONError != nil {
		log.Fatal(JSONError)
	}

	ctx.JSON(http.StatusOK, gin.H{"urls": result})
}
