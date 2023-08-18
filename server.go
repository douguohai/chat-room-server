package main

import (
	"chat-room/base"
	"chat-room/middlewares"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	lksdk "github.com/livekit/server-sdk-go"
	"log"
	"net/http"
	"time"
)

const apiKey = "devkey"
const apiSecret = "secret"

// getJoinToken 获取房间的token
// room 房间号
// identity 用户id
func getJoinToken(roomId, identity string) (string, error) {
	at := auth.NewAccessToken(apiKey, apiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     roomId,
	}
	at.AddGrant(grant).
		SetIdentity(identity).
		SetValidFor(time.Hour)

	return at.ToJWT()
}

// startRecode 开始录制
// ctx 上下文
// roomId 房间号
func startRecode(ctx *gin.Context, roomId string) {
	egressClient := lksdk.NewEgressClient("http://localhost:7880", apiKey, apiSecret)
	fileRequest := &livekit.RoomCompositeEgressRequest{
		RoomName: roomId,
		Layout:   "speaker",
		Output: &livekit.RoomCompositeEgressRequest_File{
			File: &livekit.EncodedFileOutput{
				FileType: livekit.EncodedFileType_MP4,
				Filepath: "livekit-demo/my-room-test.mp4",
				Output: &livekit.EncodedFileOutput_S3{
					S3: &livekit.S3Upload{
						Region:    "zh-east-1",
						AccessKey: "xtm",
						Secret:    "Xtm@123456",
						Bucket:    "mjts",
						Endpoint:  "http://192.168.10.240:9000",
					},
				},
			},
		},
	}
	result, err := egressClient.StartRoomCompositeEgress(ctx, fileRequest)
	if err != nil {
		log.Printf("startRecode err: %v", err)
		return
	}
	fmt.Println(result)
}

func closeRecode(ctx *gin.Context, egressID string) {
	egressClient := lksdk.NewEgressClient("http://localhost:7880", apiKey, apiSecret)
	info, err := egressClient.StopEgress(ctx, &livekit.StopEgressRequest{
		EgressId: egressID,
	})
	if err != nil {
		log.Printf("endRecode err: %v", err)
		return
	}
	fmt.Println(info)
}

func main() {
	r := gin.Default()
	r.Use(middlewares.Cors())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/getToken/:roomId/:userId", func(c *gin.Context) {
		roomId := c.Param("roomId")
		userId := c.Param("userId")
		log.Printf("roomId userId: %v,%v", roomId, userId)
		token, _ := getJoinToken(roomId, userId)
		c.JSON(200,
			base.Result{
				Code:    200,
				Message: token,
			})
	})

	r.GET("/recode/:roomId", func(c *gin.Context) {
		roomId := c.Param("roomId")
		log.Printf("roomId : %v", roomId)
		startRecode(c, roomId)
		c.JSON(200,
			base.Result{
				Code:    200,
				Message: "ok",
			})
	})

	r.GET("/end/:egressID", func(c *gin.Context) {
		egressID := c.Param("egressID")
		log.Printf("egressID : %v", egressID)
		closeRecode(c, egressID)
		c.JSON(200,
			base.Result{
				Code:    200,
				Message: "ok",
			})
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
	log.Fatal(http.ListenAndServe(":8080", nil))
}
