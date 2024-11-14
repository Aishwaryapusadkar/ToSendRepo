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

func getConnectionRequest() models.ConnectionGatewayRequest {
	request := models.ConnectionGatewayRequest{}
	request.MHeader = models.MessageRequestHeader{BodyLen: ConnectionGatewayPktLen, TemplateId: ConnetionGatewayRequestTemplateId}
	request.RHeader = models.RequestHeader{MsgSeqNum: utils.GetSequenceNumber()}
	copy(request.AppVersion[:3], []byte("2.3"))
	copy(request.Password[:8], []byte(config.GetConfig().GetString("credentaals.password")))
	request.SessionId = uint32(config.GetConfig().GetInt("credentials.sessionid")) //PROD userId
	logger.Log().Info("connection-request packet created:-", zap.Any("data", request))
	return request
}

func ConnectToGateway(c context.Context, conn net.Conn) (string, string, error) {
	logger.Log().Debug("Connection-Request Started")
	defer logger.Log().Debug("Connection-Request Ended")
	var (
		header                     models.MessageResponseHeader
		primaryConn, secondaryConn string
		response                   models.ConnectionGatewayResponseObject
		err                        error
	)

	request := getConnectionRequest()
	if err := utils.Send(conn, request, CONNECTION_GATEWAY); err != nil {
		logger.Log().Error("Failed to send the request", zap.Error(err))
		return primaryConn, secondaryConn, err
	}
	if _, err := utils.Recv(conn, &header, 8, CONNECTION_GATEWAY); err != nil {
		logger.Log().Error("Failed to read the header", zap.Error(err))
		return primaryConn, secondaryConn, err
	}
	if header.TemplateId == RejectedTemplateId {
		rejectedPkt, reason, err := parseRejectedPkt(conn, header, CONNECTION_GATEWAY)
		logger.Log(c).Info("rejected", zap.Any("data", rejectedPkt), zap.String("reason", reason), zap.Error(err))
		if err != nil {
			return primaryConn, secondaryConn, errors.New("FailedToParseRejecetedPkt")
		}
		return primaryConn, secondaryConn, errors.New(reason)
	} else if header.TemplateId == ConnetionGatewayResponseTemplateId {
		response, err = parseConnectionPkt(conn, header, CONNECTION_GATEWAY)
		if err != nil {
			logger.Log().Error("Failed to parse connnection-gateway-response ", zap.Error(err))
			return primaryConn, secondaryConn, err
		}
	} else {
		logger.Log().Error("Unknown packet received", zap.Any("headers", header))
		return primaryConn, secondaryConn, errors.New("unknown Packet")
	}

	primaryConn = utils.IntToIPV4(response.Body.GatewayId[:], response.Body.GatwaysubId)
	secondaryConn = utils.IntToIPV4(response.Body.SgatewayId[:], response.Body.SgatwaysubId)

	logger.Log(c).Info("Connection gateway request successful", zap.Any("response", response), zap.String("primaryConn", primaryConn), zap.String("secondaryConn", secondaryConn))
	return primaryConn, secondaryConn, nil
}

func parseConnectionPkt(conn net.Conn, header models.MessageResponseHeader, mType string) (models.ConnectionGatewayResponseObject, error) {
	var (
		body     models.ConnectionGatewayResponseBody
		response models.ConnectionGatewayResponseObject
	)
	if _, err := utils.Recv(conn, &body, int(header.BodyLen-ResponseHeaderSize), CONNECTION_GATEWAY); err != nil {
		logger.Log().Error("Failed to read the connection-body", zap.Error(err), zap.String("mType", mType))
		return response, err
	}
	response.Header = header
	response.Body = body
	return response, nil
}
