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

func (cu *ClubsUsecase) GetClubByID(id uint64, userID uint64) (*models.Club, error) {
	return cu.clubsRepo.GetClubByID(int64(id), userID)
}

func (cu *ClubsUsecase) GetClubs(idGt *uint64, idLte *uint64, limit *uint64, query *string) ([]*models.Club, error) {
	return cu.clubsRepo.GetClubs(idGt, idLte, limit, query)
}

func (cu *ClubsUsecase) UpdateAvatar(clubID int64, fileHeader *multipart.FileHeader) (*models.Club, error) {
	club, err := cu.clubsRepo.GetClubByID(clubID, 0)
	if err != nil {
		return nil, err
	}

	imgUrl, err := filesystem.InsertPhoto(fileHeader, "img/clubs/")
	if err != nil {
		return nil, err
	}

	oldAvatar := club.AvatarUrl
	club.AvatarUrl = imgUrl
	club, err = cu.clubsRepo.UpdateClub(club)
	if err != nil {
		return nil, err
	}

	if oldAvatar == "/img/clubs/default.webp" {
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

func (cu *ClubsUsecase) GetClubsUserByStatus(club_id int64, status string, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.UserCard, error) {
	return cu.clubsRepo.GetClubsUserByStatus(club_id, status, idGt, idLte, limit)
}

func (cu *ClubsUsecase) GetClubsCars(club_id int64, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.CarCard, error) {
	return cu.clubsRepo.GetClubsCars(club_id, idGt, idLte, limit)
}

func (cu *ClubsUsecase) GetClubsEvents(club_id int64, idGt *uint64, idLte *uint64, limit *uint64) ([]*models.EventCard, error) {
	return cu.clubsRepo.GetClubsEvents(club_id, idGt, idLte, limit)
}

func (cu *ClubsUsecase) SetUserStatusByClubID(clubID int64, userID int64, status string) error {
	return cu.clubsRepo.SetUserStatusByClubID(clubID, userID, status)
}

func (cu *ClubsUsecase) ApproveRejectUserParticipateInClub(clubID int64, userID int64, decision string) error {
	if decision == "approve" {
		return cu.clubsRepo.SetUserStatusByClubID(clubID, userID, "participant")
	}
	return cu.clubsRepo.SetUserStatusByClubID(clubID, userID, "subscriber")
}

func (cu *ClubsUsecase) GetUserStatusInClub(clubID int64, userID int64) (*models.ClubUser, error) {
	return cu.clubsRepo.GetUserStatusInClub(clubID, userID)
}

func (cu *ClubsUsecase) SetClubChatID(clubID int64, chatID int64) error {
	return cu.clubsRepo.SetClubChatID(clubID, chatID)
}

func (cu *ClubsUsecase) GetClubChatID(clubID int64, userID int64) (int64, error) {
	return cu.clubsRepo.GetClubChatID(clubID, userID)
}

func (cu *ClubsUsecase) DeleteUserFromClub(clubID int64, userID int64) error {
	return cu.clubsRepo.DeleteUserFromClub(clubID, userID)
}

func (cu *ClubsUsecase) DeleteClubByID(clubID int64) error {
	return cu.clubsRepo.DeleteClubByID(clubID)
}

func (cu *ClubsUsecase) ComplainByID(complaint models.Complaint) error {
	return cu.clubsRepo.ComplainByID(complaint)
}
