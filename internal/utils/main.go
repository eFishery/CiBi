package utils

import (
	"os"
	"fmt"
	"log"
	"strconv"
	"strings"

	"net/mail"
	"io/ioutil"

	yaml "github.com/goccy/go-yaml"
)

func (h *CiBiConfig) Load(filename string) (*CiBiConfig, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, h)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(h.NotionConfig.ApiKey, "$"){
		originalVarName := strings.TrimPrefix(h.NotionConfig.ApiKey, "$")
		h.NotionConfig.ApiKey = getEnvString(originalVarName, "EMPTY NOTION API KEY")
	}

	if strings.HasPrefix(h.NotionConfig.DatabaseID, "$"){
		originalVarName := strings.TrimPrefix(h.NotionConfig.DatabaseID, "$")
		h.NotionConfig.DatabaseID = getEnvString(originalVarName, "EMPTY NOTION DATABASE ID")
	}

	if strings.HasPrefix(h.NotionConfig.DeploymentLog, "$"){
		originalVarName := strings.TrimPrefix(h.NotionConfig.DeploymentLog, "$")
		h.NotionConfig.DeploymentLog = getEnvString(originalVarName, "EMPTY NOTION DATABASE ID")
	}

	h.FileName = filename

	return h, nil
}

func getEnvString(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		value, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		return value
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "true"
	}

	return defaultVal
}

func Tprintf(format string, params map[string]interface{}) string {
	for key, val := range params {
		format = strings.Replace(format, "%{"+key+"}s", fmt.Sprintf("%s", val), -1)
	}
	return format
}

func ValidMailAddress(email string) (string, bool) {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", false
	}
	return addr.Address, true
}

func GetCiBiPrefix() []string {
	files, err := ioutil.ReadDir("./")
	if err != nil{
		log.Fatal(err)
	}

	var f []string

	for _, file := range files{
		if strings.Contains(file.Name(), "cibi.yaml")||strings.Contains(file.Name(), "cibi.yml"){
			f = append(f, file.Name())
		}
	}

	return f
}

func GetCurrDir() string{
	currDir, err := os.Getwd() 
    if err != nil { 
        log.Fatalf("%v", err) 
    } 
    return currDir
}