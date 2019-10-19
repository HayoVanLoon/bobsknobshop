package common

type ServiceImplementationDetails interface {
	ParentName() string
	ParentVersion() string
	Name() string
	Port() int
	Active() bool
}

type serviceImplementationDetails struct {
	parentName    string
	parentVersion string
	name          string
	port          int
	active        bool
}

func (s serviceImplementationDetails) ParentName() string {
	return s.parentName
}

func (s serviceImplementationDetails) ParentVersion() string {
	return s.parentVersion
}

func (s serviceImplementationDetails) Name() string {
	return s.name
}

func (s serviceImplementationDetails) Port() int {
	return s.port
}

func (s serviceImplementationDetails) Active() bool {
	return s.active
}

func NewServiceImplementationDetails(parentName, parentVersion, name string, port int, active bool) ServiceImplementationDetails {
	return serviceImplementationDetails{parentName, parentVersion, name, port, active}
}
