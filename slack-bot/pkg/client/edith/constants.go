package edith

import (
	"fmt"
	"strings"
)

const (
	getToken         = "/api/token/"
	refreshToken     = "/api/token/refresh/"
	getUsers         = "/employees/info/"
	getMinimalUsers  = "/employees/minimal_info/"
	addUser          = "/rest-auth/registration/"
	markAttendance   = "/employees/employee_attendance/"
	getHolidayList   = "/client/holiday_list/"
	getNewsFeed      = "/announcements/news_feed/"
	getCelebrations  = "/announcements/celebration/"
	requestLeave     = "/employees/leave/"
	listLeaves       = "/employees/leave/"
	makeAnnouncement = "/announcements/announcement/"
	deleteEmployee   = "/employees/admin/"
	adminLeave       = "/employees/admin_leave/"
)

type tokenResponse struct {
	Token string `json:"token"`
}

type IUser interface {
	GetSlackID() string
	GetUsername() string
	GetRealName() string
	GetDesignation() string
}

type UserResponse struct {
	SlackID     string   `json:"slack_id"`
	Username    string   `json:"username"`
	FirstName   string   `json:"first_name"`
	MiddleName  string   `json:"middle_name"`
	LastName    string   `json:"last_name"`
	Designation string   `json:"designation"`
	PhoneNumber string   `json:"phone_number"`
	Email       string   `json:"email"`
	Image       string   `json:"image"`
	Gender      string   `json:"gender"`
	Skills      []string `json:"skills"`
	JoinDate    string   `json:"join_date"`
	Departments []string `json:"departments"`
}

func (u *UserResponse) GetSlackID() string {
	return u.SlackID
}

func (u *UserResponse) GetUsername() string {
	return u.Username
}

func (u *UserResponse) GetRealName() string {
	name := strings.Trim(u.FirstName, " ")
	if u.MiddleName != "" {
		name += fmt.Sprintf(" %s", u.MiddleName)
	}
	if u.LastName != "" {
		name += fmt.Sprintf(" %s", u.LastName)
	}
	return name
}

func (u *UserResponse) GetDesignation() string {
	return u.Designation
}

type AddUserRequest struct {
	Username    string
	FirstName   string
	MiddleName  string
	LastName    string
	Password    string
	Skills      []string
	PhoneNumber string
	Email       string
	SlackID     string
}

type AddUserResponse struct {
	Token string `json:"token"`
	User  struct {
		Pk        int    `json:"pk"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	} `json:"user"`
}

type HolidayListResponse struct {
	Date        string `json:"date"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type NewsFeedResponse struct {
	Image       string   `json:"image"`
	FirstName   string   `json:"first_name"`
	MiddleName  string   `json:"middle_name"`
	LastName    string   `json:"last_name"`
	Skills      []string `json:"skills"`
	Designation string   `json:"designation"`
	ID          int      `json:"id"`
	JoiningDate string   `json:"joining_date"`
}

type celebrationResponse struct {
	Image      string `json:"image"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	ID         int    `json:"id"`
}

type LeaveRequest struct {
	StartDateTime int64  `json:"start_date_time"`
	StopDateTime  int64  `json:"stop_date_time"`
	Kind          string `json:"kind"`
	LeaveType     string `json:"type"`
}

type defaultResponse struct {
	ID        int   `json:"id"`
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type LeaveResponse struct {
	defaultResponse
	StartDateTime float64 `json:"start_date_time"`
	StopDateTime  float64 `json:"stop_date_time"`
	Kind          string  `json:"kind"`
	LeaveType     string  `json:"type"`
	Employee      struct {
		FirstName  string `json:"first_name"`
		MiddleName string `json:"middle_name"`
		LastName   string `json:"last_name"`
	} `json:"employee"`
	ApprovedBy struct {
		FirstName  string `json:"first_name"`
		MiddleName string `json:"middle_name"`
		LastName   string `json:"last_name"`
	} `json:"approved_by"`
}

type AnnouncementRequest struct {
	Type   string `json:"type"`
	Detail string `json:"detail"`
}

type AnnouncementResponse struct {
	defaultResponse
	AnnouncementRequest
	Employee struct {
		FirstName  string `json:"first_name"`
		MiddleName string `json:"middle_name"`
		LastName   string `json:"last_name"`
	} `json:"employee"`
	Client int `json:"client"`
}
