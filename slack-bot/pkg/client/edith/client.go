package edith

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Shivin01/Edith/slack-bot/pkg/client"
	"github.com/Shivin01/Edith/slack-bot/pkg/config"
	"github.com/pkg/errors"
)

// Client is an interface representing used jenkins functions of gojenkins.
type Client interface {
	GetToken(ctx context.Context, username, password string) (*tokenResponse, error)
	RefreshToken(ctx context.Context, refresh string) (string, error)
	GetUsers(ctx context.Context, token string) ([]UserResponse, error)
	AddUser(ctx context.Context, requestData AddUserRequest) (*AddUserResponse, error)
	ModifyUser(ctx context.Context, userId int, requestData map[string]interface{}, token string) error
	GetUser(ctx context.Context, slackID, token string) ([]*UserResponse, error)
	GetMinimalUser(ctx context.Context, slackID, token string) ([]*UserResponse, error)
	MarkAttendance(ctx context.Context, token string) error
	ListHoliday(ctx context.Context, token string) ([]*HolidayListResponse, error)
	GetNewsFeeds(ctx context.Context, token string) ([]*NewsFeedResponse, error)
	GetCelebrations(ctx context.Context, token string) ([]*celebrationResponse, error)
	RequestForLeave(ctx context.Context, token string, request LeaveRequest) error
	ListLeaves(ctx context.Context, token string) ([]*LeaveResponse, error)
	MakeAnnouncement(ctx context.Context, token string, requestData AnnouncementRequest) (*AnnouncementResponse, error)
	DeleteUser(ctx context.Context, userId int, token string) error
	ListLeavesForApproval(ctx context.Context, token string) ([]*LeaveResponse, error)
	ApproveLeave(ctx context.Context, leaveId int, token string) error
	AddClient(ctx context.Context, requestData AddClientRequest) (*AddClientResponse, error)
	GetClients(ctx context.Context) ([]AddClientResponse, error)
}

// GetClient created Jenkins client with given options/credentials
func GetClient(cfg *config.Server) (Client, error) {
	return createEdithClient(cfg)
}

// implementation of Client interface. proxies to gojenkins with additional handling for inner jenkins jobs.
type edithClientImpl struct {
	client *client.APIClient
}

func createEdithClient(cfg *config.Server) (*edithClientImpl, error) {
	apiClient := client.NewAPIClient(cfg.BaseURL, cfg.MaxConnsPerHost, time.Duration(cfg.Timeout)*time.Second)
	return &edithClientImpl{apiClient}, nil
}

