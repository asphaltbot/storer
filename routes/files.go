package routes

import (
	"encoding/json"
	"github.com/asphaltbot/file-storage/util"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type URLDownload struct {
	URL string `json:"url"`
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

func RegisterUploadRoutes(e *gin.Engine) {
	e.GET("/file/:id", FetchFileByID)
	e.POST("/upload", UploadFile)
	e.POST("/download", DownloadFile)
	e.DELETE("/file/:id", DeleteFileByID)
}

func DownloadFile(c *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		c.AbortWithStatusJSON(200, gin.H{"code": 500, "message": "Unable to get request body"})
		return
	}

	var urlDownload URLDownload
	err = json.Unmarshal(bodyBytes, &urlDownload)

	if err != nil {
		c.AbortWithStatusJSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	var dir string

	if util.IsRunningInProd() {
		dir = "/home/storage/user/"
	} else {
		dir = "D:\\"
	}

	resp, err := httpClient.Get(urlDownload.URL)

	if err != nil {
		c.AbortWithStatusJSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	defer resp.Body.Close()
	id := util.RandomString(8)

	contentType := resp.Header.Get("content-type")

	if contentType == "" {
		c.AbortWithStatusJSON(200, gin.H{"code": 400, "message": "Content type could not be found for " + urlDownload.URL})
		return
	}

	extension := filepath.Ext(urlDownload.URL)

	if extension == "" {
		c.AbortWithStatusJSON(200, gin.H{"code": 500, "message": "Unable to get extension for: " + urlDownload.URL})
		return
	}

	extension = strings.Split(extension, "?")[0]

	out, err := os.Create(dir + id + extension)

	if err != nil {
		c.AbortWithStatusJSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		c.AbortWithStatusJSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "message": "Successfully downloaded", "id": id})

}

func UploadFile(c *gin.Context) {

	id := util.RandomString(8)
	form, _ := c.MultipartForm()
	file := form.File["file"][0]

	if file == nil {
		c.AbortWithStatusJSON(200, gin.H{"code": 400, "message": "Please upload a file"})
		return
	}

	var dir string

	if util.IsRunningInProd() {
		dir = "/home/storage/user/"
	} else {
		dir = "D:\\"
	}

	err := c.SaveUploadedFile(file, dir+id+filepath.Ext(file.Filename))

	if err != nil {
		c.AbortWithStatusJSON(200, gin.H{"code": 500, "message": "Unable to save uploaded file: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "id": id})

}

func DeleteFileByID(c *gin.Context) {
	fileID := c.Param("id")
	var dir string

	if util.IsRunningInProd() {
		dir = "/home/storage/user/"
	} else {
		dir = "D:\\"
	}

	fileMatch, err := filepath.Glob(dir + fileID + ".*")

	if err != nil {
		c.AbortWithStatusJSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	if len(fileMatch) == 0 {
		c.AbortWithStatusJSON(200, gin.H{"code": 404, "message": "That file could not be found"})
		return
	}

	err = os.Remove(fileMatch[0])

	if err != nil {
		c.AbortWithStatusJSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "message": "Successfully deleted"})

}

func FetchFileByID(c *gin.Context) {
	fileID := c.Param("id")
	var dir string

	if util.IsRunningInProd() {
		dir = "/home/storage/user/"
	} else {
		dir = "D:\\"
	}

	fileMatch, err := filepath.Glob(dir + fileID + ".*")

	if err != nil {
		c.AbortWithStatusJSON(200, gin.H{"code": 500, "message": err.Error()})
		return
	}

	if len(fileMatch) == 0 {
		c.AbortWithStatusJSON(200, gin.H{"code": 404, "message": "That file could not be found"})
		return
	}

	extension := path.Ext(fileMatch[0])

	if extension == ".png" || extension == ".jpg" || extension == ".bmp" || extension == ".gif" {
		c.File(fileMatch[0])
	} else {
		c.FileAttachment(fileMatch[0], fileID+extension)
	}
}
