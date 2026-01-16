package metadata

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rangodisco/yhar/config"
	"github.com/rangodisco/yhar/internal/metadata/config/database"
	"github.com/rangodisco/yhar/thirdpartyAPIs/anna/config/router"
)

var Router *gin.Engine

func init() {
	// Load main environment variables
	godotenv.Load(".env")

	// Ensure that the environment is set to test
	os.Setenv("GIN_MODE", "test")

	// Load environment variables
	config.LoadEnv()

	database.SetupDatabase()

	// Setup router
	Router = router.SetupRouter()
}
