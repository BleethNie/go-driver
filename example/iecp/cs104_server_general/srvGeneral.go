package main

import (
	"fmt"
	"github.com/BleethNie/go-driver/iecp/asdu"
	"github.com/BleethNie/go-driver/iecp/cs104"
	"log"
	"time"
)

func main() {
	inits()
	createSrv(":2404")

}

type mysrv struct{}

func createSrv(port string) {
	srv := cs104.NewServer(&mysrv{})
	srv.SetOnConnectionHandler(func(c asdu.Connect) {
		log.Println("on connect")
	})
	srv.SetConnectionLostHandler(func(c asdu.Connect) {
		log.Println("connect lost")
	})
	srv.LogMode(true)
	srv.ListenAndServer(port)
}

func inits() {
	InitRandomData()
}

// 总召处理
func (sf *mysrv) InterrogationHandler(c asdu.Connect, asduPack *asdu.ASDU, qoi asdu.QualifierOfInterrogation) error {
	log.Println("qoi:", qoi)
	fmt.Println("call InterrogationHandler()")

	var yxNum int = 0
	var ycNum int = 0

	//总召激活确认
	asduPack.SendReplyMirror(c, asdu.ActivationCon)

	var asduTmp *asdu.ASDU = asdu.NewEmptyASDU(asdu.ParamsWide)
	var asduYcTmp *asdu.ASDU = asdu.NewEmptyASDU(asdu.ParamsWide)

	//遥信总召
	asduTmp.Type = asdu.M_SP_NA_1
	asduTmp.Variable.IsSequence = false
	infoObjLength, _ := asdu.GetInfoObjSize(asdu.M_SP_NA_1)
	elementLength := infoObjLength + asdu.ParamsWide.InfoObjAddrSize
	for k, v := range AppMap["yx"] {
		asduTmp.AppendInfoObjAddr(asdu.InfoObjAddr(k + 1))

		asduTmp.AppendBytes(byte(YXMap[v]))
		yxNum = yxNum + 1

		//判断是否超过数据长度，超过则发送数据
		if (asdu.ASDUSizeMax-asduTmp.IdentifierSize()-elementLength*yxNum) <= elementLength || k == (len(AppMap["yx"])-1) {
			asduTmp.SetVariableNumber(yxNum)
			err := asduTmp.SendICReply(c, asdu.InterrogatedByStation, asduPack.CommonAddr)
			asduTmp.CleanInfoObj()
			yxNum = 0
			if err != nil {
				log.Println("falied")
			} else {
				log.Println("success")
			}
		}
	}

	//遥测总召
	asduYcTmp.Type = asdu.M_ME_NC_1
	asduYcTmp.Variable.IsSequence = false
	infoObjLength, _ = asdu.GetInfoObjSize(asdu.M_ME_NC_1)
	elementLength = infoObjLength + asdu.ParamsWide.InfoObjAddrSize
	for k, v := range AppMap["yc"] {
		asduYcTmp.AppendInfoObjAddr(asdu.InfoObjAddr(YcAddress + k))

		asduYcTmp.AppendFloat32(YCMap[v])
		//品质信息设置默认值
		asduYcTmp.AppendBytes(byte(0x3F))

		ycNum = ycNum + 1

		//判断是否超过数据长度，超过则发送数据
		if (asdu.ASDUSizeMax-asduYcTmp.IdentifierSize()-elementLength*ycNum) <= elementLength || k == (len(AppMap["yc"])-1) {
			asduYcTmp.SetVariableNumber(ycNum)
			err := asduYcTmp.SendICReply(c, asdu.InterrogatedByStation, asduPack.CommonAddr)
			asduYcTmp.CleanInfoObj()
			ycNum = 0
			if err != nil {
				log.Println("falied")
			} else {
				log.Println("success")
			}
		}
	}

	//数据发送完成
	asduPack.SendReplyMirror(c, asdu.ActivationTerm)
	return nil
}
func (sf *mysrv) CounterInterrogationHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierCountCall) error {
	fmt.Println("call CounterInterrogationHandler()")
	return nil
}
func (sf *mysrv) ReadHandler(c asdu.Connect, asduPack *asdu.ASDU, addr asdu.InfoObjAddr) error {
	fmt.Println("call ReadHandler()")
	//发送单点数据
	asduPack.Type = asdu.M_SP_NA_1
	asduPack.AppendBytes(byte(YXMap[int(addr)]))
	asduPack.SendICReply(c, asdu.Request, asduPack.CommonAddr)

	return nil
}
func (sf *mysrv) ClockSyncHandler(asdu.Connect, *asdu.ASDU, time.Time) error {
	//时钟同步
	fmt.Println("call ClockSyncHandler()")
	return nil
}
func (sf *mysrv) ResetProcessHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierOfResetProcessCmd) error {
	//重置进程
	fmt.Println("call ResetProcessHandler()")
	return nil
}
func (sf *mysrv) DelayAcquisitionHandler(asdu.Connect, *asdu.ASDU, uint16) error {
	//延迟
	fmt.Println("call DelayAcquisitionHandler()")
	return nil
}

func (sf *mysrv) ASDUHandler(asdu.Connect, *asdu.ASDU) error {
	//下发指令
	fmt.Println("call ASDUHandler()")
	return nil
}
