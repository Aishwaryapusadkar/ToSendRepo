package messages

import (
	"bse-eti-stream/pkg/config"
	"bse-eti-stream/pkg/logger"
	"bse-eti-stream/pkg/models"
	"bse-eti-stream/pkg/utils"

	//"bse-eti-stream/app/database"

	//"bytes"
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
)

var db *sql.DB

//db *sql.DB

func StartFeedCapturing(c context.Context, conn net.Conn) error {

	logger.Log().Debug("Feed-caputring Started")
	defer logger.Log().Debug("Feed-caputring Ended")
	var (
		header models.MessageResponseHeader
	)

	ctx := context.Background()
	utils.CreateWorkerPool(ctx, 5) // Initialize pool before usage
	fmt.Println("###############message header", header.BodyLen, header.TemplateId)
	wp := utils.GetWorkerPool()

	wp.Run(c, tradePktHandler)
	for {
		var (
			response interface{}
			err      error
			index    int
			reason   string
		)
		fmt.Println("###############message header", header.BodyLen, header.TemplateId)

		if _, err := utils.Recv(conn, &header, ResponseHeaderSize, TRADE); err != nil {
			logger.Log().Error("Failed to read the header", zap.Error(err))
			return err
		}

		fmt.Printf("#########################")
		fmt.Printf("TRADE TEMPLATE ID ", TradeTemplateId)
		fmt.Printf("TemplateID", header.TemplateId)

		fmt.Printf("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")

		if header.TemplateId == RejectedTemplateId {
			response, reason, err = parseRejectedPkt(conn, header, TRADE)
			logger.Log(c).Info("rejected", zap.Any("data", response), zap.String("reason", reason), zap.Error(err))

		} else if header.TemplateId == HeartbeatTemplateId {
			logger.Log().Info("HEART-BEAT received")

		} else if header.TemplateId == TradeTemplateId {
			//index = index % wp.Size()
			if err = parseTradePacket(c, index, conn, header, TRADE, wp); err != nil {
				logger.Log().Error("error while parsing trade body, looking for next pkt", zap.Error(err))
			} else {
				index++
			}
		} else if header.TemplateId == SessionLogoutTemplateId {
			logger.Log().Info("session logout")
			break
		} else {
			logger.Log().Error("Unknown packet received", zap.Any("headers", header))
		}

	}
	return nil
}

func parseTradePacket(c context.Context, index int, conn net.Conn, header models.MessageResponseHeader, mType string, wp utils.WorkerPoolInterface) error {
	var (
		body models.TradeBody
		pkt  models.TradeObject
	)

	fmt.Println("header", header)

	if _, err := utils.Recv(conn, &body, int(header.BodyLen-ResponseHeaderSize), mType); err != nil {
		logger.Log().Error("Failed to received trade body", zap.Any("headers", header), zap.String("mType", mType))
		return err
	}
	pkt.Header = header
	pkt.Body = body
	wp.EnQueue(c, index, pkt)
	return nil
}

// func TestStartFeedCapturing(c context.Context, conn net.Conn) error {
// 	logger.Log().Debug("Feed-caputring Started")
// 	defer logger.Log().Debug("Feed-caputring Ended")
// 	var (
// 		header models.MessageResponseHeader
// 	)

