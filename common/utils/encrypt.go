/*
 * 纸喵软件
 * Copyright (c) 2017~2020 http://zhimiao.org All rights reserved.
 * Author: 倒霉狐狸 <mail@xiaoliu.org>
 * Date: 2020/3/3 下午4:26
 */

package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"gitee.com/zhimiao/wechat/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// MD5 md5 encryption
func MD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

//ParseToken 解析jwtToken
func ParseToken(tokenString string) (string, error) {
	jwtSecret := []byte(common.Config.App.JwtSecret)
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*jwt.StandardClaims); ok && tokenClaims.Valid {
			if !claims.VerifyExpiresAt(time.Now().Unix(), false) {
				return "", fmt.Errorf("过期了")
			}
			if claims.Issuer != "zhimiao-wechat" {
				return "", fmt.Errorf("非法来源的签名")
			}
			return claims.Subject, nil
		}
	}
	return "", err
}

//CreateToken 生成jwtToken
func CreateToken(subject string, expire time.Duration) (string, error) {
	jwtSecret := []byte(common.Config.App.JwtSecret)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   subject,
		ExpiresAt: time.Now().Add(expire).Unix(),
		Issuer:    "zhimiao-wechat",
	})
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// 密码生成
func PasswordHash(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		logrus.Warn(err.Error())
	}
	return string(hash)
}

// 密码验证
func PasswordVerify(hashedPwd string, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		logrus.Warn(err.Error())
		return false
	}
	return true
}

//---------------DES加密  解密--------------------
func EncyptogAES(src, key string) string {
	s := []byte(src)
	k := []byte(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		print(err.Error())
		return ""
	}
	blockSize := block.BlockSize()
	paddingCount := blockSize - len(s)%blockSize
	//填充数据为：paddingCount ,填充的值为：paddingCount
	paddingStr := bytes.Repeat([]byte{byte(paddingCount)}, paddingCount)
	new_s := append(s, paddingStr...)
	blockMode := cipher.NewCBCEncrypter(block, []byte(key))
	blockMode.CryptBlocks(new_s, new_s)
	return string(new_s)

}
func DecrptogAES(src, key string) string {
	s := []byte(src)
	k := []byte(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		return ""
	}
	blockMode := cipher.NewCBCDecrypter(block, k)
	blockMode.CryptBlocks(s, s)
	n := len(s)
	count := int(s[n-1])
	return string(s[:n-count])
}
