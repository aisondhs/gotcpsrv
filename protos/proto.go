package protos

func GetFuncName(msgId uint16) string {
	var actList map[uint16](string)
	actList = make(map[uint16](string), 100)
	actList[1] = "CSGetuserReq"
	actList[2] = "CSGetuserRsp"
	if funcname, ok := actList[msgId]; ok {
    	return funcname
    }
    return ""
}
