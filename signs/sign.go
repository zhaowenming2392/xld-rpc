package signs

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"

	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

//MD5 MD5
func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	return hex.EncodeToString(has[:])
}

//Sha256 Sha256
func Sha256(str string) string {
	bytes := sha256.Sum256([]byte(str)) //计算哈希值，返回一个长度为32的数组
	return hex.EncodeToString(bytes[:]) //将数组转换成切片，转换成16进制，返回字符串
}

//HashSha256 sha256
func HmacSha256(str, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//ValidHashSha256MAC sha256
func ValidHashSha256MAC(message, messageMAC, key string) bool {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	expectedMAC := mac.Sum(nil)
	return hmac.Equal([]byte(messageMAC), expectedMAC)
}

//AES加密过程：
//  1、处理数据，对数据进行填充，采用PKCS7（当密钥长度不够时，缺几位补几个几，即使正好也要补）的方式。
//  2、对数据进行加密，采用AES加密方法中CBC加密模式
//  3、对得到的加密数据，进行base64加密，得到字符串
// 解密过程相反，key不能泄露

//pkcs7Padding 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	//判断缺少几位长度。最少1，最多 blockSize
	//len(data)%blockSize：当len(data)<blockSize=len(data)，否则是mod计算
	padding := blockSize - len(data)%blockSize //len(data)=blockSize，也要填充blockSize个
	//补足位数。把切片[]byte{byte(padding)}复制padding个；使用缺少的字节数作为需要补的字节中的内容，然后补齐字节数个
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

//pkcs7UnPadding 填充的反向操作
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("密文有误！") //肯定是有填充数量的，不可能是0
	}
	//获取填充的个数
	unPadding := int(data[length-1]) //最后一个作为数量
	return data[:(length - unPadding)], nil
}

//AesCbcEncrypt 加密
//key=16,24,32位字符串的话，分别对应AES-128，AES-192，AES-256 加密方法
func AesCbcEncrypt(data, key []byte) ([]byte, error) {
	//创建加密实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//判断加密快的大小
	blockSize := block.BlockSize()
	//填充
	encryptBytes := pkcs7Padding(data, blockSize)
	//初始化加密数据接收切片
	crypted := make([]byte, len(encryptBytes))
	//使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	//执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

//AesCbcDecrypt 解密
func AesCbcDecrypt(data, key []byte) ([]byte, error) {
	//创建实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//使用cbc
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	//初始化解密数据接收切片
	crypted := make([]byte, len(data))
	//执行解密
	blockMode.CryptBlocks(crypted, data)
	//去除填充
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}

//EncryptByAes Aes加密
//加密后 base64 再加，方便传输
//data 原文, key密钥，返回密文
func EncryptByAesCbc(data, key string) (string, error) {
	res, err := AesCbcEncrypt([]byte(data), []byte(key))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

//DecryptByAes Aes 解密
//data 密文, key密钥，返回原文
func DecryptByAesCbc(data, key string) (string, error) {
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	str, err := AesCbcDecrypt(dataByte, []byte(key))
	if err != nil {
		return "", err
	}
	return string(str), nil
}
