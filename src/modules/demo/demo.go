package demo

// Archivo único generado por Goney.
type DemoResponse struct {
	ID string `json:"id"`
}

type CreateDemoRequest struct {
	Name string `json:"name"`
}

type UpdateDemoRequest struct {
	Name string `json:"name,omitempty"`
}

type Demo struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

// Repository
type DemoRepository struct{}
func NewDemoRepository() *DemoRepository { return &DemoRepository{} }
func (r *DemoRepository) FindAll() ([]Demo, error) { return []Demo{}, nil }
func (r *DemoRepository) FindByID(id string) (*Demo, error) { return &Demo{ID: id}, nil }
func (r *DemoRepository) Create(e *Demo) (*Demo, error) { return e, nil }
func (r *DemoRepository) Update(e *Demo) (*Demo, error) { return e, nil }
func (r *DemoRepository) Delete(id string) error { return nil }

// Service
type DemoService struct { repo *DemoRepository }
func NewDemoService(repo *DemoRepository) *DemoService { return &DemoService{repo: repo} }

// Controller (placeholder)
type DemoController struct { svc *DemoService }
func NewDemoController(svc *DemoService) *DemoController { return &DemoController{svc: svc} }

// Module wiring
type DemoModule struct { Controller *DemoController; Service *DemoService; Repository *DemoRepository }
func NewDemoModule() *DemoModule {
	repo := NewDemoRepository()
	svc := NewDemoService(repo)
	ctrl := NewDemoController(svc)
	return &DemoModule{Controller: ctrl, Service: svc, Repository: repo}
}
