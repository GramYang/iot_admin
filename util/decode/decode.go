package decode

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"iot_admin/util"
	"os"
)

var Signature= getSignature()
var privateKey= getPrivateKey()

func BodyDecodeRSA(src []byte)([]byte,error){
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey,src)
}

func getPrivateKey() *rsa.PrivateKey{
	path:="util/decode/rsa_2048_priv.pem"
	file,err:=os.Open(path)
	if err!=nil{
		panic(err)
	}
	defer file.Close()
	info,_:=file.Stat()
	buf:=make([]byte,info.Size())
	_,_=file.Read(buf)
	block,_:=pem.Decode(buf)
	if block==nil{
		panic("getPrivateKey error")
	}
	privateKey,err:=x509.ParsePKCS1PrivateKey(block.Bytes)
	if err!=nil{
		panic("getPrivateKey error")
	}
	return privateKey
}

func getSignature() string{
	ctx:=md5.New()
	ctx.Write([]byte(util.Md5Key))
	return hex.EncodeToString(ctx.Sum(nil))
}