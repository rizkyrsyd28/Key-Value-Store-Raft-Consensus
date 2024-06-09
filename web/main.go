package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Sister20/if3230-tubes-dark-syster/lib/client"
	"github.com/Sister20/if3230-tubes-dark-syster/lib/pb"
	. "github.com/Sister20/if3230-tubes-dark-syster/lib/util"
	"github.com/gin-gonic/gin"
)

type KeyValueRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type KeyRequest struct {
	Key string `json:"key" binding:"required"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run ./web/client.go <server ip> <server port>")
		return
	}

	router := gin.Default()

	router.LoadHTMLGlob("web/templates/*")

	address := NewAddress(os.Args[1], os.Args[2])

	client, err := client.NewClient(address)
	if err != nil {
		fmt.Println("Address Not Found")
	}

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Dark Syster",
		})
	})
	router.PUT("/network", ChangeNetwork(client))
	router.POST("/ping", Ping(client))
	router.POST("/get", Get(client))
	router.POST("/set", Set(client))
	router.POST("/strln", StrLn(client))
	router.POST("/del", Del(client))
	router.POST("/append", Append(client))

	router.Run(":3000")
}

func ChangeNetwork(c *client.GRPCClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var rawAddress struct {
			IP   string `json:"ip" binding:"required"`
			Port string `json:"port" binding:"required"`
		}

		if err := ctx.ShouldBindJSON(&rawAddress); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Body"})
			return
		}

		address := NewAddress(rawAddress.IP, rawAddress.Port)

		if err := c.SetAddress(address); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cant Change Network"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
	}
}

func Ping(c *client.GRPCClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response, err := c.Services.KV.Ping(ctx, &pb.Empty{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something Wrong With Server"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": response.GetValue()})
	}
}

func Get(c *client.GRPCClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var input KeyRequest

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Body"})
			return
		}

		response, err := c.Services.KV.Get(ctx, &pb.KeyRequest{Key: input.Key})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something Wrong With Server"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": response.GetValue()})
	}
}

func Set(c *client.GRPCClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var input KeyValueRequest

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Body"})
			return
		}

		response, err := c.Services.KV.Set(ctx, &pb.KeyValueRequest{Key: input.Key, Value: input.Value})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something Wrong With Server"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": response.GetValue()})
	}
}

func StrLn(c *client.GRPCClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var input KeyRequest

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Body"})
			return
		}

		response, err := c.Services.KV.StrLn(ctx, &pb.KeyRequest{Key: input.Key})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something Wrong With Server"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": response.GetValue()})
	}
}

func Del(c *client.GRPCClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var input KeyRequest

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Body"})
			return
		}

		response, err := c.Services.KV.Del(ctx, &pb.KeyRequest{Key: input.Key})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something Wrong With Server"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": response.GetValue()})
	}
}

func Append(c *client.GRPCClient) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var input KeyValueRequest

		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Body"})
			return
		}

		response, err := c.Services.KV.Append(ctx, &pb.KeyValueRequest{Key: input.Key, Value: input.Value})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something Wrong With Server"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": response.GetValue()})
	}
}
