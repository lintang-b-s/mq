package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type KafkaProducerService interface {
	SendEmailNotificationsToKafka(ctx context.Context, message string) error
}

type KafkaProducerHandler struct {
	svc KafkaProducerService
}

func KafkaProducerRouter(r *server.Hertz, ps KafkaProducerService) {
	handler := &KafkaProducerHandler{
		ps,
	}

	root := r.Group("/api/v1")
	{
		ctrH := root.Group("/kafka")
		{
			ctrH.POST("/", handler.SendMessageToKafka)
		}
	}
}

func (p *KafkaProducerHandler) SendMessageToKafka(ctx context.Context, c *app.RequestContext) {
	var req sendNotificationReq

	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	err = p.svc.SendEmailNotificationsToKafka(ctx, req.Message)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, sendNotificationRes{fmt.Sprintf("message '%s' send to kafka", req.Message)})
}
