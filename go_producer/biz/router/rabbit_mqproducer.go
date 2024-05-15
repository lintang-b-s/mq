package router

import (
	"context"
	"errors"
	"fmt"
	"lintang/go_producer/biz/domain"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type RabbitMQProducerService interface {
	SendEmailNotificationsToRMQ(ctx context.Context, message string) error
}

type ProducerHandler struct {
	svc RabbitMQProducerService
}

func ProducerRouter(r *server.Hertz, ps RabbitMQProducerService) {
	handler := &ProducerHandler{
		ps,
	}

	root := r.Group("/api/v1")
	{
		ctrH := root.Group("/rabbitmq")
		{
			ctrH.POST("/", handler.SendMessageToRabbitMQ)
		}
	}
}

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

type sendNotificationReq struct {
	Message string `json:"message"`
}

type sendNotificationRes struct {
	Message string `json:"message"`
}

func (p *ProducerHandler) SendMessageToRabbitMQ(ctx context.Context, c *app.RequestContext) {
	var req sendNotificationReq

	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, ResponseError{Message: err.Error()})
		return
	}

	err = p.svc.SendEmailNotificationsToRMQ(ctx, req.Message)
	if err != nil {
		c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, sendNotificationRes{fmt.Sprintf("message '%s' send to rabbitmq", req.Message)})
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	var ierr *domain.Error
	if !errors.As(err, &ierr) {
		return http.StatusInternalServerError
	} else {
		switch ierr.Code() {
		case domain.ErrInternalServerError:
			return http.StatusInternalServerError
		case domain.ErrNotFound:
			return http.StatusNotFound
		case domain.ErrConflict:
			return http.StatusConflict
		case domain.ErrBadParamInput:
			return http.StatusBadRequest
		default:
			return http.StatusInternalServerError
		}
	}

}
