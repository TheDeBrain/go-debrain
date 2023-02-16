package protocols

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/derain/core/rules"
	"io"
	"unsafe"
)

type ResultCollect struct {
	ProtocolType   uint8  `json:"protocol_type"` // Declare the protocol type , 1 byte
	ResultListSzie uint64 `json:"result_list_szie"`
	ResultList     []byte `json:"result_list"`
	EndFlag        []byte `json:"end_flag"`
}

type Result struct {
	AddrSize      uint16 `json:"addr_size"`
	Portsize      uint16 `json:"portsize"`
	FlagSize      uint32 `json:"flag_size"`
	DescribeSize  uint64 `json:"describe_size"`
	FileBlockSize uint64 `json:"fb_size"`
	Addr          []byte `json:"addr"`
	Port          []byte `json:"port"`
	Flag          []byte `json:"flag"`
	Describe      []byte `json:"describe"`
	FileBlock     []byte `json:"file_block"`
	EndFlag       []byte `json:"end_flag"`
}

// -------------------------- result collect handle start --------------------------

func RCNew(resArr []Result) (*ResultCollect, error) {
	rsl, err := RESEecoding(resArr)
	if err != nil {
		return nil, err
	}
	rc := ResultCollect{
		rules.RESULT_COLLECT_PROTOCOL,
		uint64(len(rsl)),
		rsl,
		[]byte(rules.FILE_BLCOK_END_FLAG),
	}
	return &rc, nil
}

func RCBuf(rc *ResultCollect) (*bytes.Buffer, error) {
	buff := bytes.NewBuffer([]byte{})
	// read in protocol type
	err := binary.Write(buff, binary.BigEndian, rc.ProtocolType)
	// read in result list size
	err = binary.Write(buff, binary.BigEndian, rc.ResultListSzie)
	// read in result list
	err = binary.Write(buff, binary.BigEndian, rc.ResultList)
	if err != nil {
		return nil, err
	}
	return buff, nil
}

func RCReader(r io.Reader) (*ResultCollect, error) {
	rc := new(ResultCollect)
	// ---------------------------- protocol head ----------------------------
	// protocol type
	protocolTypeBuf := make([]byte, int(unsafe.Sizeof(rc.ProtocolType)))
	_, err := r.Read(protocolTypeBuf)
	ptBuf := bytes.NewReader(protocolTypeBuf)
	binary.Read(ptBuf, binary.BigEndian, &rc.ProtocolType)
	if err != nil {
		return rc, err
	}
	// result list size
	resultListSizeBuf := make([]byte, int(unsafe.Sizeof(rc.ResultListSzie)))
	_, err = r.Read(resultListSizeBuf)
	rlBuf := bytes.NewReader(resultListSizeBuf)
	binary.Read(rlBuf, binary.BigEndian, &rc.ResultListSzie)
	if err != nil {
		return rc, err
	}
	// ---------------------------- protocol body ----------------------------
	// result list buf
	resListBuf := make([]byte, rc.ResultListSzie)
	_, err = r.Read(resListBuf)
	rc.ResultList = resListBuf
	if err != nil {
		return rc, err
	}
	// ---------------------------- protocol foot ----------------------------
	// end flag data
	endBuf := make([]byte, len([]byte(rules.FILE_BLCOK_END_FLAG)))
	_, err = r.Read(endBuf)
	rc.EndFlag = endBuf
	if err != nil {
		return rc, err
	}
	endFlag := string(rc.EndFlag)
	if endFlag == rules.FILE_BLCOK_END_FLAG {
		return rc, nil
	}
	return rc, errors.New("illegal agreement")
}

func RCWriter(w io.Writer, rc *ResultCollect) error {
	fbArr, err := RCBuf(rc)
	if err != nil {
		return err
	}
	w.Write(fbArr.Bytes())
	return nil
}

