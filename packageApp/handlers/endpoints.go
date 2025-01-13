package handlers

import (
	"log"
	"net/http"
	// "time"
	// "packageApp/db"
	"packageApp/models"

	"github.com/gin-gonic/gin"
)

func ServeEndpointsPage(c *gin.Context) {
	endpoints, err := models.FetchEndpoints()
	if err != nil {
		log.Println("Error fetching endpoints:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching endpoints"})
		return
	}
	c.HTML(http.StatusOK, "endpoints.html", gin.H{"endpoints": endpoints})
}

func HandleAddEndpoint(c *gin.Context) {
	var endpoint models.Endpoint
	if err := c.ShouldBindJSON(&endpoint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
		return
	}

	if err := models.AddEndpoint(endpoint); err != nil {
		log.Println("Error adding endpoint:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add endpoint"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Endpoint added successfully"})
}
