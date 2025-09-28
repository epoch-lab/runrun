package handler

import (
	"log"
	"math/rand"
	"net/http"
	"runrun/internal/protocol"
	"time"

	"github.com/gin-gonic/gin"
)

// RunHandler handles the request to trigger a run.
// For now, it simulates a run for a single, hardcoded user
// without database interaction.
func RunHandler(c *gin.Context) {
	log.Println("Received request to trigger a run.")

	// --- Database part is mocked for now ---
	// In the future, we will fetch users from the database here.
	// For now, using a hardcoded example user.
	// IMPORTANT: These credentials will likely fail, this is just for demonstrating the flow.
	const fakePhone = "19871265362"
	const fakePassword = "ly373452721"

	log.Printf("Simulating run for user: %s", fakePhone)

	// 1. Generate fake client info
	clientInfo := protocol.GenerateFakeClient()

	// 2. Login to get user info and token
	// This will make a real network request to the tanmasports server.
	userInfo, err := protocol.Login(fakePhone, fakePassword, clientInfo)
	if err != nil {
		log.Printf("ERROR: Failed to login for fake user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to login for fake user. Check credentials or network.",
			"error":   err.Error(),
		})
		return
	}
	log.Printf("Successfully logged in as user ID: %d", userInfo.UserID)

	// 3. Generate random run data
	rand.Seed(time.Now().UnixNano())
	runDistance := int64(4500 + rand.Intn(500))    // 4.5km to 5.0km
	runDuration := int32(25 + rand.Intn(5))      // 25 to 29 minutes

	log.Printf("Generated run data: distance %d meters, duration %d minutes", runDistance, runDuration)

	// 4. Submit the run
	// This also makes a real network request.
	err = protocol.Submit(*userInfo, clientInfo, runDuration, runDistance)
	if err != nil {
		log.Printf("ERROR: Failed to submit run for user %d: %v", userInfo.UserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to submit run.",
			"error":   err.Error(),
		})
		return
	}

	log.Printf("Successfully submitted run for user ID: %d", userInfo.UserID)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Run process completed successfully for fake user.",
		"userId":  userInfo.UserID,
	})
}
