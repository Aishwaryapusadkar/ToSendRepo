package messages

import (
	"bse-eti-stream/pkg/config"
	"bse-eti-stream/pkg/logger"
	"bse-eti-stream/pkg/models"
	"bse-eti-stream/pkg/utils"
	"context"
	"errors"
	"net"

	"go.uber.org/zap"
)

func getSessionLoginRequest() models.SessionLoginRequest {
	request := models.SessionLoginRequest{}
	request.MHeader = models.MessageRequestHeader{BodyLen: SessionLoginPktLen, TemplateId: SessionLoginRequestTemplateId}
	request.RHeader = models.RequestHeader{MsgSeqNum: uint32(1)}
	request.SessionId = uint32(config.GetConfig().GetInt("credentials.sessionid")) //PROD userId
	copy(request.AppVersion[:3], []byte("2.3"))
	copy(request.Password[:8], []byte(config.GetConfig().GetString("credentials.password")))
	copy(request.AppUsageOrders[:1], []byte("B"))
	copy(request.AppUsageQuotes[:1], []byte("N"))
	copy(request.OrderRoutingIndicator[:1], []byte("N"))
	copy(request.AppSystemName[:9], []byte("VendorIML"))
	copy(request.AppSystemVendor[:8], []byte("Internal"))
	copy(request.AppSystemVersion[:3], []byte("2.3"))
	logger.Log().Info("session-login packet created:-", zap.Any("data", request))
	return request
}

func SessionLogin(c context.Context, conn net.Conn) error {
	logger.Log().Debug("Session-Login Started")
	defer logger.Log().Debug("Session-Login Ended")
	var (
		header      models.MessageResponseHeader
		response    models.SessionLoginResponseObject
		rejectedPkt models.RejectObject
		err         error
		reason      string
	)
	request := getSessionLoginRequest()
	if err := utils.Send(conn, request, SESSION_LOGIN); err != nil {
		logger.Log().Error("Failed to send the request", zap.Error(err))
		return err
	}
	if _, err := utils.Recv(conn, &header, 8, SESSION_LOGIN); err != nil {
		logger.Log().Error("Failed to read the header", zap.Error(err))
		return err
	}
	if header.TemplateId == RejectedTemplateId {
		rejectedPkt, reason, err = parseRejectedPkt(conn, header, SESSION_LOGIN)
		logger.Log(c).Info("rejected", zap.Any("data", rejectedPkt), zap.String("reason", reason), zap.Error(err))
		if err != nil {
			return errors.New("FailedToParseRejecetedPkt")
		}
		return errors.New(reason)
	} else if header.TemplateId == SessionLoginResponseTemplateId {
		response, err = parseSessionLoginResponsePkt(conn, header, SESSION_LOGIN)
		if err != nil {
			logger.Log().Error("Failed to parse session response", zap.Error(err))
			return err
		}
	} else {
		logger.Log().Error("Unknown packet received", zap.Any("headers", header))
		return errors.New("unknown Packet")
	}
	logger.Log().Info("Session-Login request successful", zap.Any("response", response))
	return nil
}

func parseSessionLoginResponsePkt(conn net.Conn, header models.MessageResponseHeader, mType string) (models.SessionLoginResponseObject, error) {
	var (
		body     models.SessionLoginResponseBody
		response models.SessionLoginResponseObject
	)
	if _, err := utils.Recv(conn, &body, int(header.BodyLen-ResponseHeaderSize), SESSION_LOGIN); err != nil {
		logger.Log().Error("Failed to read the session-login-body", zap.Error(err), zap.String("mType", mType))
		return response, err
	}
	response.Header = header
	response.Body = body
	return response, nil
}
