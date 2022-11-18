package test_http

import "time"

// GetDeviceResp creates a struct for a response
type GetDeviceResp struct {
	Value struct {
		Id        string    `json:"id"`
		Platform  string    `json:"platform"`
		UserId    string    `json:"userId"`
		EnteredAt time.Time `json:"enteredAt"`
	} `json:"value"`
}

// ListDevicesResp creates a struct for a response
type ListDevicesResp struct {
	Items []struct {
		Id        string    `json:"id"`
		Platform  string    `json:"platform"`
		UserId    string    `json:"userId"`
		EnteredAt time.Time `json:"enteredAt"`
	} `json:"items"`
}

// CreateDeviceReq creates a struct for a request
type CreateDeviceReq struct {
	Platform string `json:"platform"`
	UserId   string `json:"userId"`
}

// CreateDeviceResp creates a struct for a response
type CreateDeviceResp struct {
	DeviceId string `json:"deviceId"`
}

// UpdateDeviceReq creates a struct for a request
type UpdateDeviceReq struct {
	Platform string `json:"platform"`
	UserId   string `json:"userId"`
}

// UpdateDeviceResp creates a struct for a response
type UpdateDeviceResp struct {
	Success bool `json:"success"`
}

// RemoveDeviceResp creates a struct for a response
type RemoveDeviceResp struct {
	Found bool `json:"found"`
}
