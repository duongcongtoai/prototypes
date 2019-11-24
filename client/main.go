package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	pb "msc/user/proto"

	mscot "msc/opentracing"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
)

var (
	port int
)

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("micro_service")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	port = viper.GetInt("server_port")
}

func main() {
	e := echo.New()
	zipkinMw, opentracingTracer, err := mscot.NewEchoZipkinMiddleWare("client service", "primary span")
	if err != nil {
		log.Fatal(err)
	}
	// opentracingTracer := zipkintracer.Wrap(unwrappedTracer)
	opentracing.SetGlobalTracer(opentracingTracer)

	e.Use(zipkinMw)

	conn, err := mscot.TracedGrpcConn(opentracingTracer, port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	crudClient := pb.NewUserCrudClient(conn)

	e.GET("/create/:email", func(c echo.Context) error {
		email := c.Param("email")
		//context need to have span of parent span
		user, err := crudClient.CreateReturn(c.Request().Context(), &pb.User{
			Name:     "admin",
			Email:    email + "@gmail.com",
			Password: "admin",
			Status:   &pb.User_IsNew{true},
		})

		if err != nil {
			log.Fatal(err)
			return err
		}
		return c.JSON(http.StatusOK, user)
	})
	e.GET("/get/:id", func(c echo.Context) error {
		ctx := c.Request().Context()

		stringId := c.Param("id")
		id, err := strconv.Atoi(stringId)
		if err != nil {
			c.JSON(http.StatusNotFound, err)
			// return err
		}
		result, err := crudClient.Get(ctx, &pb.User{
			Status: &pb.User_Id{int64(id)},
		})
		if err != nil {

			return err
		}
		return c.JSON(http.StatusOK, result)
	})
	e.POST("/delete/:id", func(c echo.Context) error {
		stringId := c.Param("id")
		id, err := strconv.Atoi(stringId)
		if err != nil {
			return err
		}
		result, err := crudClient.Delete(context.Background(), &pb.User{
			Status: &pb.User_Id{int64(id)},
		})
		if err != nil {
			// log.Fatal(err)
			return err
		}
		return c.JSON(http.StatusOK, result)
	})
	e.Start(":8080")
}
