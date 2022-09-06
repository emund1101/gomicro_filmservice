package handler

import (
	"bytes"
	"context"
	"crypto/md5"
	pb "films/service_user/proto"
	"fmt"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
	"time"
)

//写实现方法

type UserService struct {
	Db  *gorm.DB
	Rdb *redis.Client
}

type user struct {
	Id       int
	Mobile   string
	Name     string
	Username string
	Password string
	Avadar   string
	Birthday string
}

type json_return struct {
	Code int
	Msg  string
	Data any
}

//注册实现
func (e *UserService) Register(ctx context.Context, req *pb.RegisterRequest, rsp *pb.Respone) error {
	pwd := md5.Sum([]byte(req.Password))
	p := fmt.Sprintf("%x", pwd) //将[16]byte 转string

	register := user{
		Mobile:   req.Mobile,
		Name:     req.Name,
		Username: req.Name,
		Password: p,
		Avadar:   req.Avadar,
		Birthday: req.Birthday,
	}
	e.Db.Create(&register)
	if register.Id != 0 {
		rsp.Respone = []byte("注册成功")
	} else {
		rsp.Respone = []byte("注册失败")
	}
	return nil
}

//登录的实现
func (e *UserService) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.Respone) error {
	rst := user{}
	e.Db.Table("user").Where("mobile = ?", req.Mobile).Find(&rst)

	if rst.Id == 0 {
		json, _ := jsoniter.Marshal(json_return{0, "登录失败", nil})
		rsp.Respone = json
		return nil
	}

	var md5_auth string

	//匿名函数
	f := func() error {
		rstj, _ := jsoniter.Marshal(rst)
		auth := md5.Sum(rstj)
		md5_auth = fmt.Sprintf("%x", auth) //16进制的byte
		err := e.Rdb.Set(ctx, md5_auth, string(rstj), 86400*time.Second).Err()
		if err != nil {
			fmt.Println(err)
			return err
		} else {
			return nil
		}
	}

	switch req.Type {
	case 1:
		if req.Code == "1111" {
			f()
			json, _ := jsoniter.Marshal(json_return{1, "登录成功", map[string]string{"auth": md5_auth}})
			rsp.Respone = json
		} else {
			json, _ := jsoniter.Marshal(json_return{0, "登录失败", nil})
			rsp.Respone = json
		}
	case 2:
		pwd := md5.Sum([]byte(req.Password))
		p := make([]byte, 0)
		p = append(p, pwd[0:]...) //将[16]byte 转[]byte
		if rst.Password == bytes.NewBuffer(p).String() {
			f()
			json, _ := jsoniter.Marshal(json_return{1, "登录成功", map[string]string{"auth": md5_auth}})
			rsp.Respone = json
		} else {
			json, _ := jsoniter.Marshal(json_return{0, "登录失败", nil})
			rsp.Respone = json
		}
	}

	return nil

}

//验证授权登录
func (e *UserService) Authorize(ctx context.Context, req *pb.AuthRequest, rsp *pb.Respone) error {
	//auth
	val, err := e.Rdb.Get(ctx, req.Auth).Result()
	rst := user{}
	e.Db.Table("user").Where("mobile = ?", "18111111111").Find(&rst)

	if err != nil {
		fmt.Println(err)
		json, _ := jsoniter.Marshal(json_return{0, "没有数据", nil})
		rsp.Respone = json
		return err
	} else {
		rst := user{}
		jsoniter.Unmarshal(bytes.NewBufferString(val).Bytes(), &rst)
		json, _ := jsoniter.Marshal(json_return{1, "", rst})
		rsp.Respone = json
		return nil

	}

}
