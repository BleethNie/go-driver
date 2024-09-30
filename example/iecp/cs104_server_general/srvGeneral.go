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

	///*高路2024-9-17start
	asduPack.SendReplyMirror(c, asdu.ActivationCon)

	objSize, _ := asdu.GetInfoObjSize(asduPack.Type)
	objSize += asduPack.InfoObjAddrSize

	var yxid int = 1
	switch YxType {
	case 1:
	case 2:
		yxid = 2

	}
	/*var t time.Time
	switch asduPack.Type {
	case asdu.M_SP_NA_1:
	case asdu.M_SP_TA_1:
		t = asduPack.DecodeCP24Time2a()
		objSize = objSize + 3
	case asdu.M_SP_TB_1:
		t = asduPack.DecodeCP56Time2a()
		objSize = objSize + 7
	}*/

	infolenth := getInfoCount(asduPack.Params, objSize, asduPack.Variable.IsSequence)
	var k int = 0
	if yxid == 1 {
		var info []asdu.SinglePointInfo
		for k, v := range AppMap["yx"] {
			var spi asdu.SinglePointInfo
			spi.Ioa = asdu.InfoObjAddr(k)
			//spi.Time = t
			spi.Value = YXMap[v] == 1
			spi.Qds = 0x00
			info = append(info, spi)
		}

		for range info {

			if k+infolenth < len(info) {

				err := asdu.Single(c, asduPack.Variable.IsSequence, asdu.CauseOfTransmission{Cause: asdu.InterrogatedByStation}, asduPack.CommonAddr, info[k:k+infolenth]...)
				if err != nil {
					return err
				}
			} else {
				err := asdu.Single(c, asduPack.Variable.IsSequence, asdu.CauseOfTransmission{Cause: asdu.InterrogatedByStation}, asduPack.CommonAddr, info[k:]...)
				if err != nil {
					return err
				}
			}

			k = k + infolenth
		}
	} else {
		var info []asdu.DoublePointInfo
		for k, v := range AppMap["yx"] {
			var spi asdu.DoublePointInfo
			spi.Ioa = asdu.InfoObjAddr(k)
			//spi.Time = t
			if YXMap[v] == 1 {
				spi.Value = asdu.DPIDeterminedOn
			} else {
				spi.Value = asdu.DPIDeterminedOff
			}
			spi.Qds = 0x00
			info = append(info, spi)
		}

		for range info {
			if k+infolenth < len(info) {
				err := asdu.Double(c, asduPack.Variable.IsSequence, asdu.CauseOfTransmission{Cause: asdu.InterrogatedByStation}, asduPack.CommonAddr, info[k:k+infolenth]...)
				if err != nil {
					return err
				}
			} else {
				err := asdu.Double(c, asduPack.Variable.IsSequence, asdu.CauseOfTransmission{Cause: asdu.InterrogatedByStation}, asduPack.CommonAddr, info[k:]...)
				if err != nil {
					return err
				}
				break
			}

			k = k + infolenth
		}
	}

	//遥测总召

	var ycid int = 1
	switch YCType {
	case 1:
	case 2:
		ycid = 2

	}

	if ycid == 1 {
		infoObjSize, _ := asdu.GetInfoObjSize(asdu.M_ME_NC_1)
		infoObjSize += asduPack.InfoObjAddrSize
		ycinfolength := getInfoCount(asduPack.Params, infoObjSize, asduPack.Variable.IsSequence)
		var ycinfo []asdu.MeasuredValueFloatInfo
		for k, v := range AppMap["yc"] {
			var spi asdu.MeasuredValueFloatInfo
			spi.Ioa = asdu.InfoObjAddr(k) + YcAddress
			//spi.Time = t
			spi.Value = YCMap[v+YcAddress]
			spi.Qds = 0x00
			ycinfo = append(ycinfo, spi)
		}
		var k = 0
		for range ycinfo {
			if k+ycinfolength < len(ycinfo) {
				err := asdu.MeasuredValueFloat(c, asduPack.Variable.IsSequence, asdu.CauseOfTransmission{Cause: asdu.InterrogatedByStation}, asduPack.CommonAddr, ycinfo[k:k+ycinfolength]...)
				if err != nil {
					return err
				}
			} else {
				err := asdu.MeasuredValueFloat(c, asduPack.Variable.IsSequence, asdu.CauseOfTransmission{Cause: asdu.InterrogatedByStation}, asduPack.CommonAddr, ycinfo[k:]...)
				if err != nil {
					return err
				}
			}

			k = k + ycinfolength
		}

	} else if ycid == 2 {

		fmt.Println("归一化传输~~")
	}
	/*
		//高路2024-9-17end
		//*/

	// go func() {
	// 	for {
	// 		err := asdu.Single(c, false, asdu.CauseOfTransmission{Cause: asdu.Spontaneous}, asdu.GlobalCommonAddr,
	// 			asdu.SinglePointInfo{})
	// 		if err != nil {
	// 			log.Println("falied", err)
	// 		} else {
	// 			log.Println("success", err)
	// 		}

	// 		time.Sleep(time.Second * 1)
	// 	}
	// }()

	//asduPack.SendICReply(c, asdu.InterrogatedByStation)
	err := asduPack.SendReplyMirror(c, asdu.ActivationTerm)
	if err != nil {
		return err
	}
	return nil
}
func (sf *mysrv) CounterInterrogationHandler(asdu.Connect, *asdu.ASDU, asdu.QualifierCountCall) error {
	fmt.Println("call CounterInterrogationHandler()")
	return nil
}
func (sf *mysrv) ReadHandler(c asdu.Connect, asduPack *asdu.ASDU, addr asdu.InfoObjAddr) error {
	fmt.Println("call ReadHandler()")
	///*高路2024-9-17start
	objAddr := asduPack.DecodeInfoObjAddr()
	/*var t time.Time
	asduPack.Type = asdu.M_SP_TA_1
	switch asduPack.Type {
	case asdu.M_SP_NA_1:
	case asdu.M_SP_TA_1:
		t = asduPack.DecodeCP24Time2a()
	case asdu.M_SP_TB_1:
		t = asduPack.DecodeCP56Time2a()
	}*/
	if objAddr > YcAddress {
		var spi asdu.MeasuredValueFloatInfo
		spi.Ioa = objAddr
		//spi.Time = t
		spi.Value = YCMap[int(objAddr)]
		spi.Qds = 0x00
		switch asduPack.Type {
		case asdu.M_ME_NC_1:
			err := asdu.MeasuredValueFloat(c, asduPack.Variable.IsSequence, asdu.CauseOfTransmission{Cause: asdu.Request}, asduPack.CommonAddr, spi)
			if err != nil {
				return err
			}
		case asdu.M_ME_TC_1:
			err := asdu.MeasuredValueFloatCP24Time2a(c, asdu.CauseOfTransmission{Cause: asdu.Request}, asduPack.CommonAddr, spi)
			if err != nil {
				return err
			}
		case asdu.M_ME_TF_1:
			err := asdu.MeasuredValueFloatCP56Time2a(c, asdu.CauseOfTransmission{Cause: asdu.Request}, asduPack.CommonAddr, spi)
			if err != nil {
				return err
			}
		}
	} else {
		var yxid int = 1
		switch YxType {
		case 1:
		case 2:
			yxid = 2

		}
		if yxid == 1 { //单点
			var spi asdu.SinglePointInfo
			spi.Ioa = objAddr
			//spi.Time = t
			spi.Value = YXMap[int(objAddr)] == 1
			spi.Qds = 0x00
			switch asduPack.Type {
			case asdu.M_SP_NA_1:
				err := asdu.Single(c, asduPack.Variable.IsSequence, asdu.CauseOfTransmission{Cause: asdu.Request}, asduPack.CommonAddr, spi)
				if err != nil {
					return err
				}
			case asdu.M_SP_TA_1:
				err := asdu.SingleCP24Time2a(c, asdu.CauseOfTransmission{Cause: asdu.Request}, asduPack.CommonAddr, spi)
				if err != nil {
					return err
				}
			case asdu.M_SP_TB_1:
				err := asdu.SingleCP56Time2a(c, asdu.CauseOfTransmission{Cause: asdu.Request}, asduPack.CommonAddr, spi)
				if err != nil {
					return err
				}
			}
		} else {
			//双点遥信
			var spi asdu.DoublePointInfo
			spi.Ioa = objAddr
			//spi.Time = t
			if YXMap[int(objAddr)] == 1 {
				spi.Value = asdu.DPIDeterminedOn
			} else {
				spi.Value = asdu.DPIDeterminedOff
			}
			spi.Qds = 0x00
			switch asduPack.Type {
			case asdu.M_DP_NA_1:
				err := asdu.Double(c, asduPack.Variable.IsSequence, asdu.CauseOfTransmission{Cause: asdu.InterrogatedByStation}, asduPack.CommonAddr, spi)
				if err != nil {
					return err
				}
			case asdu.M_DP_TA_1:
				err := asdu.DoubleCP24Time2a(c, asdu.CauseOfTransmission{Cause: asdu.Request}, asduPack.CommonAddr, spi)
				if err != nil {
					return err
				}
			case asdu.M_DP_TB_1:
				err := asdu.DoubleCP56Time2a(c, asdu.CauseOfTransmission{Cause: asdu.Request}, asduPack.CommonAddr, spi)
				if err != nil {
					return err
				}
			}
			/////////////////////////////
		}
	}

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

func (sf *mysrv) ASDUHandler(c asdu.Connect, asduPack *asdu.ASDU) error {
	fmt.Println("call ASDUHandler()")
	switch asduPack.Identifier.Type {
	case asdu.C_SC_NA_1, asdu.C_SC_TA_1: //单命令遥控
		if asduPack.Identifier.Coa.Cause == asdu.Activation {
			singleCommandInfo := asduPack.GetSingleCmd()
			// InSelect: true - selects, false - executes.
			if singleCommandInfo.Qoc.InSelect {
				return asduPack.SendReplyMirror(c, asdu.ActivationCon)
			} else {
				fmt.Println("调用遥控处理方法进行遥控操作，需要在执行完操作后发送一个激活终止的消息")
				return asduPack.SendReplyMirror(c, asdu.ActivationCon)
			}
		} else if asduPack.Identifier.Coa.Cause == asdu.Deactivation {
			fmt.Println("遥控中的撤销需要做什么操作")
			return asduPack.SendReplyMirror(c, asdu.ActivationCon)

		} else {
			return asduPack.SendReplyMirror(c, asdu.UnknownCOT)
		}
	case asdu.C_DC_NA_1, asdu.C_DC_TA_1: //双命令遥控
		if asduPack.Identifier.Coa.Cause == asdu.Activation {
			doubleCommandInfo := asduPack.GetDoubleCmd()
			// InSelect: true - selects, false - executes.
			if doubleCommandInfo.Qoc.InSelect {
				return asduPack.SendReplyMirror(c, asdu.ActivationCon)
			} else {
				fmt.Println("调用遥控处理方法进行遥控操作，需要在执行完操作后发送一个激活终止的消息")
				return asduPack.SendReplyMirror(c, asdu.ActivationCon)
			}
		} else if asduPack.Identifier.Coa.Cause == asdu.Deactivation {
			fmt.Println("遥控中的撤销需要做什么操作")
			return asduPack.SendReplyMirror(c, asdu.ActivationCon)

		} else {
			return asduPack.SendReplyMirror(c, asdu.UnknownCOT)
		}
	}
	return nil
}

func getInfoCount(param *asdu.Params, objSize int, isSequence bool) int {

	var infoLen int = 0
	if isSequence {
		infoLen = (asdu.ASDUSizeMax - (param.IdentifierSize() + param.InfoObjAddrSize)) / objSize
	} else {
		infoLen = (asdu.ASDUSizeMax - param.IdentifierSize()) / objSize
	}

	return infoLen

}
