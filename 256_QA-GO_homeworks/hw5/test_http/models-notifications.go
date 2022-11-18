package test_http

// GetNotificationResp creates a struct for a response
type GetNotificationResp struct {
	Notification []NotificationResp `json:"notification"`
}

// NotificationResp creates a struct for Notification in GetNotificationResp
type NotificationResp struct {
	NotificationId string `json:"notificationId"`
	Message        string `json:"message"`
}

// SendNotificationReq creates a struct for a request
type SendNotificationReq struct {
	Notification NotificationReq `json:"notification"`
}

// NotificationReq creates a struct for a request
type NotificationReq struct {
	//NotificationID     string `json:"notificationId"`
	DeviceID           string `json:"deviceId"`
	Username           string `json:"username"`
	Message            string `json:"message"`
	Lang               string `json:"lang"`
	NotificationStatus string `json:"notificationStatus"`
}

// SendNotificationResp creates a struct for a response
type SendNotificationResp struct {
	NotificationId string `json:"notificationId"`
}

// AckNotificationResp creates a struct for a response
type AckNotificationResp struct {
	Success bool `json:"success"`
}

// SubscribeNotificationResp creates a struct for a response
type SubscribeNotificationResp struct {
	Result struct {
		NotificationId string `json:"notificationId"`
		Message        string `json:"message"`
	} `json:"result"`
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Details []struct {
			TypeUrl string `json:"typeUrl"`
			Value   string `json:"value"`
		} `json:"details"`
	} `json:"error"`
}
