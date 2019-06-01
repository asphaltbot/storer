package routes

import (
	"github.com/asphaltbot/file-storage/util"
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"path/filepath"
)

func RegisterUploadRoutes(e *gin.Engine) {
	e.GET("/file/:id", FetchFileByID)
	e.POST("/upload", UploadFile)
	e.DELETE("/file/:id", DeleteFileByID)
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


	err := c.SaveUploadedFile(file, dir + id + filepath.Ext(file.Filename))

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

	if len(fileMatch) == 0{
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

	if len(fileMatch) == 0{
		c.AbortWithStatusJSON(200, gin.H{"code": 404, "message": "That file could not be found"})
		return
	}

	extension := path.Ext(fileMatch[0])
	c.FileAttachment(fileMatch[0], fileID + extension)

}
