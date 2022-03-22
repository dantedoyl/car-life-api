package usecase

import (
	"github.com/dantedoyl/car-life-api/internal/app/clients/filesystem"
	clubs "github.com/dantedoyl/car-life-api/internal/app/clubs"
	"github.com/dantedoyl/car-life-api/internal/app/models"
	"mime/multipart"
)

type ClubsUsecase struct {
	clubsRepo clubs.IClubsRepository
}

func NewClubsUsecase(repo clubs.IClubsRepository) clubs.IClubsUsecase {
	return &ClubsUsecase{
		clubsRepo: repo,
	}
}

func (cu *ClubsUsecase) CreateClub(club *models.Club) error {
	return cu.clubsRepo.InsertClub(club)
}

func (cu *ClubsUsecase) GetClubByID(id uint64) (*models.Club, error) {
	return cu.clubsRepo.GetClubByID(int64(id))
}

func (cu *ClubsUsecase) GetClubs(idGt  *uint64, idLte *uint64, limit *uint64, query *string) ([]*models.Club, error) {
	return cu.clubsRepo.GetClubs(idGt, idLte, limit, query)
}

func (cu *ClubsUsecase) UpdateAvatar(clubID int64, fileHeader *multipart.FileHeader) (*models.Club, error) {
	club, err := cu.clubsRepo.GetClubByID(clubID)
	if err != nil {
		return nil, err
	}

	imgUrl, err := filesystem.InsertPhoto(fileHeader, "static/avatar/")
	if err != nil {
		return nil, err
	}

	oldAvatar := club.AvatarUrl
	club.AvatarUrl = imgUrl
	club, err = cu.clubsRepo.UpdateClub(club)
	if err != nil {
		return nil, err
	}

	if oldAvatar == "" {
		return club, nil
	}

	err = filesystem.RemovePhoto(oldAvatar)
	if err != nil {
		return nil, err
	}

	return club, nil
}

func (cu *ClubsUsecase) GetTags() ([]models.Tag, error) {
	return cu.clubsRepo.GetTags()
}