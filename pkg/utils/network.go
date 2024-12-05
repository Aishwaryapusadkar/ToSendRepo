package utils

import (
	"bse-eti-stream/pkg/logger"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
)

func Send(conn net.Conn, data interface{}, mType string) error {
	dataToSend, err := toLitteEndian(data, mType)
	logger.Log().Error("  sending the request")
	if err != nil {
		logger.Log().Error("Failed to convert the request to littleEndian format", zap.Error(err))
		return err
	}

	if len(dataToSend) > 0 {
		logger.Log().Error(" sending length")

		size, err := conn.Write(dataToSend)
		logger.Log().Error(string(size))
		if err != nil {
			logger.Log().Error("Failed to send the data", zap.Error(err))
			return err
		}
		logger.Log().Info("Sending data", zap.String("mType", mType), zap.Int("size", size), zap.Any("pkt-format", data), zap.String("hex-format", hex.EncodeToString(dataToSend)), zap.Any("byteFormat", dataToSend))
	}
	return nil
}

func Recv(conn net.Conn, dst interface{}, size int, mType string) ([]byte, error) {
	bytes := make([]byte, size)
	fmt.Println("Reading:", size)
	n, err := conn.Read(bytes)
	fmt.Printf("Received bytes from conn.Read:")
	fmt.Println("receiving receiving ", bytes)

	//fmt.Println("RECIEVING DATA")
	if err != nil {
		logger.Log().Error("Failed to read data", zap.String("mType", mType))
		return bytes, err
	}
	if n > 0 {
		// dst = bytes
		// fmt.Println("##", dst)

		// TODO: Populate the dst with the bytes, without converting the byte ordering

		// fmt.Println("Before converting to little endian: ", dst)
		err := fromLittleEndian(bytes, dst, mType)
		// fmt.Println("After converting to little endian: ", dst)
		if err == nil {
			// logger.Log().Info("data     ...........", zap.Any("data", dst))
		} else {
			logger.Log().Error("Failed to parse data from little endian", zap.String("mType", mType), zap.Error(err))

			return bytes, err
		}
	}
	fmt.Println("$$$$$$$$$bytes", string(bytes))

	fmt.Println("dst:", dst)
	time.Sleep(time.Second * 1)

	return bytes, nil
}

func toLitteEndian(data interface{}, mType string) ([]byte, error) {
	fmt.Println("--------------", data)
	fmt.Println(mType)
	var (
		buf = new(bytes.Buffer)
	)

	if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
		logger.Log().Error("failed to convert to littl-endian", zap.Error(err), zap.String("mType", mType))
		return nil, err
	}
	return buf.Bytes(), nil
}

func fromLittleEndian(dataToRead []byte, dst interface{}, mType string) error {
	buff := bytes.NewReader(dataToRead)
	// fmt.Println("-------", dataToRead)
	// fmt.Println(mType)
	// fmt.Println("_______", dst)
	// fmt.Println("*****", buff)

	if err := binary.Read(buff, binary.LittleEndian, dst); err != nil {
		logger.Log().Error("failed to convert into object from little-endian", zap.Error(err), zap.String("mType", mType))
		return err
	}
	// fmt.Println("After binary.read: ", dst)
	return nil
}

func IntToIPV4(data1 []uint8, data2 uint32) string {
	return fmt.Sprintf("%v.%v.%v.%v:%v", data1[3], data1[2], data1[1], data1[0], data2)
}

func TestRecv(bytes []byte, dst interface{}, size int, mType string) ([]byte, error) {
	if err := fromLittleEndian(bytes, dst, mType); err != nil {
		logger.Log().Error("Failed to parse data from little endian", zap.String("mType", mType), zap.Error(err))
		return bytes, err
	}
	return bytes, nil
}
