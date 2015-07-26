package service

import proto "github.com/golang/protobuf/proto"
import "github.com/aisondhs/gotcpsrv/protos"


func CSGetuserReq(reqBytes []byte)([]byte,error) {
	reqObj := new(protos.CSGetuserReq)
	err := proto.Unmarshal(reqBytes, reqObj)
	if err != nil {
		return nil,err
	}
	uid := reqObj.GetUid()

	rspObj := new(protos.CSGetuserRsp)
	rspObj.Uid = proto.Int32(uid)
	rspObj.Name = proto.String("aison")
	rspObj.Age = proto.Int32(30)
	rspObj.City = proto.String("shenzhen")
	buffer, err := proto.Marshal(rspObj)
	if err != nil {
		return nil,err
	}
	return buffer,nil
}