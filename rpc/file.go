// @Datetime  : 2019/10/19 8:41 下午
// @Author    : psyduck
// @Purpose   :
// @TODO      : Pair programming
//
package rpc

import (
	"fmt"
	"google.golang.org/genproto/googleapis/bytestream"
	"io"
	"os"
)

var fileLoggerAddition = [2]string{"service", "file"}

func getFileDigest() {
	p.Panicw("Not implemented", fileLoggerAddition)
}

func checkFile() {
	p.Panicw("Not implemented", fileLoggerAddition)
}

func (s *Server) Write(server bytestream.ByteStream_WriteServer) error {
	p.Infow("Write method is called, start to receive file", fileLoggerAddition)
	// var fileWriteOffset int64 = 0
	var f *os.File
	for true {
		req, err := server.Recv()
		if err != nil {
			p.Errorw(fmt.Sprintf("server.Recv, %#v", err), fileLoggerAddition)
			return err
		}
		if req.FinishWrite == true {
			p.Infow("req.FinishWrite is true", fileLoggerAddition)
			f.Sync()
			f.Close()
			getFileDigest()
			checkFile()
			break
		}
		// p.Debugw("req.WriteOffset,",req.WriteOffset)
		f, err = os.OpenFile(req.ResourceName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			p.Fatalw(fmt.Sprintf("Open file %s for write failed, %#v", req.ResourceName, err), fileLoggerAddition)
			return err
		}
		//if fileWriteOffset == 0 {
		//	f, err = os.Create(req.ResourceName)
		//	if err != nil {
		//		log.Print("os.Create,",err)
		//		break
		//	}
		//} else {
		//	f, err = os.Open(req.ResourceName)
		//	if err != nil {
		//		log.Print("os.Open,",err)
		//		break
		//	}
		//}
		// n ,err := f.WriteAt(req.GetData(),fileWriteOffset)
		_, err = f.Write(req.GetData())
		//log.Print("before add fileWriteOffset,",fileWriteOffset,"n:",n,"reqData",len(req.GetData()))
		//fileWriteOffset += int64(n)
		//log.Print("after add fileWriteOffset,",fileWriteOffset,"n:",n)
	}
	return nil
}

func (s *Server) Read(in *bytestream.ReadRequest, server bytestream.ByteStream_ReadServer) error {
	p.Infow("Read method is called, start to send file", fileLoggerAddition)
	fileHandler, err := os.OpenFile(in.ResourceName, os.O_RDONLY, 888)
	if err != nil {
		p.Fatalw(fmt.Sprintf("os.OpenFile failed,%#v", err), fileLoggerAddition)
		return err
	}
	var fileReadOffset int64 = 0
	var fileBuffer [1024 * 1024]byte // 1M 的缓冲
	for true {
		n, err := fileHandler.Read(fileBuffer[:])
		if err == io.EOF {
			err = server.Send(&bytestream.ReadResponse{Data: []byte{}})
			if err != nil {
				p.Fatalw(fmt.Sprintf("file end, stop r.Send, %#v", err), fileLoggerAddition)
				return err
			}
			p.Infow(fmt.Sprintf("READ FILE TO END, n: %#v", n), fileLoggerAddition)
			break
		}
		if err != nil {
			p.Fatalw(fmt.Sprintf("fileHandler.Read, %#v", err), fileLoggerAddition)
			return err
		}
		fileReadOffset += int64(n)
		p.Debugw(fmt.Sprintf("fileReadOffset: %#v", fileReadOffset), fileLoggerAddition)
		err = server.Send(&bytestream.ReadResponse{Data: fileBuffer[:n]})
		if err != nil {
			p.Fatalw(fmt.Sprintf("server.Send, %#v", err), fileLoggerAddition)
			return err
		}
	}

	return nil
}
