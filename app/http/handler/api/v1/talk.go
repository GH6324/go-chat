package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/cache"
	"go-chat/app/http/dto"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/strutil"
	"go-chat/app/service"
	"strings"
)

type Talk struct {
	service         *service.TalkService
	talkListService *service.TalkListService
	redisLock       *cache.RedisLock
}

func NewTalkHandler(
	service *service.TalkService,
	talkListService *service.TalkListService,
	redisLock *cache.RedisLock,
) *Talk {
	return &Talk{service, talkListService, redisLock}
}

// List 会话列表
func (c *Talk) List(ctx *gin.Context) {

}

// Create 创建会话列表
func (c *Talk) Create(ctx *gin.Context) {
	params := &request.TalkListCreateRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := auth.GetAuthUserID(ctx)

	agent := strings.TrimSpace(ctx.GetHeader("user-agent"))
	if agent != "" {
		agent = strutil.Md5([]byte(agent))
	}

	lockKey := fmt.Sprintf("talk:list:%d-%d-%d-%s", uid, params.ReceiverId, params.TalkType, agent)
	if !c.redisLock.Lock(ctx.Request.Context(), lockKey, 20) {
		response.BusinessError(ctx, "创建失败")
		return
	}

	result, err := c.talkListService.Create(ctx.Request.Context(), uid, params)

	if err != nil {
		response.BusinessError(ctx, err.Error())
		return
	}

	response.Success(ctx, dto.TalkListItem{
		ID:         result.ID,
		TalkType:   result.TalkType,
		ReceiverId: result.ReceiverId,
		IsRobot:    result.IsRobot,
		Avatar:     "",
		Name:       "",
		RemarkName: "",
		UnreadNum:  0,
		MsgText:    "",
		UpdatedAt:  "",
	})
}

// Delete 删除列表
func (c *Talk) Delete(ctx *gin.Context) {
	params := &request.TalkListDeleteRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.talkListService.Delete(ctx, auth.GetAuthUserID(ctx), params.Id); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{})
}

// Top 置顶列表
func (c *Talk) Top(ctx *gin.Context) {
	params := &request.TalkListTopRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.talkListService.Top(ctx, auth.GetAuthUserID(ctx), params); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{})
}

// Disturb 会话免打扰
func (c *Talk) Disturb(ctx *gin.Context) {
	params := &request.TalkListDisturbRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.talkListService.Disturb(ctx, auth.GetAuthUserID(ctx), params); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{})
}