func (e *edithClientImpl) GetToken(ctx context.Context, username, password string) (*tokenResponse, error) {
	body := map[string]string{
		"username": username,
		"password": password,
	}

	res := &tokenResponse{}
	_, err := e.client.Post(ctx, getToken, body, &res, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (e *edithClientImpl) RefreshToken(ctx context.Context, refresh string) (string, error) {
	body := map[string]string{
		"refresh": refresh,
	}

	res := &tokenResponse{}
	_, err := e.client.Post(ctx, refreshToken, body, &res, nil)
	if err != nil {
		return "", err
	}
	return res.Token, nil
}

func (e *edithClientImpl) GetUsers(ctx context.Context, token string) ([]UserResponse, error) {
	var res []UserResponse
	_, err := e.client.GetWithToken(ctx, getUsers, token, &res, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (e *edithClientImpl) AddUser(ctx context.Context, requestData AddUserRequest) (*AddUserResponse, error) {
	res := &AddUserResponse{}
	restyRes, err := e.client.Post(ctx, addUser, requestData, res, nil)
	if err != nil {
		return nil, err
	}
	if restyRes.StatusCode() != http.StatusCreated {
		return nil, errors.New("user is not created, try again")
	}
	return res, nil
}

func (e *edithClientImpl) GetMinimalUser(ctx context.Context, slackID, token string) ([]*UserResponse, error) {
	return e.getUser(ctx, fmt.Sprintf("%s?slack_id=%s", getMinimalUsers, slackID), token)
}

func (e *edithClientImpl) GetUser(ctx context.Context, slackID, token string) ([]*UserResponse, error) {
	return e.getUser(ctx, fmt.Sprintf("%s?slack_id=%s", getUsers, slackID), token)
}

func (e *edithClientImpl) getUser(ctx context.Context, url, token string) ([]*UserResponse, error) {
	var res []*UserResponse
	_, err := e.client.GetWithToken(ctx, url, token, &res, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (e *edithClientImpl) MarkAttendance(ctx context.Context, token string) error {
	body := map[string]string{
		"date": time.Now().Format("2006-01-02"),
	}
	_, err := e.client.PostWithToken(ctx, markAttendance, body, nil, token, nil)
	if err != nil {
		return err
	}
	return nil
}

func (e *edithClientImpl) ListHoliday(ctx context.Context, token string) ([]*HolidayListResponse, error) {
	var holidayList []*HolidayListResponse
	_, err := e.client.GetWithToken(ctx, getHolidayList, token, &holidayList, nil)
	if err != nil {
		return nil, err
	}
	return holidayList, nil
}

func (e *edithClientImpl) GetNewsFeeds(ctx context.Context, token string) ([]*NewsFeedResponse, error) {
	var newsFeeds []*NewsFeedResponse
	_, err := e.client.GetWithToken(ctx, getNewsFeed, token, &newsFeeds, nil)
	if err != nil {
		return nil, err
	}
	return newsFeeds, nil
}

func (e *edithClientImpl) GetCelebrations(ctx context.Context, token string) ([]*celebrationResponse, error) {
	var celebrations []*celebrationResponse
	_, err := e.client.GetWithToken(ctx, getCelebrations, token, &celebrations, nil)
	if err != nil {
		return nil, err
	}
	return celebrations, nil
}

func (e *edithClientImpl) RequestForLeave(ctx context.Context, token string, request LeaveRequest) error {
	body := map[string]interface{}{
		"start_date_time": request.StartDateTime,
		"stop_date_time":  request.StopDateTime,
		"kind":            request.Kind,
		"type":            request.LeaveType,
	}
	_, err := e.client.PostWithToken(ctx, requestLeave, body, nil, token, nil)
	if err != nil {
		return err
	}
	return nil
}

func (e *edithClientImpl) ListLeaves(ctx context.Context, token string) ([]*LeaveResponse, error) {
	var leaves []*LeaveResponse
	_, err := e.client.GetWithToken(ctx, listLeaves, token, &leaves, nil)
	if err != nil {
		return nil, err
	}
	return leaves, nil
}

func (e *edithClientImpl) MakeAnnouncement(ctx context.Context, token string, requestData AnnouncementRequest) (*AnnouncementResponse, error) {
	res := &AnnouncementResponse{}
	body := map[string]string{
		"type":   requestData.Type,
		"detail": requestData.Detail,
	}
	_, err := e.client.PostWithToken(ctx, makeAnnouncement, body, res, token, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (e *edithClientImpl) DeleteUser(ctx context.Context, userID int, token string) error {
	_, err := e.client.Delete(ctx, fmt.Sprintf("%s%d/", deleteEmployee, userID), token, nil)
	if err != nil {
		return err
	}
	return nil
}

func (e *edithClientImpl) ListLeavesForApproval(ctx context.Context, token string) ([]*LeaveResponse, error) {
	var leaves []*LeaveResponse
	_, err := e.client.GetWithToken(ctx, adminLeave, token, &leaves, nil)
	if err != nil {
		return nil, err
	}
	return leaves, nil
}

func (e *edithClientImpl) ApproveLeave(ctx context.Context, leaveId int, token string) error {
	_, err := e.client.Patch(ctx, fmt.Sprintf("%s%d/", adminLeave, leaveId), token, nil, nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (e *edithClientImpl) AddClient(ctx context.Context, requestData AddClientRequest) (*AddClientResponse, error) {
	res := &AddClientResponse{}
	_, err := e.client.PostWithToken(ctx, clientList, requestData, res, "", nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (e *edithClientImpl) GetClients(ctx context.Context) ([]AddClientResponse, error) {
	var res []AddClientResponse
	_, err := e.client.GetWithToken(ctx, clientList, "", &res, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (e *edithClientImpl) ModifyUser(ctx context.Context, userId int, requestData map[string]interface{}, token string) error {
	_, err := e.client.Patch(ctx, fmt.Sprintf("%s%d/", modifyEmployee, userId), token, requestData, nil, nil)
	if err != nil {
		return err
	}
	return nil
}
