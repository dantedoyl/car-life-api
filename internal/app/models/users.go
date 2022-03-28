package models

type User struct {
	ID uint64
	Name string
	Surname string
	OwnClubs []ClubCard
	ParticipantClubs []ClubCard
	ParticipantEvents []EventCard
}

type UserCard struct {
	ID uint64
	Name string
}

type CarCard struct {
	AvatarUrl string
	Name      string
}
