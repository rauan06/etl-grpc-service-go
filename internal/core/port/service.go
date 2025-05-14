package port

type Service interface {
	Run()
	GetServiceName() string
	Status() int
	Stop()
}
