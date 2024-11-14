package messages

import (
	"bse-eti-stream/pkg/logger"
	"bse-eti-stream/pkg/models"
	"bse-eti-stream/pkg/utils"
	"net"

	"go.uber.org/zap"
)

func parseRejectedPkt(conn net.Conn, header models.MessageResponseHeader, mType string) (models.RejectObject, string, error) {
	var (
		body      models.RejectBody
		pkt       models.RejectObject
		bodyBytes []byte
		err       error
	)
	logger.Log().Debug("rejected header", zap.Any("data", header), zap.String("mtype", mType))

	if bodyBytes, err = utils.Recv(conn, &body, int(header.BodyLen-ResponseHeaderSize), mType); err != nil {
		logger.Log().Error("Failed to read the rejected body", zap.Error(err), zap.String("mType", mType))
		return pkt, "", err
	}
	pkt.Header = header
	pkt.Body = body
	temp := make([]byte, pkt.Body.TextLen)
	copy(temp, bodyBytes[64:64+pkt.Body.TextLen])
	return pkt, string(temp), nil
}
