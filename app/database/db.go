package database

import (
	"database/sql"
	"fmt"
	"log"

	//_ "github.com/lib/pq"
	//"go.uber.org/zap"
)

// Database connection parameters
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "open_position"
)
var db *sql.DB

//func Init(){}

// InitializeDB initializes a connection to the PostgreSQL database and returns the database instance
func InitializeDB() (*sql.DB, error) {
	// Construct connection string
	dbConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Open the database connection
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Println("Failed to connect to database:", err)
		return nil, err
	}

	// Ping to check if the database connection is alive
	err = db.Ping()
	if err != nil {
		log.Println("Failed to ping database:", err)
		return nil, err
	}

	log.Println("Connected to the database successfully.")
	return db, nil
}

func GetDB() *sql.DB{
	if db == nil{
		log.Println("Warning : GetDB called but db is nil , make sure Init is called first")
	}
	return db
}

// BxtTrade represents a trade record in the bxt_trade table
// type BxtTrade struct {
//     BxtMmbrId        string  
//     BxtTrdrId        string  
//     BxtScripCd       string  
//     BxtScripId       string  
//     BxtRate          float64 
//     BxtQty           int32   
//     BxtTrdStts       string  
//     BxtClrngMmbrCd   string  
//     BxtTm            string  
//     BxtDt            string  
//     BxtClntId        string  
//     BxtOrdrId        string  
//     BxtOrdrTyp       string  
//     BxtFlw           string  
//     BxtTrdId         string  
//     BxtInstitutId    string  
//     BxtIsinCd        string  
//     BxtScripGrp      string  
//     BxtSttlmntNum    string  
//     BxtOrdrTimeStmp  string  
//     BxtAoPoFlg       string  
//     BxtLoctnId       string  
//     BxtTrdModfTm     string  
//     BxtSssnId        string  
//     BxtCpCd          string  
//     BxtCpConfmFlg    string  
//     BxtInsrtDt       string  
// }


// type Bxttrade struct {
//     BxtSecurityId           int64
//     BxtRelatedSecurityId    int64
//     BxtPrice                float64
//     BxtLastPrice            float64
//     BxtSideLastPrice        float64
//     BxtTransactionTime      int64
//     BxtOrderId              int64
//     BxtSenderLocationId     int64
//     BxtCIOrdId              int64
//     BxtMsgTag               int32
//     BxtTradeId              int64
//     BxtOrigTradeId          int64
//     BxtBusinessUnitId       int64
//     BxtSessionId            int64
//     BxtOwnerUserId          int64
//     BxtPartyIdClearingUnit  int64
//     BxtCumQty               int32
//     BxtLeavesQty            int32
//     BxtMarketSegmentId      int32
//     BxtRelatedSymbol        int32
//     BxtLastQty              int32
//     BxtSideLastQty          int32
//     BxtSideTradeId          int64
//     BxtMatchDate            int32
//     BxtMatch                int32
//     BxtStrategyLinkId       int32
//     BxtTotNumTradeReports   int32
//     BxtMultiLegReportingType int32
//     BxtTradeReportType      int32
//     BxtTransferReason       int32
//     BxtPartyIdBeneficiary   string
//     BxtPartyIdTakeupTradingFirm string
//     BxtPartyIdOrderOriginatingFirm string
//     BxtAccountType          int32
//     BxtAggressorSide        int32
//     BxtMatchType            int32
//     BxtMatchSubType         int32
//     BxtSide                 int32
//     BxtAggressorIndicator   int32
//     BxtTradingCapacity      int32
//     BxtAccount              string
//     BxtPositionEffect       string
//     BxtFreeText1            string
//     BxtFreeText2            string
//     BxtFreeText3            string
//     BxtOrderCategory        string
//     BxtOrderType            int32
//     BxtRelatedProductComplex int32
//     BxtOrderSide            int32
//     BxtPartyClearingOrganisation string
//     BxtPartyExecutingFirm   string
//     BxtPartyExecutingTrader string
//     BxtPartyClearingFirm    string
//     BxtFiller5              string
// }







