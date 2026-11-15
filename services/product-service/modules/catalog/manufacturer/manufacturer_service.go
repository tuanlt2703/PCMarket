package manufacturer

type ManufacturerService struct {
	repo *ManufacturerRepository
}

func NewManufacturerService(repo *ManufacturerRepository) *ManufacturerService {
	return &ManufacturerService{repo: repo}
}

func (svc *ManufacturerService) GetAll() ([]Manufacturer, error) {
	return svc.repo.GetAll()
}

func (svc *ManufacturerService) GetByID(id int) (Manufacturer, error) {
	return svc.repo.GetByID(id)
}

func (svc *ManufacturerService) Create(newComponent CreateManufacturer) (Manufacturer, error) {
	return svc.repo.Create(newComponent)
}

func (svc *ManufacturerService) Update(id int, updateComponent UpdateManufacturer) (Manufacturer, error) {
	return svc.repo.Update(id, updateComponent)
}

func (svc *ManufacturerService) Delete(id int) error {
	return svc.repo.Delete(id)
}
