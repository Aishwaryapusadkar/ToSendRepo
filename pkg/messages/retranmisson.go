package messages

import (
	"bse-eti-stream/pkg/logger"
	"bse-eti-stream/pkg/models"
	"bse-eti-stream/pkg/utils"
	"errors"
	"net"

	"go.uber.org/zap"
)

func getRetransmissionRequest(userId uint32, partitionData models.PartitionInfo) models.ReTransmissionRequest {
	request := models.ReTransmissionRequest{}
	request.MHeader = models.MessageRequestHeader{BodyLen: RetransmissionPktLen, TemplateId: SessionLoginRequestTemplateId}
	request.RHeader = models.RequestHeader{MsgSeqNum: utils.GetSequenceNumber()}
	request.Scope = userId //For which TradeId retranmission is required
	request.RefAppId = uint8(1)
	request.PartionId = partitionData.Id
	request.AppBegnSeqno = partitionData.BegSeqNo
	request.AppEndSeqno = partitionData.EndSeqno
	logger.Log().Info("Retransmission packet created:-", zap.Any("data", request))
	return request
}

func Retransmit(conn net.Conn, userId uint32, paritionData models.PartitionInfo) error {
	logger.Log().Debug("Retransmission Started")
	defer logger.Log().Debug("Retransmission Ended")
	var (
		header   models.MessageResponseHeader
		response interface{}
		err      error
		reason   string
	)
	request := getRetransmissionRequest(userId, paritionData)
	if err := utils.Send(conn, request, RETRANSMISSION); err != nil {
		logger.Log().Error("Failed to send the request", zap.Error(err))
		return err
	}
	if _, err := utils.Recv(conn, &header, 8, RETRANSMISSION); err != nil {
		logger.Log().Error("Failed to read the header", zap.Error(err))
		return err
	}
	if header.TemplateId == RejectedTemplateId {
		response, reason, err = parseRejectedPkt(conn, header, RETRANSMISSION)
		logger.Log().Info("rejected", zap.Any("data", response), zap.String("reason", reason), zap.Error(err))
		if err != nil {
			return errors.New("FailedToParseRejecetedPkt")
		}
		return errors.New(reason)
	} else if header.TemplateId == RetransmissionResponseTemplateId {
		response, err = parseRetransmissionResponsePkt(conn, header, RETRANSMISSION)
	} else {
		logger.Log().Error("Unknown packet received", zap.Any("headers", header))
		return errors.New("unknown Packet")
	}
	if err != nil {
		logger.Log().Error("Failed to parse the ", zap.Error(err))
		return err
	}
	logger.Log().Info("Retransmission request successful", zap.Uint32("userid", userId), zap.Any("parition-info", paritionData), zap.Any("response", response))
	return nil
}

func parseRetransmissionResponsePkt(conn net.Conn, header models.MessageResponseHeader, mType string) (models.ReTransmissionResponseObject, error) {
	var (
		body     models.ReTransmissionResponseBody
		response models.ReTransmissionResponseObject
	)
	if _, err := utils.Recv(conn, &body, int(header.BodyLen-ResponseHeaderSize), RETRANSMISSION); err != nil {
		logger.Log().Error("Failed to read the Retransmission-body", zap.Error(err), zap.String("mType", mType))
		return response, err
	}
	response.Header = header
	response.Body = body
	return response, nil
}
