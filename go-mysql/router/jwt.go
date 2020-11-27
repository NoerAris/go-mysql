package router

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/errors"
	"io"
	"log"
	"strings"
	"time"
)

const RolePublic = `ROLE_P00`
const RoleMemberRegistered = `ROLE_S01`
const RoleMemberApproved = `ROLE_S02`

const JwtHeaderKey = `Authorization`

//const JwtUserId = `uid`
const JwtObaMemberNumber = `oba`
const JwtEmail = `eml`
const JwtBranch = `bid`
const JwtRoles = `rol`

//const JwtUniqueId = `qid`
//const JwtSubject = `sub`

const RsaSecretKey = "4$up3rDup3r$ecr3tK3y234"
const JwtAesDecryptKey = "E1BB465D57CAE78091F9CE83DF"

//const SIGNATURE = "gFCWe4saM1HcIf8qqaPjnbh5aiVFPyuuLMtcFhPDntrnqa0Ad52VIHKFyjo1HpCau6DaOQJHhmqge7cXV0ScbA"

var Jakarta *time.Location

func init() {
	Jakarta, _ = time.LoadLocation("Asia/Jakarta")
}

type JwtModel struct {
	Email     string
	Branch    string
	ObaNumber string
	RawToken  string
}

func GetClaim(jwtoken string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(jwtoken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method : %v", token.Header["alg"])
		}
		return []byte(RsaSecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
func CreateJwt(mapClaims jwt.MapClaims) (string, error) {
	mapClaims["nbf"] = time.Now().In(Jakarta)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(RsaSecretKey))
	if err != nil {
		return "Error Creating Key", err
	}
	return tokenString, nil
}
func CreateJwtPublic() (string, error) {
	return CreateJwt(jwt.MapClaims{})
}
func JwtValidator(ctx context.Context) {
	claims, err := GetClaimFromContext(ctx)
	if err != nil {
		ctx.StatusCode(RESPONSE_CODE_UNAUTHORIZED)
		ctx.JSON(GetErrorResponse(err.Error()))
	} else {
		if err := claims.Valid(); err != nil {
			ctx.StatusCode(RESPONSE_CODE_UNAUTHORIZED)
			ctx.JSON(GetErrorResponse(err.Error()))
			return
		}
		ctx.Next()
	}
}
func GetClaimFromContext(ctx context.Context) (jwt.MapClaims, error) {
	authToken := ctx.GetHeader(JwtHeaderKey)
	claims, err := GetClaim(strings.Replace(authToken, "Bearer ", "", 1))
	return claims, err
}
func GetClaimObjectFromContext(ctx context.Context) (*JwtModel, error) {
	claim, err := GetClaimFromContext(ctx)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed get claim from header : %v", err))
	}
	branchId, check := claim[JwtBranch].(string)
	if !check {
		return nil, errors.New(fmt.Sprintf("claim doesnt have branch_id value : %v", check))
	}
	branch_id, err := DecryptClaimValue(branchId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("claim doesnt have branch_id value : %v", err))
	}
	emailString, check := claim[JwtEmail].(string)
	if !check {
		return nil, errors.New(fmt.Sprintf("claim doesnt have email value : %v", check))
	}
	email, err := DecryptClaimValue(emailString)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("claim doesnt have email value : %v", err))
	}
	obaNumberString, check := claim[JwtObaMemberNumber].(string)
	if !check {
		return nil, errors.New(fmt.Sprintf("claim doesnt have oba number value : %v", check))
	}
	obaNumber, err := DecryptClaimValue(obaNumberString)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("claim doesnt have oba number value : %v", err))
	}
	return &JwtModel{Email: email, Branch: branch_id, ObaNumber: obaNumber, RawToken: ctx.GetHeader(JwtHeaderKey)}, nil
}

func DesEncryption(key, iv, plainText []byte) ([]byte, error) {

	block, err := des.NewCipher(key)

	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	if len(plainText)%blockSize != 0 {
		panic("Need a multiple of the blocksize")
	}
	ciphertext := make([]byte, len(PKCS5Padding(plainText, blockSize)))
	for len(plainText) > 0 {
		block.Encrypt(ciphertext, PKCS5Padding(plainText, blockSize))
		plainText = plainText[blockSize:]
		ciphertext = ciphertext[blockSize:]
	}
	return ciphertext, nil
}

func AesDecryption(key, iv, cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {

		log.Println("error create new chiper:=", err)
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(cipherText))
	blockMode.CryptBlocks(origData, cipherText)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}
func AesEncryption(key, source []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {

		log.Println("error create new chiper:=", err)
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(source))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	blockMode := cipher.NewCBCEncrypter(block, iv)
	source = PKCS5Padding(source, block.BlockSize())
	origData := make([]byte, len(source))
	blockMode.CryptBlocks(origData, source)
	result := append(iv, origData...)
	return result, nil
}
func DecryptClaimValue(value string) (string, error) {
	encryptedText := value
	encryptedData, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		log.Println("error base64 decode:=", err)
		return "", err
	}
	key, err := hex.DecodeString(JwtAesDecryptKey)
	if err != nil {
		log.Println("error hex decode:=", err)
		return "", err
	}
	res, err := AesDecryption(key, encryptedData[0:16], encryptedData[16:])
	if err != nil {
		log.Println("error decrypt:=", err)
		return "", err
	}
	return string(res[:]), nil
}
func EncryptClaimValue(value string) (string, error) {
	key, err := hex.DecodeString(JwtAesDecryptKey)
	if err != nil {
		log.Println("error hex decode:=", err)
		return "", err
	}
	res, err := AesEncryption(key, []byte(value))
	if err != nil {
		log.Println("error decrypt:=", err)
		return "", err
	}

	encryptedData := base64.StdEncoding.EncodeToString(res)
	if err != nil {
		log.Println("error base64 encode:=", err)
		return "", err
	}
	return encryptedData, nil
}
func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
