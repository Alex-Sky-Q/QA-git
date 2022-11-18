//go:build e2e_test_gomock

package test_http

import (
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestDeviceLifecycleMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	clientMock := NewMockClient(ctrl)

	runner.Run(t, "HTTP - Test device lifecycle (CRUD)", func(t provider.T) {
		t.Title("Test device lifecycle (CRUD)")
		t.Description("Basic device lifecycle - create, update, get and remove device")
		platform := "ios"
		userId := "1"
		newPlatform := "android"
		newUserId := "2"
		deviceId := ""

		// Create a device
		t.WithNewStep("Create a device", func(sCtx provider.StepCtx) {
			// Expected mock response
			var mockResp http.Response
			mockResp.StatusCode = 200
			bodyReader := strings.NewReader("{\"deviceId\":\"478\"}")
			mockResp.Body = io.NopCloser(bodyReader)

			clientMock.EXPECT().CreateDevice(platform, userId).Return(&mockResp, nil)

			// Test mock
			resp, err := clientMock.CreateDevice(platform, userId)
			sCtx.Require().NoError(err)
			sCtx.Assert().Equal(http.StatusOK, resp.StatusCode, "Status code should be correct. Want: %v. Got %v",
				http.StatusOK, resp.StatusCode)
			var createDevResp CreateDeviceResp
			err = json.NewDecoder(resp.Body).Decode(&createDevResp)
			sCtx.Require().NoError(err)
			deviceId = createDevResp.DeviceId
			sCtx.Assert().NotEmpty(deviceId, "Device ID is empty")
			sCtx.WithNewAttachment("Response", allure.Text, []byte(fmt.Sprintf("%+v", createDevResp)))
		})

		// Check that device was created
		t.WithNewStep("Check that device was created", func(sCtx provider.StepCtx) {
			// Expected mock response
			var mockResp http.Response
			mockResp.StatusCode = 200
			bodyReader := strings.NewReader("{\"value\":{\"id\":\"478\", \"platform\":\"ios\", \"userId\":\"1\", " +
				"\"enteredAt\":\"2022-11-01T00:36:10.353993Z\"}}")
			mockResp.Body = io.NopCloser(bodyReader)

			clientMock.EXPECT().GetDevice(deviceId).Return(&mockResp, nil)

			// Test mock
			resp, err := clientMock.GetDevice(deviceId)
			sCtx.Require().NoError(err)
			sCtx.Assert().Equal(http.StatusOK, resp.StatusCode, "Status code should be correct. Want: %v. Got %v",
				http.StatusOK, resp.StatusCode)

			deviceResp := new(GetDeviceResp)
			err = json.NewDecoder(resp.Body).Decode(deviceResp)
			sCtx.Require().NoError(err)

			sCtx.Assert().Equal(platform, deviceResp.Value.Platform, "Platform should be correct. Got: %v. Want: %v",
				deviceResp.Value.Platform, platform)
			sCtx.Assert().Equal(userId, deviceResp.Value.UserId, "UserId should be correct. Got: %v. Want: %v",
				deviceResp.Value.UserId, userId)
			sCtx.Assert().Equal(platform, deviceResp.Value.Platform, "Platform should be correct. Got: %v. Want: %v",
				deviceResp.Value.Platform, platform)
			sCtx.Assert().True(deviceResp.Value.EnteredAt.Before(time.Now().UTC()), "Time should be correct. Got: %v. Want <= %v",
				deviceResp.Value.EnteredAt, time.Now().UTC())
			sCtx.WithNewAttachment("Response", allure.Text, []byte(fmt.Sprintf("%+v", deviceResp)))
		})

		// Update device
		t.WithNewStep("Update device", func(sCtx provider.StepCtx) {
			// Expected mock response
			var mockResp http.Response
			mockResp.StatusCode = 200
			bodyReader := strings.NewReader("{\"success\":true}")
			mockResp.Body = io.NopCloser(bodyReader)

			clientMock.EXPECT().UpdateDevice(deviceId, newPlatform, newUserId).Return(&mockResp, nil)

			// Test mock
			resp, err := clientMock.UpdateDevice(deviceId, newPlatform, newUserId)
			sCtx.Require().NoError(err)
			sCtx.Assert().Equal(http.StatusOK, resp.StatusCode, "Status code should be correct. Want: %v. Got %v",
				http.StatusOK, resp.StatusCode)

			UpdateResp := new(UpdateDeviceResp)
			err = json.NewDecoder(resp.Body).Decode(&UpdateResp)
			sCtx.Require().NoError(err)
			sCtx.Assert().True(UpdateResp.Success, "Success should be true")
		})

		// Check that device was updated
		t.WithNewStep("Check that device was updated", func(sCtx provider.StepCtx) {
			// Expected mock response
			var mockResp http.Response
			mockResp.StatusCode = 200
			bodyReader := strings.NewReader("{\"value\":{\"id\":\"478\", \"platform\":\"android\", \"userId\":\"2\", " +
				"\"enteredAt\":\"2022-11-01T00:36:10.353993Z\"}}")
			mockResp.Body = io.NopCloser(bodyReader)

			clientMock.EXPECT().GetDevice(deviceId).Return(&mockResp, nil)

			// Test mock
			resp, err := clientMock.GetDevice(deviceId)
			sCtx.Require().NoError(err)

			deviceRespUpd := new(GetDeviceResp)
			err = json.NewDecoder(resp.Body).Decode(&deviceRespUpd)
			sCtx.Require().NoError(err)

			sCtx.Assert().Equal(newPlatform, deviceRespUpd.Value.Platform, "Platform should be correct. Got: %v. Want: %v",
				deviceRespUpd.Value.Platform, newPlatform)
			sCtx.Assert().Equal(newUserId, deviceRespUpd.Value.UserId, "UserId should be correct. Got: %v. Want: %v",
				deviceRespUpd.Value.UserId, newUserId)
			sCtx.WithNewAttachment("Response", allure.Text, []byte(fmt.Sprintf("%+v", deviceRespUpd)))
		})

		// Remove device
		t.WithNewStep("Remove device", func(sCtx provider.StepCtx) {
			// Expected mock response
			var mockResp http.Response
			mockResp.StatusCode = 200
			bodyReader := strings.NewReader("{\"found\":true}")
			mockResp.Body = io.NopCloser(bodyReader)

			clientMock.EXPECT().RemoveDevice(deviceId).Return(&mockResp, nil)

			// Test mock
			resp, err := clientMock.RemoveDevice(deviceId)
			sCtx.Require().NoError(err)
			sCtx.Assert().Equal(http.StatusOK, resp.StatusCode, "Status code should be correct. Want: %v. Got %v",
				http.StatusOK, resp.StatusCode)

			RemDeviceResp := new(RemoveDeviceResp)
			err = json.NewDecoder(resp.Body).Decode(&RemDeviceResp)
			sCtx.Require().NoError(err)
			sCtx.Assert().True(RemDeviceResp.Found, "Found should be true")
		})

		// Check that device was removed
		t.WithNewStep("Check that device was removed", func(sCtx provider.StepCtx) {
			// Expected mock response
			var mockResp http.Response
			mockResp.StatusCode = 404
			bodyReader := strings.NewReader("{\"code\":5, \"message\":\"device not found\", \"details\":[]}")
			mockResp.Body = io.NopCloser(bodyReader)

			clientMock.EXPECT().GetDevice(deviceId).Return(&mockResp, nil)

			// Test mock
			resp, err := clientMock.GetDevice(deviceId)
			sCtx.Require().NoError(err)
			sCtx.Assert().Equal(http.StatusNotFound, resp.StatusCode, "Status code should be correct. Want: %v. Got %v",
				http.StatusNotFound, resp.StatusCode)
		})
	})
}