func RCEecoding(rc ResultCollect) ([]byte, error) {
	b, err := json.Marshal(rc)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func RCDecoding(rcArr []byte) (*ResultCollect, error) {
	rc := new(ResultCollect)
	err := json.Unmarshal(rcArr, rc)
	if err != nil {
		return nil, err
	}
	return rc, nil
}

// -------------------------- result collect handle end --------------------------

// -------------------------- result handle start --------------------------

func RESNew(addr string, port string, flag string, des string, fbSlice []*FileBlock) (*Result, error) {
	var allFbSlice []byte
	if fbSlice != nil {
		for _, fb := range fbSlice {
			fbuf, _ := FBBuf(fb)
			allFbSlice = append(allFbSlice, fbuf.Bytes()...)
		}
	}
	res := Result{
		uint16(len([]byte(addr))),
		uint16(len([]byte(port))),
		uint32(len([]byte(flag))),
		uint64(len([]byte(des))),
		uint64(len(allFbSlice)),
		[]byte(addr),
		[]byte(port),
		[]byte(flag),
		[]byte(des),
		allFbSlice,
		[]byte(rules.FILE_BLCOK_END_FLAG),
	}
	return &res, nil
}

func RESBuf(res *Result) (*bytes.Buffer, error) {
	buff := bytes.NewBuffer([]byte{})
	err := binary.Write(buff, binary.BigEndian, res.AddrSize)
	err = binary.Write(buff, binary.BigEndian, res.Portsize)
	err = binary.Write(buff, binary.BigEndian, res.FlagSize)
	err = binary.Write(buff, binary.BigEndian, res.DescribeSize)
	err = binary.Write(buff, binary.BigEndian, res.FileBlockSize)
	err = binary.Write(buff, binary.BigEndian, res.Addr)
	err = binary.Write(buff, binary.BigEndian, res.Port)
	err = binary.Write(buff, binary.BigEndian, res.Flag)
	err = binary.Write(buff, binary.BigEndian, res.Describe)
	err = binary.Write(buff, binary.BigEndian, res.FileBlock)
	err = binary.Write(buff, binary.BigEndian, []byte(rules.FILE_BLCOK_END_FLAG))
	if err != nil {
		return nil, err
	}
	return buff, nil
}

func RESReader(r io.Reader) (*Result, error) {
	rc := new(Result)
	// ----------------------------  head ----------------------------
	asBuf := make([]byte, int(unsafe.Sizeof(rc.AddrSize)))
	_, err := r.Read(asBuf)
	addrSizeBuf := bytes.NewReader(asBuf)
	binary.Read(addrSizeBuf, binary.BigEndian, &rc.AddrSize)
	if err != nil {
		return rc, err
	}
	psBuf := make([]byte, int(unsafe.Sizeof(rc.Portsize)))
	_, err = r.Read(psBuf)
	portSizeBuf := bytes.NewReader(psBuf)
	binary.Read(portSizeBuf, binary.BigEndian, &rc.Portsize)
	if err != nil {
		return rc, err
	}
	flagSizeBuf := make([]byte, int(unsafe.Sizeof(rc.FlagSize)))
	_, err = r.Read(flagSizeBuf)
	fsBuf := bytes.NewReader(flagSizeBuf)
	binary.Read(fsBuf, binary.BigEndian, &rc.FlagSize)
	if err != nil {
		return rc, err
	}
	describeSizeBuf := make([]byte, int(unsafe.Sizeof(rc.DescribeSize)))
	_, err = r.Read(describeSizeBuf)
	desSizeBuf := bytes.NewReader(describeSizeBuf)
	binary.Read(desSizeBuf, binary.BigEndian, &rc.DescribeSize)
	if err != nil {
		return rc, err
	}
	fbSizeBuf := make([]byte, int(unsafe.Sizeof(rc.FileBlockSize)))
	_, err = r.Read(fbSizeBuf)
	fBSizeBuf := bytes.NewReader(fbSizeBuf)
	binary.Read(fBSizeBuf, binary.BigEndian, &rc.FileBlockSize)
	if err != nil {
		return rc, err
	}
	// ----------------------------  body ----------------------------
	aBuf := make([]byte, rc.AddrSize)
	_, err = r.Read(aBuf)
	rc.Addr = aBuf
	if err != nil {
		return rc, err
	}
	pBuf := make([]byte, rc.Portsize)
	_, err = r.Read(pBuf)
	rc.Port = pBuf
	if err != nil {
		return rc, err
	}
	fBuf := make([]byte, rc.FlagSize)
	_, err = r.Read(fBuf)
	rc.Flag = fBuf
	if err != nil {
		return rc, err
	}
	dBuf := make([]byte, rc.DescribeSize)
	_, err = r.Read(dBuf)
	rc.Describe = dBuf
	if err != nil {
		return rc, err
	}
	fbBuf := make([]byte, rc.FileBlockSize)
	_, err = r.Read(fbBuf)
	rc.FileBlock = fbBuf
	if err != nil {
		return rc, err
	}
	// ----------------------------  foot ----------------------------
	// end flag data
	endBuf := make([]byte, len([]byte(rules.FILE_BLCOK_END_FLAG)))
	_, err = r.Read(endBuf)
	rc.EndFlag = endBuf
	if err != nil {
		return rc, err
	}
	endFlag := string(rc.EndFlag)
	if endFlag == rules.FILE_BLCOK_END_FLAG {
		return rc, nil
	}
	return rc, errors.New("illegal agreement")
}

func RESWriter(w io.Writer, res *Result) error {
	fbArr, err := RESBuf(res)
	if err != nil {
		return err
	}
	w.Write(fbArr.Bytes())
	return nil
}

func RESEecoding(res []Result) ([]byte, error) {
	b, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func RESDecoding(resArr []byte) ([]Result, error) {
	res := new([]Result)
	err := json.Unmarshal(resArr, res)
	if err != nil {
		return nil, err
	}
	return *res, nil
}

// -------------------------- result handle end --------------------------
