package docs

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/redbonzai/user-management-api/pkg/logger"
	"github.com/swaggo/swag"
	"go.uber.org/zap"
)

type swaggerInfo struct {
	Title       string
	Description string
	Version     string
	Host        string
	BasePath    string
}

func (swaggerInfo *swaggerInfo) ReadDoc() string {
	swaggerDir := "docs"
	swaggerFile := "swagger.yaml"

	filePath := filepath.Join(swaggerDir, swaggerFile)
	file, err := os.Open(filePath)
	if err != nil {
		logger.Fatal("failed to open Swagger file:", zap.Error(err))
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Fatal("failed to read Swagger file:", zap.Error(err))
	}

	return string(content)
}

var swaggerInfoInstance = swaggerInfo{
	Title:       "User Management API",
	Description: "User CRUD Operations The Badass Way.",
	Version:     "1.0.0",
	Host:        "localhost:8080",
	BasePath:    "/",
}

func init() {
	swag.Register(swag.Name, &swaggerInfoInstance)
}
