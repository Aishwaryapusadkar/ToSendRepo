package utils

import (
	//"bse-eti-stream/pkg/config"
	"bse-eti-stream/pkg/logger"
	"bse-eti-stream/app/database"

	"context"
	//"encoding/csv"
	"fmt"
	"os"
	"time"
	"sync"
	"log"

	"go.uber.org/zap"

	"database/sql"
	
)

var (
	//cwo          *csv.Writer
	uwo          *os.File
	userLogFiles = make(map[string]bool)
)

// func CreateCsvWriter(c context.Context, headers []string) error {
// 	logger.Log(c).Debug("STARTED")
// 	defer logger.Log(c).Debug("ENDED")

// 	if cwo == nil {
// 		dir := time.Now().Format("02-01-2006")
// 		err := os.Mkdir(dir, 0777)
// 		if err != nil {
// 			if os.IsExist(err) {
// 				logger.Log(c).Info("directory already exists", zap.String("dir-name", time.Now().Format("02-01-2006")))
// 			} else {
// 				logger.Log(c).Error("Failed to create directory", zap.Error(err))
// 				return err
// 			}
// 		}
// 		filename := fmt.Sprintf("%v/%s", dir, config.GetConfig().GetString("file.data-file"))
// 		fp, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
// 		if err != nil {
// 			logger.Log().Error("Failed to create/open file", zap.Error(err), zap.String("filename", filename))
// 			return err
// 		}
// 		logger.Log(c).Info("Filed opened", zap.String("filename", filename))
// 		cwo = csv.NewWriter(fp)
// 		if err := cwo.Write(headers); err != nil {
// 			logger.Log(c).Error("Failed to write csv header ", zap.Error(err), zap.Any("data", headers), zap.String("Filename", filename))
// 			return err
// 		}
// 		cwo.Flush()

// 		logger.Log().Info("writer object created & headers created")
// 	}
// 	return nil
// }

// func AddToCsv(data []string) error {
// 	logger.Log().Debug("STARTED")
// 	defer logger.Log().Debug("ENDED")
// 	if err := cwo.Write(data); err != nil {
// 		logger.Log().Error("Failed to write to csv", zap.Error(err), zap.Any("data", data))
// 		return err
// 	}
// 	return nil
// }

// func CSVFlush() error {
// 	if cwo != nil {
// 		cwo.Flush()
// 		if err := cwo.Error(); err != nil {
// 			logger.Log().Error("Failed to flush", zap.Error(err))
// 			return err
// 		}
// 	}
// 	return nil
// }

// In-memory structure to store data temporarily
// var tempTable = struct {
// 	sync.Mutex
// 	data []interface{}
// }{}

// CreateTradeTable creates a table in the database using the provided structure
// func CreateTradeTable(c context.Context, db *sql.DB, tableName string) error {
// 	log.Println("STARTED")
// 	defer log.Println("ENDED")

// 	// Define the SQL query to create the table
// 	createTableSQL := fmt.Sprintf(`
// 		CREATE TABLE IF NOT EXISTS %s (
// 			BXT_TRDR_ID          VARCHAR(50),
// 			BXT_SCRIP_CD         VARCHAR(50),
// 			BXT_SCRIP_ID         VARCHAR(50),
// 			BXT_RATE             DECIMAL(18, 2),
// 			BXT_QTY              INTEGER,
// 			BXT_TRD_STTS         VARCHAR(20),
// 			BXT_CLRNG_MBR_CD     VARCHAR(50),
// 			BXT_TM               VARCHAR(20),
// 			BXT_DT               DATE,
// 			BXT_CLNT_ID          VARCHAR(50),
// 			BXT_ORDR_ID          VARCHAR(50),
// 			BXT_ORDR_TYP         VARCHAR(20),
// 			BXT_FLW              VARCHAR(20),
// 			BXT_TRD_ID           VARCHAR(50),
// 			BXT_INSTITUT_ID      VARCHAR(50),
// 			BXT_ISIN_CD          VARCHAR(50),
// 			BXT_SCRIP_GRP        VARCHAR(50),
// 			BXT_STTLMNT_NUM      INTEGER,
// 			BXT_ORDR_TIME_STMP   TIMESTAMP,
// 			BXT_AO_PO_FLG        VARCHAR(10),
// 			BXT_LOCTN_ID         VARCHAR(50),
// 			BXT_TRD_MODF_TM      TIMESTAMP,
// 			BXT_SSSN_ID          VARCHAR(50),
// 			BXT_CP_CD            VARCHAR(50),
// 			BXT_CP_CONFM_FLG     VARCHAR(10),
// 			BXT_INSRT_DT         TIMESTAMP DEFAULT SYSDATE
// 		);`, tableName)

