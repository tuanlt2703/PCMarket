package model

type ModelService struct {
	repo *ModelRepository
}

func NewModelService(repo *ModelRepository) *ModelService {
	return &ModelService{repo: repo}
}

func (svc *ModelService) GetAll() ([]Model, error) {
	return svc.repo.GetAll()
}

func (svc *ModelService) GetByID(id int) (Model, error) {
	return svc.repo.GetByID(id)
}

func (svc *ModelService) Create(createModel CreateModel) (Model, error) {
	return svc.repo.Create(createModel)
}

func (svc *ModelService) Update(id int, updateModel UpdateModel) (Model, error) {
	return svc.repo.Update(id, updateModel)
}

func (svc *ModelService) Delete(id int) error {
	return svc.repo.Delete(id)
}
