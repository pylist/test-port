package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// warmUpConnection 先尝试连接，不计算延迟，用于预热网络连接
func WarmUpConnection(address string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

type CheckTCPPortResult struct {
	Open    bool
	Latency time.Duration
}

// checkTCPPort 测试指定的地址和端口是否可以建立 TCP 连接，并返回连接延迟。
func CheckTCPPort(address string, timeout time.Duration) (CheckTCPPortResult, error) {
	result := CheckTCPPortResult{}
	// 先进行一次预连接
	if err := WarmUpConnection(address, timeout); err != nil {
		return result, err
	}

	// 现在测量实际连接延迟
	startTime := time.Now()
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return result, err
	}
	defer conn.Close()
	result = CheckTCPPortResult{
		Open:    true,
		Latency: time.Since(startTime),
	}
	return result, nil
}

type TestPortResult struct {
	Address string `json:"address"`
	Open    bool   `json:"open"`
	Latency string `json:"latency"`
	Message string `json:"message"`
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/testPort", func(c *gin.Context) {
		var req struct {
			Address string `json:"address"`
			Timeout int    `json:"timeout"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		timeout := time.Duration(req.Timeout) * time.Second
		log.Println(timeout, req.Address)
		testResult, err := CheckTCPPort(req.Address, timeout)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result := TestPortResult{
			Address: req.Address,
			Open:    testResult.Open,
			Latency: testResult.Latency.String(),
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    "SUCCESS",
			"message": "",
			"data":    result,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
