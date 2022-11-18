package test_http

type GetNotificationResp struct {
	Notification []NotificationResp `json:"notification"`
}
type NotificationResp struct {
	NotificationId string `json:"notificationId"`
	Message        string `json:"message"`
}

type SendNotificationReq struct {
	Notification NotificationReq `json:"notification"`
}
type NotificationReq struct {
	//NotificationID     string `json:"notificationId"`
	DeviceID           string `json:"deviceId"`
	Username           string `json:"username"`
	Message            string `json:"message"`
	Lang               string `json:"lang"`
	NotificationStatus string `json:"notificationStatus"`
}

type SendNotificationResp struct {
	NotificationId string `json:"notificationId"`
}

type AckNotificationResp struct {
	Success bool `json:"success"`
}

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
