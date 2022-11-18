package test_http

import "time"

type GetDeviceResp struct {
	Value struct {
		Id        string    `json:"id"`
		Platform  string    `json:"platform"`
		UserId    string    `json:"userId"`
		EnteredAt time.Time `json:"enteredAt"`
	} `json:"value"`
}

type ListDevicesResp struct {
	Items []struct {
		Id        string    `json:"id"`
		Platform  string    `json:"platform"`
		UserId    string    `json:"userId"`
		EnteredAt time.Time `json:"enteredAt"`
	} `json:"items"`
}

type CreateDeviceReq struct {
	Platform string `json:"platform"`
	UserId   string `json:"userId"`
}

type CreateDeviceResp struct {
	DeviceId string `json:"deviceId"`
}

type UpdateDeviceReq struct {
	Platform string `json:"platform"`
	UserId   string `json:"userId"`
}

type UpdateDeviceResp struct {
	Success bool `json:"success"`
}

type RemoveDeviceResp struct {
	Found bool `json:"found"`
}
