package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"msc/db"
	"net"

	"github.com/spf13/viper"

	userPb "msc/user/proto"

	"msc/user"

	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"

	mscot "msc/opentracing"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

type userCrudServer struct {
	db *gorm.DB
}

func toUserProtobuf(u *user.User) (userPb.User, error) {
	user := userPb.User{}
	user.Status = &userPb.User_Id{Id: int64(u.ID)}
	user.Name = u.Name
	user.Email = u.Email
	return user, nil
}

func fromUserProtobuf(pbu *userPb.User) (user.User, error) {
	newUser := user.User{}
	password, err := processPassword([]byte(pbu.Password))
	if err != nil {
		return newUser, err
	}
	newUser.Password = password
	newUser.Name = pbu.Name
	newUser.Email = pbu.Email
	return newUser, nil
}

func (u userCrudServer) Create(ctx context.Context, pbu *userPb.User) (*userPb.Result, error) {
	newUser, err := fromUserProtobuf(pbu)
	if err != nil {
		return &userPb.Result{Value: false}, errors.New("Error while encoding password")
	}
	u.db.Create(&newUser)
	if u.db.NewRecord(newUser) {
		return &userPb.Result{Value: false}, errors.New("Failed to insert new User")
	}
	return &userPb.Result{Value: true}, nil
}
func (u userCrudServer) Get(ctx context.Context, pbu *userPb.User) (*userPb.Result, error) {
	user, err := fromUserProtobuf(pbu)
	if err != nil {
		return &userPb.Result{Value: false}, err
	}
	u.db.Find(&user)
	result, err := toUserProtobuf(&user)
	if err != nil {
		return &userPb.Result{Value: false}, err
	}
	return &userPb.Result{Value: true, User: &result}, nil

}
func (u userCrudServer) Delete(ctx context.Context, pbu *userPb.User) (*userPb.Result, error) {
	deleteUser, err := fromUserProtobuf(pbu)
	if err != nil {
		return &userPb.Result{Value: false}, err
	}
	u.db.Delete(&deleteUser)
	if u.db.NewRecord(deleteUser) {
		return &userPb.Result{Value: true}, nil
	}
	return &userPb.Result{Value: false}, errors.New("Faield to delete user")
}

func (u userCrudServer) CreateReturn(ctx context.Context, pbu *userPb.User) (*userPb.Result, error) {
	newUser, err := fromUserProtobuf(pbu)
	if err != nil {
		return &userPb.Result{Value: false}, err
	}
	u.db.Create(&newUser)
	pbu.Status = &userPb.User_Id{int64(newUser.ID)}
	if u.db.NewRecord(newUser) {
		return &userPb.Result{Value: false}, errors.New("Failed to insert new User")
	}
	return &userPb.Result{Value: true, User: pbu}, nil
}

func (u userCrudServer) Edit(ctx context.Context, pbu *userPb.User) (*userPb.Result, error) {
	if !pbu.GetIsNew() {
		return &userPb.Result{Value: false}, errors.New("User passed in edit function is still new")
	}
	editedUser, err := fromUserProtobuf(pbu)
	if err != nil {
		return &userPb.Result{Value: false}, errors.New("Error while encoding password")
	}
	editedUser.ID = uint(pbu.GetId())

	u.db.Save(&editedUser)
	return &userPb.Result{Value: true}, nil
}

func (u userCrudServer) EditReturn(ctx context.Context, pbu *userPb.User) (*userPb.Result, error) {
	if !pbu.GetIsNew() {
		return &userPb.Result{Value: false}, errors.New("User passed in edit function is still new")
	}
	editedUser, err := fromUserProtobuf(pbu)
	if err != nil {
		return &userPb.Result{Value: false}, errors.New("Error while encoding password")
	}
	editedUser.ID = uint(pbu.GetId())

	u.db.Save(&editedUser)
	return &userPb.Result{Value: true, User: pbu}, nil
}

func main() {
	viper.SetConfigFile("micro_service.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	port := viper.GetInt("server_port")
	db, closeFunc, err := db.OpenDb()
	if err != nil {
		//To be adjusted in the future
		log.Fatal(err)
	}
	crudService := userCrudServer{db}

	if err != nil {
		log.Fatal(err)
	}
	t, err := mscot.NewTracer("grpc server service")

	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(zipkingrpc.NewServerHandler(t)),
	)

	userPb.RegisterUserCrudServer(grpcServer, &crudService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer.Serve(lis)
	defer closeFunc()
}

func processPassword(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
