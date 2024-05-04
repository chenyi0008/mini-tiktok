package helper

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
	"log"
	"mini-tiktok/service/core/define"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"time"
)

func GenerateToken(id uint, username string, second int) (string, error) {
	uc := define.UserClaim{
		Id:       id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(second)).Unix(),
		},
	}
	//fmt.Printf("uc: %+v", uc)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenString, err := token.SignedString([]byte("mini-tiktok"))
	if err != nil {
		log.Println("err:", err)
		return "", err
	}

	//tokenString := fmt.Sprintf("%v", token)
	return tokenString, nil
}

func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func AnalyzeToken(token string) (*define.UserClaim, error) {
	uc := new(define.UserClaim)
	claims, err := jwt.ParseWithClaims(token, uc, func(token *jwt.Token) (interface{}, error) {
		return []byte(define.JwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return uc, errors.New("token is invalid")
	}
	return uc, err
}

// CosUpload 文件上传到腾讯云
func CosUpload(r *http.Request) (string, error) {
	u, _ := url.Parse(define.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})

	file, fileHeader, err := r.FormFile("data")

	// Ext取后缀 扩展名 extension
	key := define.TencentFilePrefix + GetUUID() + path.Ext(fileHeader.Filename)

	if err != nil {
		panic(err)
	}

	_, err = client.Object.Put(
		context.Background(), key, file, nil,
	)

	if err != nil {
		panic(err)
	}
	return define.CosBucket + "/" + key, nil
}

func FileUploadToJpg(data *[]byte) (string, error) {
	u, _ := url.Parse(define.CosBucket)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  define.TencentSecretID,
			SecretKey: define.TencentSecretKey,
		},
	})

	// Ext取后缀 扩展名 extension
	key := define.TencentFilePrefix + GetUUID() + ".jpg"

	_, err := client.Object.Put(
		context.Background(), key, bytes.NewReader(*data), nil,
	)

	if err != nil {
		panic(err)
	}
	return define.CosBucket + "/" + key, nil
}

func GetUUID() string {
	return uuid.NewV4().String()
}

// Ffmpeg 视频封面截取
func Ffmpeg(videoURL string, frameNum int) ([]byte, error) {
	// 创建一个临时文件来存储输出图像
	outputFile := "output.jpg"

	// 使用 ffmpeg 从视频中获取指定帧并将其输出到临时文件
	cmd := exec.Command("ffmpeg",
		"-i", videoURL,
		"-vf", fmt.Sprintf("select='gte(n,%d)'", frameNum),
		"-vframes", "1",
		outputFile)

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// 读取临时文件的内容到缓冲区
	buf, err := os.ReadFile(outputFile)
	if err != nil {
		log.Fatal(err)
	}

	// 删除临时文件
	err = os.Remove(outputFile)
	if err != nil {
		log.Println("Error removing temporary file:", err)
	}

	return buf, nil
}

// TransformUnixToDate 将时间戳（int64）转化为2023-09-24T04:45:05+08:00（string）
func TransformUnixToDate(date int64) string {
	timeTemplate := "2006-01-02 15:04:05"
	return time.Unix(date, 0).Format(timeTemplate)
}
