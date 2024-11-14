package messages

import (
	"bse-eti-stream/pkg/config"
	"bse-eti-stream/pkg/logger"
	"bse-eti-stream/pkg/models"
	"bse-eti-stream/pkg/utils"
	"context"
	"errors"
	"fmt"
	"net"

	"go.uber.org/zap"
)

func getUserLoginRequest() models.UserLoginRequest {
	request := models.UserLoginRequest{}
	request.MHeader = models.MessageRequestHeader{BodyLen: UserLoginPktLen, TemplateId: UserLoginRequestTemplateId}
	request.RHeader = models.RequestHeader{MsgSeqNum: utils.GetSequenceNumber()}
	request.Username = uint32(config.GetConfig().GetInt("credentials.sessionid")) //PROD userId
	copy(request.Password[:8], []byte(config.GetConfig().GetString("credentials.password")))
	logger.Log().Info("User-Login packet created:-", zap.Any("data", request))

	return request
}

func UserLogin(c context.Context, conn net.Conn) error {
	logger.Log().Debug("User-Login Started")
	defer logger.Log().Debug("User-Login Ended")
	var (
		header      models.MessageResponseHeader
		response    models.UserLoginResponseObject
		rejectedPkt models.RejectObject
		err         error
		reason      string
	)
	request := getUserLoginRequest()
	logger.Log().Info("User-Login packet created:- 2    ###############", zap.Any("data", request))
	if err := utils.Send(conn, request, USER_LOGIN); err != nil {
		logger.Log().Error("Failed to send the request", zap.Error(err))
		return err
	}
	fmt.Println("==================header", &header)
	if _, err := utils.Recv(conn, &header, 8, USER_LOGIN); 
	err != nil {

		fmt.Println("#################", &header)
		logger.Log().Error("Failed to read the header", zap.Error(err))
		return err
	}
	logger.Log().Error("Success to receive the response")

	if header.TemplateId == RejectedTemplateId {
		rejectedPkt, reason, err = parseRejectedPkt(conn, header, USER_LOGIN)
		logger.Log(c).Info("rejected", zap.Any("data", rejectedPkt), zap.String("reason", reason), zap.Error(err))
		if err != nil {
			return errors.New("FailedToParseRejecetedPkt")
		}
		return errors.New(reason)
	} else if header.TemplateId == UserLoginResponseTemplateId {
		response, err = parseUserLoginResponsePkt(conn, header, USER_LOGIN)
		if err != nil {
			logger.Log().Error("Failed to parse user-login-response ", zap.Error(err))
			return err
		}
	} else {
		logger.Log().Error("Unknown packet received", zap.Any("headers", header))
		return errors.New("unknown Packet")
	}
	logger.Log().Info("User-Login request successful", zap.Any("response", response))
	return nil
}

func parseUserLoginResponsePkt(conn net.Conn, header models.MessageResponseHeader, mType string) (models.UserLoginResponseObject, error) {
	var (
		body     models.UserLoginResponseBody
		response models.UserLoginResponseObject
	)
	if _, err := utils.Recv(conn, &body, int(header.BodyLen-ResponseHeaderSize), USER_LOGIN); err != nil {
		logger.Log().Error("Failed to read the user-login-body", zap.Error(err), zap.String("mType", mType))
		return response, err
	}
	response.Header = header
	response.Body = body
	return response, nil
}
