package test_http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/pkg/logger"
	"net/http"
	"net/http/httputil"
	"time"
)

const baseURL = "http://localhost:8080"
const actDeviceURL = "/api/v1/devices"
const actNotifURL = "/api/v1/notification"

// Client for the Act Device API
type Client interface {
	Do(req *http.Request) (*http.Response, error)
	// Device methods
	CreateDevice(platform string, userId string) (*http.Response, error)
	GetDevice(deviceId string) (*http.Response, error)
	ListDevices(page string, perPage string) (*http.Response, error)
	UpdateDevice(deviceId string, platform string, userId string) (*http.Response, error)
	RemoveDevice(deviceId string) (*http.Response, error)
	// Notification methods
	SendNotification(deviceId string, message string) (*http.Response, error)
	GetNotification(deviceId string) (*http.Response, error)
	AckNotification(notifId string) (*http.Response, error)
	SubscribeNotification(deviceId string) (*http.Response, error)
}

type client struct {
	client  *retryablehttp.Client
	baseUrl string
}

// NewHTTPClient creates a new HTTP client
func NewHTTPClient(retryMax int, timeout time.Duration) Client {
	c := &retryablehttp.Client{
		HTTPClient:      &http.Client{Timeout: timeout},
		RetryMax:        retryMax,
		RetryWaitMin:    1 * time.Second,
		RetryWaitMax:    10 * time.Second,
		CheckRetry:      retryablehttp.DefaultRetryPolicy,
		Backoff:         retryablehttp.DefaultBackoff,
		RequestLogHook:  requestHook,
		ResponseLogHook: responseHook,
	}

	client := &client{client: c, baseUrl: baseURL}
	return client
}

func requestHook(_ retryablehttp.Logger, req *http.Request, retry int) {
	dump, err := httputil.DumpRequest(req, true) // better way
	if err != nil {
		logger.ErrorKV(req.Context(), "can't dump request")
	}

	logger.InfoKV(
		req.Context(),
		fmt.Sprintf("Retry request %d", retry),
		"request", string(dump),
		"url", req.URL.String(),
	)
}

func responseHook(_ retryablehttp.Logger, res *http.Response) {
	dump, err := httputil.DumpResponse(res, true) // better way
	if err != nil {
		logger.ErrorKV(res.Request.Context(), "can't dump response")
	}

	logger.InfoKV(
		res.Request.Context(),
		"Responded",
		"response", string(dump),
		"url", res.Request.URL.String(),
		"status_code", res.StatusCode,
	)
}

func (c *client) Do(request *http.Request) (*http.Response, error) {
	req, err := retryablehttp.FromRequest(request)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}

func (c *client) CreateDevice(platform string, userId string) (*http.Response, error) {
	data, err := json.Marshal(CreateDeviceReq{
		Platform: platform,
		UserId:   userId,
	})
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(data)
	resp, err := c.client.Post(c.baseUrl+actDeviceURL, "application/json", body)
	return resp, err
}

func (c *client) GetDevice(deviceId string) (*http.Response, error) {
	resp, err := c.client.Get(baseURL + actDeviceURL + "/" + deviceId)
	return resp, err
}

func (c *client) ListDevices(page string, perPage string) (*http.Response, error) {
	resp, err := c.client.Get(baseURL + actDeviceURL + "?page=" + page + "&perPage=" + perPage)
	return resp, err
}

func (c *client) UpdateDevice(deviceId string, platform string, userId string) (*http.Response, error) {
	data, err := json.Marshal(UpdateDeviceReq{
		Platform: platform,
		UserId:   userId,
	})
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPut, baseURL+actDeviceURL+"/"+deviceId, body)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	return resp, err
}

func (c *client) RemoveDevice(deviceId string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, baseURL+actDeviceURL+"/"+deviceId, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	return resp, err
}

func (c *client) SendNotification(deviceId string, message string) (*http.Response, error) {
	data, err := json.Marshal(SendNotificationReq{
		NotificationReq{
			DeviceID:           deviceId,
			Username:           "1",
			Message:            message,
			Lang:               "LANG_ENGLISH",
			NotificationStatus: "STATUS_CREATED",
		},
	})
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(data)
	resp, err := c.client.Post(c.baseUrl+actNotifURL, "application/json", body)
	return resp, err
}

func (c *client) GetNotification(deviceId string) (*http.Response, error) {
	resp, err := c.client.Get(baseURL + actNotifURL + "?deviceId=" + deviceId)
	return resp, err
}

func (c *client) AckNotification(notifId string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, baseURL+actNotifURL+"/ack/"+notifId, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	return resp, err
}

func (c *client) SubscribeNotification(deviceId string) (*http.Response, error) {
	return nil, nil
}
