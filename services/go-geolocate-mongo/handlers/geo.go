package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-geolocate-mongo/models"
	"go-geolocate-mongo/pkg/mongo"
)

var googleAPIKey = os.Getenv("GOOGLE_API_KEY")

type geolocateReq struct {
	ConsiderIP bool `json:"considerIp"`
}

type geolocateResp struct {
	Location struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"location"`
	Accuracy float64 `json:"accuracy"`
}

type geocodeResp struct {
	Results []struct {
		FormattedAddress string `json:"formatted_address"`
	} `json:"results"`
	Status string `json:"status"`
}

// GET /detect/ip
func DetectIP(c *gin.Context) {
	// Call Google Geolocation API with considerIp=true
	url := fmt.Sprintf("https://www.googleapis.com/geolocation/v1/geolocate?key=%s", googleAPIKey)
	payload := geolocateReq{ConsiderIP: true}
	b, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		c.JSON(resp.StatusCode, gin.H{"error": string(body)})
		return
	}
	var gr geolocateResp
	if err := json.Unmarshal(body, &gr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Save to DB
	loc := models.Location{
		Lat:       gr.Location.Lat,
		Lng:       gr.Location.Lng,
		Accuracy:  gr.Accuracy,
		Source:    "ip",
		CreatedAt: time.Now(),
	}
	col := mongo.Collection(os.Getenv("MONGO_DB"), os.Getenv("MONGO_COLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = col.InsertOne(ctx, loc)

	c.JSON(http.StatusOK, gin.H{"lat": gr.Location.Lat, "lng": gr.Location.Lng, "accuracy": gr.Accuracy})
}

// POST /reverse { lat, lng }
func ReverseGeocode(c *gin.Context) {
	var body struct {
		Lat    float64 `json:"lat" binding:"required"`
		Lng    float64 `json:"lng" binding:"required"`
		UserID string  `json:"userId,omitempty"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat & lng required"})
		return
	}
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?latlng=%f,%f&key=%s", body.Lat, body.Lng, googleAPIKey)
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		c.JSON(resp.StatusCode, gin.H{"error": string(b)})
		return
	}
	var gr geocodeResp
	if err := json.Unmarshal(b, &gr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	address := ""
	if len(gr.Results) > 0 {
		address = gr.Results[0].FormattedAddress
	}
	// store
	loc := models.Location{
		UserID:    body.UserID,
		Lat:       body.Lat,
		Lng:       body.Lng,
		Address:   address,
		Source:    "client",
		CreatedAt: time.Now(),
	}
	col := mongo.Collection(os.Getenv("MONGO_DB"), os.Getenv("MONGO_COLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = col.InsertOne(ctx, loc)

	c.JSON(http.StatusOK, gin.H{"address": address, "status": gr.Status})
}

// POST /save { lat, lng, address, userId }
func SaveLocation(c *gin.Context) {
	var payload models.Location
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.CreatedAt.IsZero() {
		payload.CreatedAt = time.Now()
	}
	if payload.Source == "" {
		payload.Source = "manual"
	}
	col := mongo.Collection(os.Getenv("MONGO_DB"), os.Getenv("MONGO_COLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := col.InsertOne(ctx, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"saved": true})
}

// GET /locations?userId=...
func ListLocations(c *gin.Context) {
	userId := c.Query("userId")
	col := mongo.Collection(os.Getenv("MONGO_DB"), os.Getenv("MONGO_COLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{}
	if userId != "" {
		filter["userId"] = userId
	}
	opts := options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(50)
	cur, err := col.Find(ctx, filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var results []models.Location
	if err := cur.All(ctx, &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"locations": results})
}
