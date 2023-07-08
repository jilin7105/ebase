package main

import (
	"context"
	"fmt"
	"github.com/jilin7105/ebase"
	"github.com/jilin7105/ebase/examp/grpc/proto/user"
	"github.com/jilin7105/ebase/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"os"
	"time"
)

type AuthService struct {
}

func NewService() *AuthService {
	return &AuthService{}
}

type User struct {
	Id       int32
	Name     string
	Password string
}

func (*AuthService) Login(ctx context.Context, name, password string) (User, error) {
	return User{
		Id:       1,
		Name:     name,
		Password: password,
	}, nil
}

type AuthClr struct {
	user.UnimplementedAuthServiceServer
}

func validateLoginRequest(in *user.LoginRequest) error {
	return nil
}

func NewAuthClr() *AuthClr {
	return &AuthClr{}
}

func (a AuthClr) Login(ctx context.Context, in *user.LoginRequest) (*user.LoginResponse, error) {
	if err := validateLoginRequest(in); err != nil {
		return nil, err
	}
	logger.Info("登录接口被调用 request_id : %s", ctx.Value("eb-grpc-request-id"))
	modeluser, err := NewService().Login(ctx, in.Username, in.Password)
	if err != nil {
		return nil, err
	}
	return &user.LoginResponse{
		User: &user.User{
			Id:   modeluser.Id,
			Name: modeluser.Name,
		},
	}, nil
}

func (a AuthClr) Register(ctx context.Context, in *user.RegisterRequest) (*user.RegisterResponse, error) {
	panic("implement me")
}

//使用go run main.go  启动测试服务
func main() {
	path, _ := os.Getwd()
	ebase.SetProjectPath(path)
	ebase.Init()
	eb := ebase.GetEbInstance()
	s, err := eb.GetGrpcServer()
	if err != nil {
		panic(err)
	}
	//注册
	user.RegisterAuthServiceServer(s, NewAuthClr())

	//测试一下
	go func(port int) {
		time.Sleep(5) //保证服务启动
		cc, err := grpc.Dial(fmt.Sprintf(":%d", port), grpc.WithInsecure())
		if err != nil {
			panic(err)
		}

		cli := user.NewAuthServiceClient(cc)

		md := metadata.Pairs("eb-grpc-request-id", "test1")
		// 创建带metadata的上下文
		cxt := metadata.NewOutgoingContext(context.Background(), md)
		login, err := cli.Login(cxt, &user.LoginRequest{Username: "admin", Password: "admin"})
		if err != nil {
			return
		}
		log.Println("->", login.Token, login.User.Id, login.User.Name)
	}(eb.Config.GrpcServer.Port)

	eb.Run()
}
