package util

import (
	"math/rand"
	"regexp"
)

var _letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

//SubStr doc
//@Method SubStr @SummaryCut string ends with length
//@Param  (string) source string
//@Param  (int)    start pos
//@Param  (int)    sub length
//@Return (string) sub string
func SubStr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

//SubStr2 doc
//@Method SubStr2 doc : Cut string ends with index
//@Param  (string) source string
//@Param  (int)    start pos
//@Param  (int)    end pos
//@Return (string) sub string
func SubStr2(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}
	return string(rs[start:end])
}

//RandStr doc
//@Method RandStr doc : Randomly generate a string of length n
//@Param (int) length
//@Return (string)
func RandStr(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = _letterRunes[rand.Intn(len(_letterRunes))]
	}
	return string(b)
}

//VerifyPasswordFormat doc
//@Summary Verify password is valid
//@Param (string) password
//@Return (bool) is valid
func VerifyPasswordFormat(pwd string) bool {
	b, e := regexp.MatchString("^([a-zA-Z_-].*)([0-9].*)$", pwd)
	if e != nil {
		panic(e)
	}

	if b {
		if len(pwd) >= 8 && len(pwd) <= 16 {
			return true
		}
		return false
	}

	return false
}

//VerifyAccountFormat doc
//@Summary Verify account is valid
//@Param (string) account
//@Return (bool) is valid
func VerifyAccountFormat(account string) bool {
	b, e := regexp.MatchString("^[a-zA-Z0-9_-]{8,16}$", account)
	if e != nil {
		panic(e)
	}
	return b
}

//VerifyCaptchaFormat doc
//@Summary Verify captcha is valid
//@Param (string) captcha
//@Return (bool) is valid
func VerifyCaptchaFormat(captcha string) bool {
	b, e := regexp.MatchString("^[0-9]{6,6}", captcha)
	if e != nil {
		panic(e)
	}
	return b
}

//VerifyEmailFormat doc
//@Summary Verify is email
//@Param (string) email
//@Return (bool) is valid
func VerifyEmailFormat(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

//VerifyMobileFormat doc
//@Summary Verify is mobile
//@Param (string) mobile
//@Return (bool) is valid
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}
