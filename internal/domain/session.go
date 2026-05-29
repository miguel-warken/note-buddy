package domain

type Session struct {
	ID int
	Subject string
	Title string
	StartedAt string
	Status string
	Entries []SessionEntry
}