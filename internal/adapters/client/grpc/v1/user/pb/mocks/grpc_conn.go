package mocks

import (
	"context"
	"log"
	"net"

	user "github.com/bagus-aulia/inventhier/internal/adapters/client/grpc/v1/user/pb"
	"github.com/bxcodec/faker"
	mock "github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

type MockUserServer struct {
	user.UnimplementedUserV1HandlerServer
	mock.Mock
}

func (_m *MockUserServer) GetUser(ctx context.Context, req *user.Params) (*user.Data, error) {
	ret := _m.Called(ctx, req)

	var r0 *user.Data
	if rf, ok := ret.Get(0).(func(ctx context.Context, req *user.Params) *user.Data); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.Data)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(ctx context.Context, req *user.Params) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func Dialer(setup func() *MockUserServer) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()

	user.RegisterUserV1HandlerServer(server, setup())

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func MockGrpcConn() (*grpc.ClientConn, error) {
	setup := func() *MockUserServer {
		var mockUser *user.User
		faker.FakeData(&mockUser)

		mockServer := new(MockUserServer)
		mockServer.On("GetUser", mock.Anything, mock.Anything).Return(&user.Data{
			Status: 200,
			Data:   mockUser,
			Error:  nil,
		}, nil)

		return mockServer
	}

	conn, err := grpc.NewClient(
		"bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(Dialer(setup)),
	)
	if err != nil {
		log.Fatal(err)
	}

	return conn, err
}

func MockGrpcConnError() (*grpc.ClientConn, error) {
	setup := func() *MockUserServer {
		mockServer := new(MockUserServer)
		mockServer.On("GetUser", mock.Anything, mock.Anything).Return(&user.Data{
			Status: 500,
			Data:   nil,
			Error: &user.Error{
				Message: "Internal Server Error",
				Reason:  "internal_server_error",
			},
		}, nil)

		return mockServer
	}

	conn, err := grpc.NewClient(
		"bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(Dialer(setup)),
	)
	if err != nil {
		log.Fatal(err)
	}

	return conn, err
}
