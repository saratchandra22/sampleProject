package services

import (
	"github.com/saratchandra13/sampleProject/pkg/domain/entity"
)

type beerList []*entity.Beer

func (al *appLogic) ListBeer() (beerList, error) {
	beerList, err := al.beerRepo.GetAllBeer()
	if err != nil {
		return nil, err
	}
	return beerList, nil
}
