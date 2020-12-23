package decode

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"iot_admin/util"
	"os"
)

var Signature = getSignature()
var privateKey = getPrivateKey()

//普通的rsa解密，加密文本长度不能超过秘钥长度
func BodyDecodeRSA(src []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, src)
}

//分区的rsa解密
func BodyDecodeRSADivisional(src []byte) ([]byte, error) {
	keySize, srcSize := privateKey.Size(), len(src)
	var offSet = 0
	var buffer = bytes.Buffer{}
	for offSet < srcSize {
		//endIndex的初值是keySize，并以keySize累加
		endIndex := offSet + keySize
		if endIndex > srcSize {
			endIndex = srcSize
		}
		bytesOnce, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, src[offSet:endIndex])
		if err != nil {
			return nil, err
		}
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	return buffer.Bytes(), nil
}

func getPrivateKey() *rsa.PrivateKey {
	path := "util/decode/rsa_2048_priv.pem"
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	_, _ = file.Read(buf)
	block, _ := pem.Decode(buf)
	if block == nil {
		panic("getPrivateKey error")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic("getPrivateKey error")
	}
	return privateKey
}

func getSignature() string {
	ctx := md5.New()
	ctx.Write([]byte(util.Md5Key))
	return hex.EncodeToString(ctx.Sum(nil))
}
