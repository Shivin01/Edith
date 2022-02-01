package db

import (
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	ServerID          int    `gorm:"server_id"`
	Name              string `gorm:"name"`
	RegisteredName    string `gorm:"registered_name"`
	LeaveCount        int    `gorm:"leave_count"`
	NoticePeriodCount int    `gorm:"notice_period_count"`
}

type User struct {
	ID              string           `gorm:"id:primaryKey"`
	ServerID        int              `gorm:"server_id"`
	Username        string           `gorm:"username"`
	FullName        string           `gorm:"full_name"`
	AccessToken     string           `gorm:"access_token"`
	Designation     string           `gorm:"designation"`
	CustomVariables []CustomVariable `gorm:"foreignKey:UserRefer"`
	CustomCommands  []CustomCommand  `gorm:"foreignKey:UserRefer"`
	ClientName      string           `gorm:"client_name"`
	Gender          string           `gorm:"gender"`
}

func (u *User) GetSlackID() string {
	return u.ID
}

func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) GetRealName() string {
	return u.FullName
}

func (u *User) GetDesignation() string {
	return u.Designation
}

func (u *User) IsAdmin() bool {
	return u.Designation == "hr" || u.Designation == "manager" || u.Designation == "admin"
}

type CustomVariable struct {
	gorm.Model
	Name      string `gorm:"name"`
	Value     string `gorm:"value"`
	UserRefer string
}

type CustomCommand struct {
	gorm.Model
	Alias     string `gorm:"alias"`
	Command   string `gorm:"command"`
	UserRefer string
}

type Stat struct {
	gorm.Model
	TotalCommands        int `gorm:"total_commands"`
	UnauthorizedCommands int `gorm:"unauthorized_commands"`
	UnknownCommands      int `gorm:"unknown_commands"`
	Interactions         int `gorm:"interactions"`
}

type FallbackQueue struct {
	gorm.Model
	Channel         string `gorm:"channel"`
	User            string `gorm:"user"`
	Timestamp       string `gorm:"timestamp"`
	Thread          string `gorm:"thread"`
	InternalMessage bool   `gorm:"internal_message"`
	UpdatedMessage  bool   `gorm:"updated_message"`
	QueueKey        string `gorm:"queue_key"`
	Text            string `gorm:"text"`
}
