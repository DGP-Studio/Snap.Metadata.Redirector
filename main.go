package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"time"
)

var currentHost string

func JiHuHasBannedFiles() {
	CdnHost := os.Getenv("cdn_host")
	JiHuHost := os.Getenv("JiHuLab_host")

	for {
		fmt.Printf("Checking JiHuBan Status...")
		apiURL := "https://api.github.com/repos/Masterain98/JiHuBanChecker/issues?labels=legal&state=open"
		GitHubToken := os.Getenv("github_token")
		req, err := http.NewRequest("GET", apiURL, nil)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("Authorization", "Bearer "+GitHubToken)
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}

		var issues []map[string]interface{}
		err = json.Unmarshal(body, &issues)

		if len(issues) > 0 {
			currentHost = CdnHost
			fmt.Printf("Legal issues found")
		} else {
			currentHost = JiHuHost
			fmt.Printf("No legal issue found")
		}

		time.Sleep(5 * time.Minute)
	}
}

func main() {
	r := gin.Default()
	go JiHuHasBannedFiles()

	r.GET("/*path", func(c *gin.Context) {
		if c.Param("path") == "online" {
			c.JSON(http.StatusOK, gin.H{
				"message": "redirect server is running",
			})
		} else if c.Param("path") == "current-target" {
			c.JSON(200, gin.H{"host": currentHost})
		} else {
			hostBURL := "https://" + currentHost + c.Param("path")
			c.Redirect(302, hostBURL)
		}
	})

	if err := r.Run(":8080"); err != nil {
		fmt.Println("Gin server encountered an error:", err)
	}
}
