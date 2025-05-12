package port

type PipelineService interface {
	Run() error
	Status() error
	Stop() error
}