// 	i := 0
// 	for {
// 		if i > 3 {
// 			break
// 		}
// 		h := [][]byte{
// 			{96, 1, 0, 0, 4, 41, 0, 0},
// 			{96, 1, 0, 0, 4, 41, 0, 0},
// 			{96, 1, 0, 0, 4, 41, 0, 0},
// 			{96, 1, 0, 0, 4, 41, 0, 0},
// 		}
// 		b := [][]byte{
// 			{60, 55, 165, 181, 88, 86, 195, 23, 101, 42, 0, 0, 0, 0, 0, 0, 167, 10, 0, 0, 4, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 130, 34, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 128, 64, 142, 83, 56, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 128, 49, 52, 142, 181, 88, 86, 195, 23, 194, 84, 80, 152, 174, 68, 195, 23, 4, 57, 107, 147, 191, 188, 19, 0, 85, 68, 4, 0, 0, 0, 0, 0, 88, 76, 142, 181, 88, 86, 195, 23, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 85, 68, 4, 0, 193, 12, 1, 0, 255, 255, 255, 255, 221, 0, 0, 0, 69, 43, 157, 0, 47, 43, 157, 0, 79, 59, 0, 0, 200, 0, 0, 0, 0, 0, 0, 0, 22, 0, 0, 0, 0, 0, 0, 128, 200, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 212, 143, 24, 1, 21, 216, 52, 1, 223, 12, 1, 0, 255, 255, 255, 255, 0, 0, 0, 128, 255, 255, 1, 0, 1, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 30, 4, 255, 1, 1, 1, 65, 49, 67, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 56, 53, 49, 48, 55, 55, 49, 54, 48, 55, 32, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 66, 50, 48, 48, 48, 49, 52, 51, 55, 56, 32, 32, 0, 5, 255, 255, 69, 67, 65, 71, 49, 48, 51, 32, 32, 84, 48, 48, 50, 48, 55, 48, 32, 32, 32, 32, 0, 0, 0, 0, 0, 0, 0},
// 			{6, 203, 120, 190, 88, 86, 195, 23, 102, 42, 0, 0, 0, 0, 0, 0, 167, 10, 0, 0, 4, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 127, 31, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 112, 56, 57, 0, 0, 0, 0, 0, 112, 56, 57, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 128, 0, 69, 100, 190, 88, 86, 195, 23, 74, 223, 76, 152, 174, 68, 195, 23, 4, 57, 107, 147, 191, 188, 19, 0, 47, 189, 11, 0, 0, 0, 0, 0, 129, 222, 16, 55, 124, 82, 195, 23, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 47, 189, 11, 0, 194, 12, 1, 0, 255, 255, 255, 255, 221, 0, 0, 0, 69, 43, 157, 0, 51, 43, 157, 0, 79, 59, 0, 0, 47, 5, 0, 0, 113, 10, 0, 0, 22, 0, 0, 0, 0, 0, 0, 128, 111, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 56, 144, 24, 1, 21, 216, 52, 1, 224, 12, 1, 0, 255, 255, 255, 255, 0, 0, 0, 128, 255, 255, 1, 0, 1, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 30, 4, 255, 2, 0, 1, 65, 49, 67, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 56, 53, 48, 52, 51, 55, 57, 49, 48, 57, 32, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 66, 55, 48, 48, 48, 49, 49, 56, 53, 48, 32, 32, 49, 2, 255, 255, 69, 67, 65, 71, 49, 48, 51, 32, 32, 84, 48, 48, 50, 49, 49, 48, 32, 32, 32, 32, 0, 0, 0, 0, 0, 0, 0},
// 			{37, 57, 193, 212, 88, 86, 195, 23, 33, 24, 0, 0, 0, 0, 0, 0, 167, 10, 0, 0, 4, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 90, 40, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 104, 176, 194, 41, 0, 0, 0, 0, 104, 176, 194, 41, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 128, 5, 31, 171, 212, 88, 86, 195, 23, 244, 188, 76, 152, 174, 68, 195, 23, 4, 57, 107, 147, 191, 188, 19, 0, 152, 81, 7, 0, 0, 0, 0, 0, 166, 45, 171, 212, 88, 86, 195, 23, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 152, 81, 7, 0, 203, 136, 0, 0, 255, 255, 255, 255, 221, 0, 0, 0, 46, 43, 157, 0, 65, 43, 157, 0, 79, 59, 0, 0, 90, 0, 0, 0, 0, 0, 0, 0, 9, 1, 0, 0, 0, 0, 0, 128, 90, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 28, 220, 132, 0, 21, 216, 52, 1, 204, 136, 0, 0, 255, 255, 255, 255, 0, 0, 0, 128, 255, 255, 1, 0, 1, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 30, 4, 255, 2, 1, 1, 65, 49, 67, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 56, 53, 49, 48, 52, 57, 57, 49, 50, 54, 32, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 66, 52, 48, 48, 48, 49, 53, 48, 48, 50, 32, 32, 49, 2, 255, 255, 69, 67, 65, 71, 49, 48, 51, 32, 32, 84, 48, 48, 50, 50, 53, 48, 32, 32, 32, 32, 0, 0, 0, 0, 0, 0, 0},
// 			{29, 72, 247, 226, 88, 86, 195, 23, 173, 21, 0, 0, 0, 0, 0, 0, 167, 10, 0, 0, 2, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 91, 65, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 128, 64, 101, 20, 191, 24, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 128, 205, 65, 227, 226, 88, 86, 195, 23, 122, 237, 133, 228, 250, 85, 195, 23, 4, 57, 107, 147, 191, 188, 19, 0, 157, 81, 7, 0, 0, 0, 0, 0, 64, 82, 227, 226, 88, 86, 195, 23, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 157, 81, 7, 0, 126, 216, 6, 0, 255, 255, 255, 255, 221, 0, 0, 0, 46, 43, 157, 0, 65, 43, 157, 0, 79, 59, 0, 0, 7, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 128, 7, 0, 0, 0, 0, 0, 0, 128, 0, 0, 0, 128, 136, 105, 83, 6, 21, 216, 52, 1, 15, 233, 6, 0, 255, 255, 255, 255, 0, 0, 0, 128, 255, 255, 1, 0, 1, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 30, 4, 255, 1, 1, 1, 65, 49, 67, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 56, 53, 48, 50, 49, 51, 57, 57, 56, 53, 32, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 66, 52, 48, 48, 48, 49, 53, 48, 48, 51, 32, 32, 0, 5, 255, 255, 69, 67, 65, 71, 49, 48, 51, 32, 32, 84, 48, 48, 50, 50, 53, 48, 32, 32, 32, 32, 0, 0, 0, 0, 0, 0, 0},
// 		}
// 		var (
// 			response interface{}
// 			err      error
// 			index    int
// 			reason   string
// 		)
// 		if _, err := utils.TestRecv(h[i], &header, ResponseHeaderSize, TRADE); err != nil {
// 			logger.Log().Error("Failed to read the header", zap.Error(err))
// 			return err
// 		}
// 		if header.TemplateId == RejectedTemplateId {
// 			response, reason, err = parseRejectedPkt(conn, header, TRADE)
// 			logger.Log(c).Info("rejected", zap.Any("data", response), zap.String("reason", reason), zap.Error(err))
// 		} else if header.TemplateId == HeartbeatTemplateId {
// 			logger.Log().Info("HEART-BEAT received")
// 		} else if header.TemplateId == TradeTemplateId {
// 			if err = testparseTradePacket(c, index, b[i], header, TRADE); err != nil {
// 				logger.Log().Error("error while parsing trade body, looking for next pkt", zap.Error(err))
// 			} else {
// 				index++
// 			}
// 		} else if header.TemplateId == SessionLogoutTemplateId {
// 			logger.Log().Info("session logout")
// 			break
// 		} else {
// 			logger.Log().Error("Unknown packet received", zap.Any("headers", header))
// 		}
// 		i++
// 	}
// 	return nil
// }

