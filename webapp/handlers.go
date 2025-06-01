package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func genSTSCredsHandler(c *gin.Context) {
	var genSTSJson struct {
		ObjectName string `json:"object_name"`
	}

	if err := c.BindJSON(&genSTSJson); err != nil {
		return
	}

	stsCreds, err := generateSTSCredentials(genSTSJson.ObjectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "generate failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_key_id":     stsCreds.AccessKeyID,
		"secret_access_key": stsCreds.SecretAccessKey,
		"session_token":     stsCreds.SessionToken,
	})
	return
}

func getObjectHandler(c *gin.Context) {
	var getObjectJson struct {
		AccessKeyID     string `json:"access_key_id"`
		SecretAccessKey string `json:"secret_access_key"`
		SessionToken    string `json:"session_token"`
		ObjectName      string `json:"object_name"`
	}

	getObjectJson.AccessKeyID = c.Query("access_key_id")
	getObjectJson.SecretAccessKey = c.Query("secret_access_key")
	getObjectJson.SessionToken = c.Query("session_token")
	getObjectJson.ObjectName = c.Query("object_name")

	object, err := getObject(
		&stsCredentialsValue{
			AccessKeyID:     getObjectJson.AccessKeyID,
			SecretAccessKey: getObjectJson.SecretAccessKey,
			SessionToken:    getObjectJson.SessionToken,
		},
		bucketName,
		getObjectJson.ObjectName,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get object failed"})
		return
	}
	defer object.Close()

	stat, err := object.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get object failed"})
		return
	}

	c.Header("Content-Type", stat.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", stat.Size))

	if _, err = io.Copy(c.Writer, object); err != nil {
		return
	}

	return
}

func putObjectHandler(c *gin.Context) {
	accessKeyID := c.PostForm("access_key_id")
	secretAccessKey := c.PostForm("secret_access_key")
	sessionToken := c.PostForm("session_token")

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "get file failed"})
		return
	}
	objectName := fileHeader.Filename

	avatar, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "get file failed"})
		return
	}
	defer avatar.Close()

	err = putObject(
		&stsCredentialsValue{
			accessKeyID,
			secretAccessKey,
			sessionToken,
		},
		bucketName,
		objectName,
		avatar,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "put object failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "put object success",
		"object-name": objectName,
	})
	return
}
