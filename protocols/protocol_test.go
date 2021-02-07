package protocols

import (
	"fmt"
	"testing"
)

func TestProtocol_AppendNumber(t *testing.T){
	var pb Protocol
	pb.AppendNumber(true)
	pb.AppendNumber(int8(1))
	pb.AppendNumber(uint16(2))
	pb.AppendNumber(int32(4))
	pb.AppendNumber(float32(4.4))
	fmt.Println("content:",pb.Content)
}

func TestProtocol_AppendStringUint8(t *testing.T){
	str := "hello world"
	var pb Protocol
	pb.AppendStringUint8(str)
	//buff := make([]byte, 300)
	//buff[1] = 1
	//pb.AppendStringUint8(string(buff))

	fmt.Println("len:", len(str), "content:",pb.Content)
}

func TestProtocol_AppendStringUint16(t *testing.T){
	str := "hello world"
	var pb Protocol
	pb.AppendStringUint16(str)
	//buff := make([]byte, 70000)
	//buff[1] = 1
	//pb.AppendStringUint16(string(buff))

	fmt.Println("len:", len(str), "content:",pb.Content)
}

func TestProtocol_GetNumber(t *testing.T){
	var pb Protocol
	pb.AppendNumber(ProtocolHeader{3, 4})
	pb.AppendNumber(float32(4.4))
	var (
		cmd uint16
		length uint16
		f float32
	)
	err := pb.GetNumber(&cmd)
	if err != nil {
		fmt.Println(err)
	}
	pb.GetNumber(&length)
	if err != nil {
		fmt.Println(err)
	}
	err = pb.GetNumber(&f)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("cmd:",cmd,", len:", length, "float:", f)
}

func TestProtocol_GetStringUint8(t *testing.T){
	str := "hello world"
	var pb Protocol
	pb.AppendStringUint8(str)

	newStr, _:= pb.GetStringUint8()
	fmt.Println("new:",newStr)
}

func TestProtocol_GetStringUint16(t *testing.T){
	str := "hello world"
	var pb Protocol
	pb.AppendStringUint16(str)

	newStr, _:= pb.GetStringUint16()
	fmt.Println("new:",newStr)
}

func TestProtocol_Bytes(t *testing.T){
	str := "hello world"
	var pb Protocol
	pb.AppendNumber(true)
	pb.AppendNumber(int8(1))
	pb.AppendNumber(uint16(2))
	pb.AppendNumber(float32(4.4))
	pb.AppendStringUint16(str)
	fmt.Println("buff:",pb.Bytes(3))
	var (
		b bool
		i8 int8
		u16 uint16
		f32 float32
	)
	pb.GetNumber(&b)
	pb.GetNumber(&i8)
	pb.GetNumber(&u16)
	pb.GetNumber(&f32)
	newStr, _:= pb.GetStringUint16()
	fmt.Println("bool:",b)
	fmt.Println("int8:",i8)
	fmt.Println("uint16:",u16)
	fmt.Println("float32:",f32)
	fmt.Println("string:",newStr)
}
