package strs

import (
	"math/rand"
	"strconv"
	"time"
	"unsafe"
)

//NumRand 随机数字字符串
func NumRand(l int) string {
	//将时间戳设置成种子数
	rand.Seed(time.Now().UnixNano())
	var yzm string
	//生成10个0-99之间的随机数
	for i := 0; i < l; i++ {
		yzm += strconv.Itoa(rand.Intn(10))
	}

	return yzm
}

//StrRand 随机字符串
func StrRand(n int) string {
	const (
		letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)

	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
