package utils

import (
	"bse-eti-stream/app/database"
	"bse-eti-stream/pkg/logger"
	"bse-eti-stream/pkg/models"
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"go.uber.org/zap"
)

// //////////////////**********************///////////////////////
// Database connection parameters
// const (
// 	host     = "localhost"
// 	port     = 5432
// 	user     = "postgres"
// 	password = "password"
// 	dbname   = "open_position"
// )

var db *sql.DB

//*************////
//func Init(){}
// func InitializeDB() (*sql.DB, error) {

// 	dbConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

// 	db, err := sql.Open("postgres", dbConnStr)
// 	if err != nil {
// 		fmt.Println("Failed to connect to database:", err)
// 		return nil, err
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		fmt.Println("Failed to ping database:", err)
// 		return nil, err
// 	}

// 	fmt.Println("Connected to the database successfully.")
// 	return db, nil
// }

// func GetDB() *sql.DB {
// 	if db == nil {
// 		fmt.Println("Warning : GetDB called but db is nil , make sure Init is called first")
// 	}
// 	return db
// }

///////***************////////////////////

var (
	wp *WorkerPool
)

type WorkerPool struct {
	context  context.Context
	Capacity int
	Pool     []chan interface{}
	wg       *sync.WaitGroup
	db       *sql.DB
}

type WorkerPoolInterface interface {
	EnQueue(context.Context, int, interface{})
	DeQueue(context.Context, int, func(interface{}) (interface{}, error))
	Run(context.Context, func(interface{}) (interface{}, error))
	Size() int
	Wait()
}

func CreateWorkerPool(c context.Context, size int) {
	fmt.Println("inside the workerpools")
	if size == 0 {
		size = 1
	}
	wp = &WorkerPool{
		Capacity: size,
	}
	wp.Pool = make([]chan interface{}, size)
	for i := 0; i < size; i++ {
		wp.Pool[i] = make(chan interface{}, size)
	}
	wp.context = c
	wp.wg = &sync.WaitGroup{}

	dbConn, err := database.InitializeDB()
	if err != nil {
		logger.Log(c).Error("Unable to connect to DB")
	} else {
		wp.db = dbConn
	}

	logger.Log(c).Info("worker pool is created of size", zap.Int("", size))
}

func GetWorkerPool() WorkerPoolInterface {
	fmt.Println(wp, "wp ia")
	return wp
}
func (w *WorkerPool) Wait() {
	w.wg.Wait()
}
func (w *WorkerPool) Size() int {
	return w.Capacity
}

func (w *WorkerPool) Run(c context.Context, task func(interface{}) (interface{}, error)) {
	w.wg.Add(w.Capacity)
	for i := 0; i < wp.Capacity; i++ {
		logger.Log(c).Info("worker started", zap.Int("index", i))
		go wp.DeQueue(c, i, task)
		// wp.DeQueue(c, i, task)
	}
}

func hasField(m map[string]interface{}, field string) bool {
	_, exists := m[field]
	return exists
}

// var tempTable1 = struct {
// 	sync.Mutex
// 	data []models.Bxttrade
// }{}

func ginContext(ctx *gin.Context) *gin.Context {
	return ctx
}

