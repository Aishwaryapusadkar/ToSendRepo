package main

import (
	"bse-eti-stream/pkg/config"
	"bse-eti-stream/pkg/logger"
	"bse-eti-stream/pkg/messages"

	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var (
		environment string
	)
	ctx := context.Background()
	env := os.Args[1]
	if env != "uat" {
		environment = "prod"
	} else {
		environment = "uat"
	}
	fmt.Println("ENV:", environment)
	config.Load(environment)

	if err := start(ctx); err != nil {
		log.Fatal("Failed to start server, err:", err)
		os.Exit(1)
	}
}

func start(c context.Context) error {
	var (
		conn net.Conn
		//gConn net.Conn
	)
	cnf := config.GetConfig()
	logLevel, err := strconv.Atoi(cnf.GetString("log.Level"))
	if err != nil {
		log.Println("Invalid log config: ", err)
		return err
	}
	logger.LoggerInit(cnf.GetString("log.path"), zapcore.Level(logLevel))

	// utils.CreateWorkerPool(c, cnf.GetInt("workerpool.size"))
	// if err := utils.CreateTradeTable(c, open_position, BxtTrade); err != nil {
	// 	log.Println("Failed to create the table", err)
	// 	return err
	// }

	//tcp connection done through API gateway
	conn, err = net.Dial("tcp", fmt.Sprintf("%s:%s", config.GetConfig().GetString("gateway.host"), config.GetConfig().GetString("gateway.port")))
	if err != nil {
		logger.Log(c).Error("Failed to connect tcp connection", zap.Error(err))
		return err
	}
	defer conn.Close()

	// Log successful TCP connection
	//logger.Log(c).Info("Connected to TCP gateway successfully.")

	// Database connection details from configuration
	// dbHost := config.GetConfig().GetString("db.host")
	// dbPort := config.GetConfig().GetString("db.port")
	// dbUser := config.GetConfig().GetString("db.user")
	// dbPassword := config.GetConfig().GetString("db.password")
	// dbName := config.GetConfig().GetString("db.name")

	//Construct the PostgreSQL connection string
	//dbConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	//    dbConnStr := "host=localhost port=5432 user=postgres password=password dbname=open_position sslmode=disable"

	//     // Open the database connection
	//     db, err := sql.Open("postgres", dbConnStr)
	//     if err != nil {
	//         logger.Log(c).Error("Failed to connect to database", zap.Error(err))
	//         return err
	//     }
	//     defer db.Close()

	//     // Ping to check if the database connection is alive
	//     err = db.Ping()
	//     if err != nil {
	//         logger.Log(c).Error("Failed to ping database", zap.Error(err))
	//         return err
	//     }

	// Log successful database connection
	// logger.Log(c).Info("Connected to the database successfully.")

	// //API gateway connection
	// primaryConnStr, secondaryConnStr, err := messages.ConnectToGateway(c, conn)
	// if err != nil {
	// 	logger.Log(c).Error("Gateway connection failed", zap.Error(err))
	// 	return err
	// }
	// logger.Log(c).Info("Gateway connected", zap.String("primaryHost", primaryConnStr), zap.String("secondaryHost", secondaryConnStr))

	// gConn, err = net.Dial("tcp", fmt.Sprintf("%s", primaryConnStr))
	// if err != nil {
	// 	logger.Log(c).Error("Failed to connect tcp connection to primary-gateway looking for secondaryConn", zap.Error(err))
	// 	gConn, err = net.Dial("tcp", fmt.Sprintf("%s", secondaryConnStr))
	// 	if err != nil {
	// 		logger.Log(c).Error("Failed to connect tcp connection to secondary-gateway", zap.Error(err))
	// 		return err
	// 	} else {
	// 		logger.Log(c).Info("tpc session connection made on ", zap.String("host", secondaryConnStr))
	// 	}
	// } else {
	// 	logger.Log(c).Info("tpc session connection made on ", zap.String("host", primaryConnStr))
	// }
	// defer gConn.Close()

	// if err := messages.SessionLogin(c, gConn); err != nil {
	// 	logger.Log(c).Error("Session login failed", zap.Error(err))
	// 	return err
	// }

	if err := messages.UserLogin(c, conn); err != nil {
		logger.Log(c).Error("User login failed", zap.Error(err))
		return err
	}

	if err := messages.Subscribe(c, conn);
	err != nil {
		logger.Log(c).Error("Subscription failed", zap.Error(err))
		return err
	}

	if err := messages.StartFeedCapturing(c, conn); err != nil {
		logger.Log(c).Error("Feed Capturing failed", zap.Error(err))
		return err
	}
	return nil
}
