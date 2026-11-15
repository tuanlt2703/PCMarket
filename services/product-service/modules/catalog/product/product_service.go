package product

type ProductService struct {
	repo *ProductRepository
}

func NewProductService(repo *ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (svc *ProductService) GetAll() ([]ProductModel, error) {
	return svc.repo.GetAll()
}

func (svc *ProductService) GetByID(id int) (ProductModel, error) {
	return svc.repo.GetByID(id)
}

func (svc *ProductService) Create(createProductModel CreateProductModel) (ProductModel, error) {
	return svc.repo.Create(createProductModel)
}

func (svc *ProductService) Update(id int, updateProductModel UpdateProductModel) error {
	return svc.repo.Update(id, updateProductModel)
}

func (svc *ProductService) Delete(id int) error {
	return svc.repo.Delete(id)
}
