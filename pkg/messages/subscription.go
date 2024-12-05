package messages

import (
	"bse-eti-stream/pkg/logger"
	"bse-eti-stream/pkg/models"
	"bse-eti-stream/pkg/utils"
	"context"
	"errors"
	"fmt"
	"net"

	"go.uber.org/zap"
)

func getSubscriptionRequest() models.SubscribeRequest {
	request := models.SubscribeRequest{}
	request.MHeader = models.MessageRequestHeader{BodyLen: SubscriptionPktLen, TemplateId: SubscriptionRequestTemplateId}
	request.RHeader = models.RequestHeader{MsgSeqNum: utils.GetSequenceNumber()}
	request.SubscriptionScope = uint32(4294967295) //not value to fetch all
	request.RefAppId = uint8(1)
	logger.Log().Info("Subscription packet created:-", zap.Any("data", request))
	return request
}

func Subscribe(c context.Context, conn net.Conn) error {
	logger.Log().Debug("Subscription Started")
	defer logger.Log().Debug("Subscription Ended")
	var (
		header      models.MessageResponseHeader
		response    models.SubscribeResponseObject
		rejectedPkt models.RejectObject
		err         error
		reason      string
	)
	request := getSubscriptionRequest()
	if err := utils.Send(conn, request, SUBSCRIPTION); err != nil {
		logger.Log().Error("Failed to send the request", zap.Error(err))
		return err
	}
	if _, err := utils.Recv(conn, &header, 8, SUBSCRIPTION); 
	
	err != nil {
		logger.Log().Error("Failed to read the header", zap.Error(err))
		return err
	}
	if header.TemplateId == RejectedTemplateId {
		rejectedPkt, reason, err = parseRejectedPkt(conn, header, SUBSCRIPTION)
		logger.Log(c).Info("rejected", zap.Any("data", rejectedPkt), zap.String("reason", reason), zap.Error(err))
		if err != nil {
			return errors.New("FailedToParseRejecetedPkt")
		}
		return errors.New(reason)
	} else if header.TemplateId == SubscriptionResponseTemplateId {
		fmt.Println("Header templateid:", header.TemplateId)
		response, err = parseSubscriptionResponsePkt(conn, header, SUBSCRIPTION)
		if err != nil {
			logger.Log().Error("Failed to parse the ", zap.Error(err))
			return err
		}
	} else {
		logger.Log().Error("Unknown packet received", zap.Any("headers", header))
		return errors.New("unknown Packet")
	}

	logger.Log().Info("Subscription request successful", zap.Any("response", response))
	return nil
}

func parseSubscriptionResponsePkt(conn net.Conn, header models.MessageResponseHeader, mType string) (models.SubscribeResponseObject, error) {
	var (
		body     models.SubscribeResponseBody
		response models.SubscribeResponseObject
	)
	fmt.Println("mmmmmmmmmmmmm", body)
	if _, err := utils.Recv(conn, &body, int(header.BodyLen-ResponseHeaderSize), SUBSCRIPTION); err != nil {
		logger.Log().Error("Failed to read the Subscription-body", zap.Error(err), zap.String("mType", mType))
		return response, err
	}
	response.Header = header
	response.Body = body
	return response, nil
}