//	func testparseTradePacket(c context.Context, index int, bytes []byte, header models.MessageResponseHeader, mType string) error {
//		logger.Log(c).Debug("STARTED")
//		defer logger.Log(c).Debug("ENDED")
//		var (
//			body models.TradeBody
//			pkt  models.TradeObject
//		)
//		if _, err := utils.TestRecv(bytes, &body, int(header.BodyLen-ResponseHeaderSize), mType); err != nil {
//			logger.Log().Error("Failed to received trade body", zap.Any("headers", header), zap.String("mType", mType))
//			return err
//		}
//		pkt.Header = header
//		pkt.Body = body
//		_, y := tradePktHandler(pkt)
//		return y
//	}
func tradePktHandler(data interface{}) (interface{}, error) {
	logger.Log().Debug("STARTED")
	defer logger.Log().Debug("ENDED")
	var (
		dataToAdd []string
	)
	cnf := config.GetConfig()
	temp := data.(models.TradeObject)
	fmt.Printf("\nfeed %+v\n", temp)
	dataToAdd = append(dataToAdd, cnf.GetString("credentials.membercode")) //memeber-id
	dataToAdd = append(dataToAdd, getTraderId(temp.Body.SessionId))        //trade-id
	dataToAdd = append(dataToAdd, getTraderId(temp.Body.SessionId))        //session-id
	//dataToAdd = append(dataToAdd, string(temp.Body.ClientCode[:]))                 //clientId
	dataToAdd = append(dataToAdd, AccountTypeMapping[uint(temp.Body.AccountType)]) //clientType
	//dataToAdd = append(dataToAdd, string(temp.Body.CPCCode[0:12]))                 //cp code
	// if !bytes.Equal(temp.Body.CPCCode[:], EmptyCPCCode) {                          //cp code confirmation
	// 	dataToAdd = append(dataToAdd, "Y")
	// } else {
	// 	dataToAdd = append(dataToAdd, "N")
	// }
	dataToAdd = append(dataToAdd, fmt.Sprintf("%v", temp.Body.OrderId))                          //orderId
	dataToAdd = append(dataToAdd, fmt.Sprintf("%v", temp.Body.TradeId))                          //tradeId
	dataToAdd = append(dataToAdd, fmt.Sprintf("%d", temp.Body.SecurityId))                       //scrip-code
	dataToAdd = append(dataToAdd, fmt.Sprintf("%.4f", float64(temp.Body.LastPrice)/100000000.0)) //Rate
	dataToAdd = append(dataToAdd, fmt.Sprintf("%d", temp.Body.LastQty))                          //qty
	dataToAdd = append(dataToAdd, BuySellMapping[uint(temp.Body.Side)])                          //buy-sell
	dataToAdd = append(dataToAdd, AOPOFlagMapping[temp.Body.MultiLegReportingType])              //ao/po flag
	dataToAdd = append(dataToAdd, fmt.Sprintf("%v", temp.Body.SenderLocationId))                 //locationId
	//dataToAdd = append(dataToAdd, time.Unix(0, int64(temp.Body.Activitytime)).Format("15:04:05"))      //"Order Time Stamp
	dataToAdd = append(dataToAdd, time.Unix(0, int64(temp.Body.TransactionTime)).Format("15:04:05"))   //time
	dataToAdd = append(dataToAdd, time.Unix(0, int64(temp.Body.TransactionTime)).Format("02-01-2006")) //date
	dataToAdd = append(dataToAdd, OrderTypeMapping[temp.Body.OrderType])                               //order-type
	//dataToAdd = append(dataToAdd, time.Unix(0, int64(temp.Body.Activitytime)).Format("15:04:05"))      //"TradeModification time

	// if db == nil {
	// 	return "", errors.New("database connection is nil")
	// }

	if err := utils.AddToTempTable(dataToAdd, db); err != nil {
		logger.Log().Error("Failed to add into temp table", zap.Error(err))
		return "", err
	}
	return data, nil

}

func getTraderId(data uint32) string {
	return fmt.Sprintf("%d", data%uint32(1000))
}
