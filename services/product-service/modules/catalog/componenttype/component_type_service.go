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