// 	// Execute the SQL statement to create the table
// 	_, err := db.ExecContext(c, createTableSQL)
// 	if err != nil {
// 		log.Println("Failed to create table:", err)
// 		return err
// 	}

// 	log.Println("Table created successfully with name:", tableName)
// 	return nil
// }

// AddToTempTable adds data to the temporary in-memory table
// 

// Assuming tempTable is a struct holding a slice and a mutex for safe concurrent access
var tempTable = struct {
	sync.Mutex
	data []database.BxtTrade
}{}

// AddToTempTable accepts data, converts it to the BxtTrade struct, and inserts it into the bxt_trade table.
func AddToTempTable(data interface{}, db *sql.DB) error {
	tempTable.Lock()
	defer tempTable.Unlock()

	// Cast data to BxtTrade struct
	bxtData, ok := data.(database.BxtTrade)
	if !ok {
		log.Println("Invalid data type. Expected BxtTrade struct.")
		return nil
	}

	// Append the data to the slice (this part is optional if you just want to insert directly)
	tempTable.data = append(tempTable.data, bxtData)

	// Insert data into the bxt_trade table
	query := `
		INSERT INTO bxt_trade (
			BXT_MMBR_ID, BXT_TRDR_ID, BXT_SCRIP_CD, BXT_SCRIP_ID, BXT_RATE, BXT_QTY, BXT_TRD_STTS,
			BXT_CLRNG_MBR_CD, BXT_TM, BXT_DT, BXT_CLNT_ID, BXT_ORDR_ID, BXT_ORDR_TYP, BXT_FLW, BXT_TRD_ID,
			BXT_INSTITUT_ID, BXT_ISIN_CD, BXT_SCRIP_GRP, BXT_STTLMNT_NUM, BXT_ORDR_TIME_STMP, BXT_AO_PO_FLG,
			BXT_LOCTN_ID, BXT_TRD_MODF_TM, BXT_SSSN_ID, BXT_CP_CD, BXT_CP_CONFM_FLG, BXT_INSRT_DT
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, CURRENT_TIMESTAMP
		)`

	// Execute the SQL insert statement with the BxtTrade data
	_, err := db.Exec(query,
		bxtData.BxtMmbrId, bxtData.BxtTrdrId, bxtData.BxtScripCd, bxtData.BxtScripId, bxtData.BxtRate,
		bxtData.BxtQty, bxtData.BxtTrdStts, bxtData.BxtClrngMmbrCd, bxtData.BxtTm, bxtData.BxtDt,
		bxtData.BxtClntId, bxtData.BxtOrdrId, bxtData.BxtOrdrTyp, bxtData.BxtFlw, bxtData.BxtTrdId,
		bxtData.BxtInstitutId, bxtData.BxtIsinCd, bxtData.BxtScripGrp, bxtData.BxtSttlmntNum,
		bxtData.BxtOrdrTimeStmp, bxtData.BxtAoPoFlg, bxtData.BxtLoctnId, bxtData.BxtTrdModfTm,
		bxtData.BxtSssnId, bxtData.BxtCpCd, bxtData.BxtCpConfmFlg,
	)

	if err != nil {
		log.Println("Failed to insert data into bxt_trade:", err)
		return err
	}

	log.Println("Data successfully added to bxt_trade table")
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