func (w *WorkerPool) EnQueue(c context.Context, index int, data interface{}) {
	fmt.Println("Before adding data into the chan")
	select {
	case w.Pool[index] <- data:
		fmt.Println("Data sent into the channel")

	default:
		// fmt.Println("Data not received yet!")
	}

	logger.Log(c).Debug("Added to Queue", zap.Any("indexdAt", index))
	logger.Log(c).Debug("fetched from Queue", zap.Any("id", index), zap.Any("data", data))

	fmt.Println("Processing data to insert into the database", data)

	data1 := data.(models.TradeObject)

	fmt.Println("data1 :", data1)

	//data2 := data.(models.TradeHeader)

	//fmt.Println("data2 :", data2)
	if wp.db == nil {
		fmt.Println("database connection is null")
	}

	tradeHeader := models.TradeHeader{
		SecurityId:                 data1.Body.SecurityId,
		RelatedSecurityId:          data1.Body.RelatedSecurityId,
		Price:                      data1.Body.Price,
		LastPrice:                  data1.Body.LastPrice,
		SideLastPrice:              data1.Body.SideLastPrice,
		TransactionTime:            data1.Body.TransactionTime,
		OrderId:                    data1.Body.OrderId,
		SenderLocationId:           data1.Body.SenderLocationId,
		CIOrdId:                    data1.Body.SenderLocationId,
		MsgTag:                     data1.Body.MsgTag,
		TradeId:                    data1.Body.TradeId,
		OrigTradeId:                data1.Body.OrigTradeId,
		BusinessUnitId:             data1.Body.BusinessUnitId,
		SessionId:                  data1.Body.SessionId,
		OwnerUserId:                data1.Body.OwnerUserId,
		PartyIdClearingUnit:        data1.Body.PartyIdClearingUnit,
		CumQty:                     data1.Body.CumQty,
		LeavesQty:                  data1.Body.LeavesQty,
		MarketSegmentId:            data1.Body.MarketSegmentId,
		RelatedSymbol:              data1.Body.RelatedSymbol,
		LastQty:                    data1.Body.LastQty,
		SideLastQty:                data1.Body.SideLastQty,
		SideTradeId:                data1.Body.SideTradeId,
		MatchDate:                  data1.Body.MatchDate,
		TradeMatch:                 data1.Body.TradeMatch,
		StrategyLinkId:             data1.Body.StrategyLinkId,
		TotNumTradeReports:         data1.Body.TotNumTradeReports,
		MultiLegReportingType:      data1.Body.MultiLegReportingType,
		TradeReportType:            data1.Body.TradeReportType,
		TrasnferReason:             data1.Body.TrasnferReason,
		PartyIdBeneficiery:         data1.Body.PartyIdBeneficiery,
		PartyIdTakeupTradingfirm:   data1.Body.PartyIdTakeupTradingfirm,
		PartyIdOrderOrignatingFirm: data1.Body.PartyIdOrderOrignatingFirm,
		AccountType:                data1.Body.AccountType,
		AggresorSide:               data1.Body.AggresorSide,
		MatchType:                  data1.Body.MatchType,
		MatchSubType:               data1.Body.MatchSubType,
		Side:                       data1.Body.Side,
		AggresorIndicator:          data1.Body.AggresorIndicator,
		TradingCapacity:            data1.Body.TradingCapacity,
		Account:                    data1.Body.Account,
		PositionEffect:             data1.Body.PositionEffect,
		OrderCategory:              data1.Body.OrderCategory,
		OrderType:                  data1.Body.OrderType,
		RelatedproductComplex:      data1.Body.RelatedproductComplex,
		OrderSide:                  data1.Body.OrderSide,
		PartyClearingOrganisation:  data1.Body.PartyClearingOrganisation,
		PartyExecutingFirm:         data1.Body.PartyExecutingFirm,
		PartyExecutingTrader:       data1.Body.PartyExecutingTrader,
		PartyClearingFirm:          data1.Body.PartyClearingFirm,
	}

	// // tempTable.data = append(tempTable.data, database.BxtTrade(data1)) // Replace `BxtTradeData` with your actual struct type
	// switch v := data.(type) {
	// case models.Bxttrade:
	// 	fmt.Println(v.BxtTradeId)
	// 	// 	switch z:=v.RHeader{
	// 	// 		cA
	// 	// 	}
	// }
	// ctx := ginContext()
	// var request [][]map[string]interface{}
	// if castedData, ok := interface{}(data).([][]map[string]interface{}); ok {
	// 	request = castedData
	// 	fmt.Println("Request:", request)
	// } else {
	// 	fmt.Println("Failed to cast data to [][]map[string]interface{}")
	// }
	// if err := ctx.BindJSON(&request); err != nil {
	// 	fmt.Printf("Invalid request")
	// }
	fmt.Println("I AM HERE======================", tradeHeader.SecurityId)
	fmt.Println("I AM HERE======================", tradeHeader.Price)

	fmt.Println("PartyIdBeneficiery ", string(tradeHeader.PartyIdBeneficiery[:]))

	fmt.Println("I AM everywhere======================", data1.Body.SecurityId,
		data1.Body.RelatedSecurityId,
		data1.Body.Price,
		data1.Body.LastPrice,
		data1.Body.SideLastPrice,
		data1.Body.TransactionTime,
		data1.Body.OrderId,
		data1.Body.SenderLocationId,
		data1.Body.CIOrdId,
		data1.Body.MsgTag,
		data1.Body.TradeId,
		data1.Body.OrigTradeId,
		data1.Body.BusinessUnitId,
		data1.Body.SessionId,
		data1.Body.OwnerUserId,
		data1.Body.PartyIdClearingUnit,
		data1.Body.CumQty,
		data1.Body.LeavesQty,
		data1.Body.MarketSegmentId,
		data1.Body.RelatedSymbol,
		data1.Body.LastQty,
		data1.Body.SideLastQty,
		data1.Body.SideTradeId,
		data1.Body.MatchDate,
		data1.Body.TradeMatch,
		data1.Body.StrategyLinkId,
		data1.Body.TotNumTradeReports,
		data1.Body.MultiLegReportingType,
		data1.Body.TradeReportType,
		data1.Body.TrasnferReason,
		data1.Body.PartyIdBeneficiery,
		data1.Body.PartyIdTakeupTradingfirm,
		data1.Body.PartyIdOrderOrignatingFirm,
		data1.Body.AccountType,
		data1.Body.AggresorSide,
		data1.Body.MatchType,
		data1.Body.MatchSubType,
		data1.Body.Side,
		data1.Body.AggresorIndicator,
		data1.Body.TradingCapacity,
		data1.Body.Account,
		data1.Body.PositionEffect,
		data1.Body.FreeText1,
		data1.Body.FreeText2,
		data1.Body.FreeText3,
		data1.Body.OrderCategory,
		data1.Body.OrderType,
		data1.Body.RelatedproductComplex,
		data1.Body.OrderSide,
		data1.Body.PartyClearingOrganisation,
		data1.Body.PartyExecutingFirm,
		data1.Body.PartyExecutingTrader,
		data1.Body.PartyClearingFirm,
		data1.Body.Filler5,
	)

	query := `INSERT INTO bxttrade (bxtsecurityid,bxtrelatedsecurityid, bxtprice, bxtlastprice, bxtsidelastprice,
	bxttransactiontime, bxtorderid, bxtsenderlocationid, bxtciordid, bxtmsgtag,
	bxttradeid, bxtorigtradeid, bxtbusinessunitid, bxtsessionid, bxtowneruserid,
	bxtpartyidclearingunit, bxtcumqty, bxtleavesqty, bxtmarketsegmentid, bxtrelatedsymbol,
	bxtlastqty, bxtsidelastqty, bxtsidetradeid, bxtmatchdate, bxttradematch,
	bxtstrategylinkid, bxttotnumtradereports, bxtmultilegreportingtype,
	bxttradereporttype, bxttransferreason, bxtpartyidbeneficiary, bxtpartyidtakeuptradingfirm,
	bxtpartyidorderoriginatingfirm, bxtaccounttype, bxtaggressorside, bxtmatchtype,
	bxtmatchsubtype, bxtside, bxtaggressorindicator, bxttradingcapacity, bxtaccount,
	bxtpositioneffect, bxtfreetext1, bxtfreetext2, bxtfreetext3, bxtordercategory,
	bxtordertype, bxtrelatedproductcomplex, bxtorderside, bxtpartyclearingorganisation,
	bxtpartyexecutingfirm, bxtpartyexecutingtrader, bxtpartyclearingfirm, bxtfiller5) 
	VALUES( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
	$21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39,
	$40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54)`

	_, err := wp.db.Exec(query,
		data1.Body.SecurityId,
		data1.Body.RelatedSecurityId,
		data1.Body.Price,
		data1.Body.LastPrice,
		data1.Body.SideLastPrice,
		data1.Body.TransactionTime,
		data1.Body.OrderId,
		data1.Body.SenderLocationId,
		data1.Body.CIOrdId,
		data1.Body.MsgTag,
		data1.Body.TradeId,
		data1.Body.OrigTradeId,
		data1.Body.BusinessUnitId,
		data1.Body.SessionId,
		data1.Body.OwnerUserId,
		data1.Body.PartyIdClearingUnit,
		data1.Body.CumQty,
		data1.Body.LeavesQty,
		data1.Body.MarketSegmentId,
		data1.Body.RelatedSymbol,
		data1.Body.LastQty,
		data1.Body.SideLastQty,
		data1.Body.SideTradeId,
		data1.Body.MatchDate,
		data1.Body.TradeMatch,
		data1.Body.StrategyLinkId,
		data1.Body.TotNumTradeReports,
		data1.Body.MultiLegReportingType,
		data1.Body.TradeReportType,
		data1.Body.TrasnferReason,
		string(data1.Body.PartyIdBeneficiery[:]),
		string(data1.Body.PartyIdTakeupTradingfirm[:]),
		string(data1.Body.PartyIdOrderOrignatingFirm[:]),
		data1.Body.AccountType,
		data1.Body.AggresorSide,
		data1.Body.MatchType,
		data1.Body.MatchSubType,
		data1.Body.Side,
		data1.Body.AggresorIndicator,
		data1.Body.TradingCapacity,
		string(data1.Body.Account[:]),
		string(data1.Body.PositionEffect[:]),
		string(data1.Body.FreeText1[:]),
		string(data1.Body.FreeText2[:]),
		string(data1.Body.FreeText3[:]),
		string(data1.Body.OrderCategory[:]),
		data1.Body.OrderType,
		data1.Body.RelatedproductComplex,
		data1.Body.OrderSide,
		string(data1.Body.PartyClearingOrganisation[:]),
		string(data1.Body.PartyExecutingFirm[:]),
		string(data1.Body.PartyExecutingTrader[:]),
		string(data1.Body.PartyClearingFirm[:]),
		string(data1.Body.Filler5[:]),
	)
	if err != nil {
		fmt.Println("Failed to insert data into bxttrade:", err)
		return
	}
	fmt.Println("Data successfully added to bxttrade table")

	fmt.Println("Data successfully added to bxt_trade table")
	//return nil
}

func (w *WorkerPool) DeQueue(c context.Context, index int, task func(interface{}) (interface{}, error)) {
	fmt.Println("inside a deque")
	defer w.wg.Done()
	data := <-w.Pool[index]
	fmt.Println("after fetching data from the channel")
	task(data)
	response, err := task(data)
	logger.Log(c).Debug("fetched from Queue", zap.Any("id", index), zap.Any("data", data), zap.Any("response", response), zap.Error(err))
}
