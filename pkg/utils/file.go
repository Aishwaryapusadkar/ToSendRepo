package utils

import (
	//"bse-eti-stream/app/database"
	"bse-eti-stream/pkg/logger"
	"bse-eti-stream/pkg/models"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"

	"database/sql"
)

var (
	uwo          *os.File
	userLogFiles = make(map[string]bool)
)

// Assuming tempTable is a struct holding a slice and a mutex for safe concurrent access
var tempTable = struct {
	sync.Mutex
	data []models.Bxttrade
}{}

// AddToTempTable accepts data, converts it to the BxtTrade struct, and inserts it into the bxt_trade table.
func AddToTempTable(data interface{}, db *sql.DB) error {
	tempTable.Lock()
	defer tempTable.Unlock()

	// Cast data to BxtTrade struct
	bxtData, ok := data.(models.Bxttrade)
	if !ok {
		log.Println("Invalid data type. Expected BxtTrade struct.")
		return nil
	}

	// Append the data to the slice (this part is optional if you just want to insert directly)
	tempTable.data = append(tempTable.data, bxtData)



	// query := `
    // INSERT INTO Bxttrade (
    //     bxtsecurityid, bxtrelatedsecurityid, bxtprice, bxtlastprice, bxtsidelastprice,
    // 	bxttransactiontime, bxtorderid, bxtsenderlocationid, bxtciordid, bxtmsgtag,
    // 	bxttradeid, bxtorigtradeid, bxtbusinessunitid, bxtsessionid, bxtowneruserid,
    // 	bxtpartyidclearingunit, bxtcumqty, bxtleavesqty, bxtmarketsegmentid, bxtrelatedsymbol,
    // 	bxtlastqty, bxtsidelastqty, bxtsidetradeid, bxtmatchdate, bxttradematch, 
    // 	bxtstrategylinkid, bxttotnumtradereports, bxtmultilegreportingtype, 
    // 	bxttradereporttype, bxttransferreason, bxtpartyidbeneficiary, bxtpartyidtakeuptradingfirm,
    // 	bxtpartyidorderoriginatingfirm, bxtaccounttype, bxtaggressorside, bxtmatchtype,
    // 	bxtmatchsubtype, bxtside, bxtaggressorindicator, bxttradingcapacity, bxtaccount,
    // 	bxtpositioneffect, bxtfreetext1, bxtfreetext2, bxtfreetext3, bxtordercategory,
    // 	bxtordertype, bxtrelatedproductcomplex, bxtorderside, bxtpartyclearingorganisation,
    // 	bxtpartyexecutingfirm, bxtpartyexecutingtrader, bxtpartyclearingfirm, bxtfiller5,
    // 	bxtinsrtdt
    // ) VALUES (
    //     $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
    //     $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, 
    //     $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, CURRENT_TIMESTAMP
    // )`

	// _, err := db.Exec(query,
	// 	bxtData.BxtSecurityId, 
	// 	bxtData.BxtRelatedSecurityId,
	// 	bxtData.BxtPrice,
	// 	bxtData.BxtLastPrice,
	// 	bxtData.BxtSideLastPrice,
	// 	bxtData.BxtTransactionTime,
	// 	bxtData.BxtOrderId,
	// 	bxtData.BxtSenderLocationId,
	// 	bxtData.BxtCIOrdId,
	// 	bxtData.BxtMsgTag,
	// 	bxtData.BxtTradeId,
	// 	bxtData.BxtOrigTradeId,
	// 	bxtData.BxtBusinessUnitId,
	// 	bxtData.BxtSessionId,
	// 	bxtData.BxtOwnerUserId,
	// 	bxtData.BxtPartyIdClearingUnit,
	// 	bxtData.BxtCumQty,
	// 	bxtData.BxtLeavesQty,
	// 	bxtData.BxtMarketSegmentId,
	// 	bxtData.BxtRelatedSymbol,
	// 	bxtData.BxtLastQty,
	// 	bxtData.BxtSideLastQty,
	// 	bxtData.BxtSideTradeId,
	// 	bxtData.BxtMatchDate,
	// 	bxtData.BxtMatch,
	// 	bxtData.BxtStrategyLinkId,
	// 	bxtData.BxtTotNumTradeReports,
	// 	bxtData.BxtMultiLegReportingType,
	// 	bxtData.BxtTradeReportType,
	// 	bxtData.BxtTransferReason,
	// 	bxtData.BxtPartyIdBeneficiary,
	// 	bxtData.BxtPartyIdTakeupTradingFirm,
	// 	bxtData.BxtPartyIdOrderOriginatingFirm,
	// 	bxtData.BxtAccountType,
	// 	bxtData.BxtAggressorSide,
	// 	bxtData.BxtMatchType,
	// 	bxtData.BxtMatchSubType,
	// 	bxtData.BxtSide,
	// 	bxtData.BxtAggressorIndicator,
	// 	bxtData.BxtTradingCapacity,
	// 	bxtData.BxtAccount,
	// 	bxtData.BxtPositionEffect,
	// 	bxtData.BxtFreeText1,
	// 	bxtData.BxtFreeText2,
	// 	bxtData.BxtFreeText3,
	// 	bxtData.BxtOrderCategory,
	// 	bxtData.BxtOrderType,
	// 	bxtData.BxtRelatedProductComplex,
	// 	bxtData.BxtOrderSide,
	// 	bxtData.BxtPartyClearingOrganisation,
	// 	bxtData.BxtPartyExecutingFirm,
	// 	bxtData.BxtPartyExecutingTrader,
	// 	bxtData.BxtPartyClearingFirm,
	// 	bxtData.BxtFiller5,
	// )
	

	// if err != nil {
	// 	log.Println("Failed to insert data into BxtTrade:", err)
	// 	return err
	// }

	// log.Println("Data successfully added to BxtTrade table")
	return nil
}

func AddUserLog(c context.Context, userId string, partitionId uint16, seqNo uint64) error {
	var (
		err error
	)
	filename := fmt.Sprintf("%v/%s_%d.log", time.Now().Format("02-01-2006"), userId, partitionId)
	if !userLogFiles[filename] {
		uwo, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			logger.Log().Error("Failed to create/open user-log file", zap.Error(err), zap.String("filename", filename))
			return err
		}
		userLogFiles[filename] = true
	} else {
		if _, err = uwo.WriteString(fmt.Sprintf("%v\n", seqNo)); err != nil {
			logger.Log(c).Error("Failed to append in user-log file", zap.Error(err), zap.String("filename", filename), zap.Any("seqno", seqNo))
			return err
		} else {
			logger.Log(c).Debug("UserLog Added", zap.String("file", filename), zap.Any("seq", seqNo))
		}
	}
	return nil
}
