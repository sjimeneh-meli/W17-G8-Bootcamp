package repositories

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/error_message"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/models"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
)

type SellerRepository interface {
	GetAll() (map[int]models.Seller, error)
	Save(seller models.Seller) (models.Seller, error)
	Update(id int, seller models.Seller) (models.Seller, error)
	Delete(id int) error
}

type JsonSellerRepository struct {
	storage loader.Storage[models.Seller]
}

func NewJSONSellerRepository(storage loader.Storage[models.Seller]) *JsonSellerRepository {

	return &JsonSellerRepository{
		storage: storage,
	}
}

func (r *JsonSellerRepository) GetAll() (map[int]models.Seller, error) {
	data, err := r.storage.ReadAll()
	//fmt.Println(data)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (r *JsonSellerRepository) Save(seller models.Seller) (models.Seller, error) {
	sellers, err := r.GetAll()
	if err != nil {
		return models.Seller{}, err
	}
	var id int
	for _, value := range sellers {
		if value.Id > id {
			id = value.Id
		}
		if value.CID == seller.CID {
			return models.Seller{}, error_message.ErrAlreadyExists
		}
	}
	seller.Id = id + 1
	sellers[seller.Id] = seller

	itemSlice := loader.NewJSONStorage[models.Seller]("").MapToSlice(sellers)

	err = r.storage.WriteAll(itemSlice)

	if err != nil {
		return models.Seller{}, err
	}

	return seller, nil
}

func (r *JsonSellerRepository) Update(id int, seller models.Seller) (models.Seller, error) {
	sellers, err := r.GetAll()
	if err != nil {
		return models.Seller{}, err
	}
	_, ok := sellers[id]

	if !ok {
		return models.Seller{}, error_message.ErrNotFound
	}

	existingSeller := sellers[id]

	if seller.CID != "" {
		existingSeller.CID = seller.CID
	}
	if seller.CompanyName != "" {
		existingSeller.CompanyName = seller.CompanyName
	}
	if seller.Address != "" {
		existingSeller.Address = seller.Address
	}
	if seller.Telephone != "" {
		existingSeller.Telephone = seller.Telephone
	}

	sellers[id] = existingSeller

	itemSlice := loader.NewJSONStorage[models.Seller]("").MapToSlice(sellers)

	err = r.storage.WriteAll(itemSlice)

	if err != nil {
		return models.Seller{}, err
	}

	return existingSeller, nil

}

func (r *JsonSellerRepository) Delete(id int) error {
	sellers, err := r.GetAll()
	if err != nil {
		return err
	}

	_, ok := sellers[id]

	if !ok {
		return error_message.ErrNotFound
	}

	delete(sellers, id)

	itemSlice := loader.NewJSONStorage[models.Seller]("").MapToSlice(sellers)

	errWrite := r.storage.WriteAll(itemSlice)

	if errWrite != nil {
		return errWrite
	}

	return nil
}
