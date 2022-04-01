package test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	DB              *gorm.DB
	testingDataPath string
)

const (
	testingDataFile = "dbadmin/testing_data.sql"
)

type Endpoint struct {
	HTTPMethod string
	Path       string
	Handler    gin.HandlerFunc
}

func BuildTestRouter(endpoints []Endpoint, middlewares ...gin.HandlerFunc) *gin.Engine {
	engine := gin.New()

	engine.Use(middlewares...)
	for _, endpoint := range endpoints {
		engine.Handle(endpoint.HTTPMethod, endpoint.Path, endpoint.Handler)
	}
	return engine
}

// ReverseWalkFindFile searches for file walking the path
// in reservse way, ./../..
func ReverseWalkFindFile(filename string) (string, error) {
	found := false
	ogCwd, _ := os.Getwd()
	datapath := ""
	prevCwd := ""
	for !found {
		cwd, _ := os.Getwd()
		if cwd == prevCwd {
			err := os.Chdir(ogCwd)
			if err != nil {
				return "", err
			}
			return "", fmt.Errorf("Cannot find file: %s", err)
		}
		datapath = filepath.Join(cwd, filename)
		_, err := os.Stat(datapath)
		if err == nil {
			found = true
		} else {
			prevCwd = cwd
			err := os.Chdir("..")
			if err != nil {
				return "", err
			}
		}
	}
	err := os.Chdir(ogCwd)
	if err != nil {
		return "", err
	}
	return datapath, nil
}

func ResetDB() error {
	if testingDataPath == "" {
		var err error
		// Because golang runs tests always in the current tested module
		// folder, we must return to the root folder and find testing data file
		path, err := ReverseWalkFindFile(testingDataFile)
		if err != nil {
			return err
		}
		// cache the full path for current run
		testingDataPath = path
	}
	buf, err := os.ReadFile(testingDataPath)
	if err != nil {
		return err
	}
	plainDb, err := DB.DB()
	if err != nil {
		return err
	}
	_, err = plainDb.Exec("TRUNCATE TABLE account, cluster, image, repository, repository_image, cluster_image, cve, image_cve CASCADE")
	if err != nil {
		return err
	}
	_, err = plainDb.Exec(string(buf))
	return err
}
