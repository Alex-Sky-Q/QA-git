//go:build e2e_test

package test_http

import (
	"encoding/json"
	"fmt"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"net/http"
	"testing"
	"time"
)

func TestDeviceLifecycleHTTP(t *testing.T) {
	client := NewHTTPClient(5, 1*time.Second)
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
			resp, err := client.CreateDevice(platform, userId)
			sCtx.Require().NoError(err)
			// Post returns 200, not 201
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
			resp, err := client.GetDevice(deviceId)
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
			resp, err := client.UpdateDevice(deviceId, newPlatform, newUserId)
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
			resp, err := client.GetDevice(deviceId)
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
			resp, err := client.RemoveDevice(deviceId)
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
			resp, err := client.GetDevice(deviceId)
			sCtx.Require().NoError(err)
			sCtx.Assert().Equal(http.StatusNotFound, resp.StatusCode, "Status code should be correct. Want: %v. Got %v",
				http.StatusNotFound, resp.StatusCode)
		})
	})
}

func TestDevicesListHTTP(t *testing.T) {
	client := NewHTTPClient(5, 1*time.Second)
	runner.Run(t, "Test devices list", func(t provider.T) {
		t.Title("HTTP - Test devices list")
		t.Description("List devices and check all devices are shown")
		platform1 := "ios"
		userId1 := "1"
		platform2 := "android"
		userId2 := "2"
		platform3 := "android"
		userId3 := "3"
		var deviceIds []string
		deviceToDel := ""

		t.WithNewStep("Create devices", func(sCtx provider.StepCtx) {
			resp, err := client.CreateDevice(platform1, userId1)
			sCtx.Require().NoError(err)
			DeviceResp := new(CreateDeviceResp)
			err = json.NewDecoder(resp.Body).Decode(&DeviceResp)
			sCtx.Require().NoError(err)
			deviceIds = append(deviceIds, DeviceResp.DeviceId)

			resp, err = client.CreateDevice(platform2, userId2)
			sCtx.Require().NoError(err)
			DeviceResp = new(CreateDeviceResp)
			err = json.NewDecoder(resp.Body).Decode(&DeviceResp)
			sCtx.Require().NoError(err)
			deviceIds = append(deviceIds, DeviceResp.DeviceId)

			resp, err = client.CreateDevice(platform3, userId3)
			sCtx.Require().NoError(err)
			DeviceResp = new(CreateDeviceResp)
			err = json.NewDecoder(resp.Body).Decode(&DeviceResp)
			sCtx.Require().NoError(err)
			deviceToDel = DeviceResp.DeviceId

			resp, err = client.RemoveDevice(deviceToDel)
			sCtx.Require().NoError(err)
		})

		// List devices
		var DeviceListResp *ListDevicesResp
		t.WithNewStep("List devices", func(sCtx provider.StepCtx) {
			resp, err := client.ListDevices("1", "15")
			sCtx.Require().NoError(err)
			DeviceListResp = new(ListDevicesResp)
			err = json.NewDecoder(resp.Body).Decode(&DeviceListResp)
			sCtx.Require().NoError(err)
		})

		// Check that deleted device not in list
		t.WithNewStep("Check that deleted device not in list", func(sCtx provider.StepCtx) {
			for _, device := range DeviceListResp.Items {
				deviceId := device.Id
				if deviceId == deviceToDel {
					t.Fatalf("Deleted device is in list: %v", deviceToDel)
				}
			}
		})

		// Check that devices are in list
		t.WithNewStep("Check that devices are in list", func(sCtx provider.StepCtx) {
			for _, x := range deviceIds {
				found := false
				for _, device := range DeviceListResp.Items {
					deviceId := device.Id
					if x == deviceId {
						found = true
						break
					}
				}
				if !found {
					t.Fatal("Id not found")
				}
			}
		})
	})
}
