package componenttype

type ComponentTypeService struct {
	repo *ComponentTypeRepository
}

func NewComponentTypeService(repo *ComponentTypeRepository) *ComponentTypeService {
	return &ComponentTypeService{repo: repo}
}

func (svc *ComponentTypeService) GetAll() ([]ComponentType, error) {
	return svc.repo.GetAll()
}

func (svc *ComponentTypeService) GetByID(id int) (ComponentType, error) {
	return svc.repo.GetByID(id)
}

func (svc *ComponentTypeService) Create(newComponent CreateComponentType) (ComponentType, error) {
	return svc.repo.Create(newComponent)
}

func (svc *ComponentTypeService) Update(id int, updateComponent UpdateComponentType) (ComponentType, error) {
	return svc.repo.Update(id, updateComponent)
}

func (svc *ComponentTypeService) Delete(id int) error {
	return svc.repo.Delete(id)
}
