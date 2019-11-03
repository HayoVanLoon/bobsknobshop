package v1

import (
	"google.golang.org/grpc"
	"log"
)

type ServiceImplementation interface {
	ParentName() string
	ParentVersion() string
	Name() string
	Port() int
	Active() bool
	Service() string
}

type serviceImplementation struct {
	parentName    string
	parentVersion string
	name          string
	port          int
	active        bool
}

func (s serviceImplementation) ParentName() string {
	return s.parentName
}

func (s serviceImplementation) ParentVersion() string {
	return s.parentVersion
}

func (s serviceImplementation) Name() string {
	return s.name
}

func (s serviceImplementation) Port() int {
	return s.port
}

func (s serviceImplementation) Active() bool {
	return s.active
}

func (s serviceImplementation) Service() string {
	return s.parentName + "-" + s.parentVersion + "-service-" + s.name
}

func NewServiceImplementation(parentName, parentVersion, name string, port int, active bool) ServiceImplementation {
	return serviceImplementation{parentName, parentVersion, name, port, active}
}

// Provides a connection-closing function
func CloseConnFn(conn *grpc.ClientConn) func() {
	return func() {
		if err := conn.Close(); err != nil {
			log.Printf("WARN: error closing connection: %v", err)
		}
	}
}

// Establishes a connection to a service
func GetConn(s string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(s, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
