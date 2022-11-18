//go:build e2e_test_noallure

package test_http

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

// Basic device lifecycle - create, update, get and remove device
func TestDeviceLifecycle(t *testing.T) {
	client := NewHTTPClient(5, 1*time.Second)
	t.Run("Test device lifecycle (CRUD)", func(t *testing.T) {
		platform := "ios"
		userId := "1"
		newPlatform := "android"
		newUserId := "2"

		// Create a device
		resp, err := client.CreateDevice(platform, userId)
		if err != nil {
			t.Fatal(err)
		}

		deviceId := ""

		// Post returns 200, write to issues
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Status code is not %v. Got %v instead", http.StatusOK, resp.StatusCode)
		} else {
			var createDevResp CreateDeviceResp
			err = json.NewDecoder(resp.Body).Decode(&createDevResp)
			if err != nil {
				t.Fatal(err)
			}
			deviceId = createDevResp.DeviceId
			if deviceId == "" {
				t.Fatal("Device ID is empty")
			}
		}

		// Check that device was created
		resp, err = client.GetDevice(deviceId)
		if err != nil {
			t.Fatalf("%s", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("%v", resp.StatusCode)
		}
		deviceResp := new(GetDeviceResp)
		err = json.NewDecoder(resp.Body).Decode(deviceResp)
		if err != nil {
			t.Fatal(err)
		}
		if deviceResp.Value.Platform != platform {
			t.Fatalf("Platform mismatch. Got: %v. Want: %v", deviceResp.Value.Platform, platform)
		} else if deviceResp.Value.UserId != userId {
			t.Fatalf("UserId mismatch. Got: %v. Want: %v", deviceResp.Value.UserId, userId)
		} else if !deviceResp.Value.EnteredAt.Before(time.Now()) {
			t.Fatalf("Time mismatch. Got: %v. Want: %v", deviceResp.Value.EnteredAt, time.Now())
		}

		// Update device
		resp, err = client.UpdateDevice(deviceId, newPlatform, newUserId)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Status code is not %v. Got %v instead", http.StatusOK, resp.StatusCode)
		}
		UpdateResp := new(UpdateDeviceResp)
		err = json.NewDecoder(resp.Body).Decode(&UpdateResp)
		if err != nil {
			t.Fatal(err)
		}
		if !UpdateResp.Success {
			t.Fatal("Success is not true")
		}

		// Check that device was updated
		resp, err = client.GetDevice(deviceId)
		if err != nil {
			t.Fatal(err)
		}
		deviceRespUpd := new(GetDeviceResp)
		err = json.NewDecoder(resp.Body).Decode(&deviceRespUpd)
		if err != nil {
			t.Fatal(err)
		}
		if deviceRespUpd.Value.Platform != newPlatform {
			t.Fatalf("Platform mismatch. Got: %v. Want: %v", deviceRespUpd.Value.Platform, newPlatform)
		} else if deviceRespUpd.Value.UserId != newUserId {
			t.Fatalf("UserId mismatch. Got: %v. Want: %v", deviceRespUpd.Value.UserId, newUserId)
		}

		// Remove device
		resp, err = client.RemoveDevice(deviceId)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Status code is not %v. Got %v instead", http.StatusOK, resp.StatusCode)
		}
		RemDeviceResp := new(RemoveDeviceResp)
		err = json.NewDecoder(resp.Body).Decode(&RemDeviceResp)
		if !RemDeviceResp.Found {
			t.Fatal("Found is not true")
		}

		// Check that device was removed
		resp, err = client.GetDevice(deviceId)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("Status code is not %v. Got %v instead", http.StatusOK, resp.StatusCode)
		}
	})
}

// Test devices list
func TestDevicesList(t *testing.T) {
	client := NewHTTPClient(5, 1*time.Second)
	t.Run("Test devices list", func(t *testing.T) {
		platform1 := "ios"
		userId1 := "1"
		platform2 := "android"
		userId2 := "2"
		platform3 := "android"
		userId3 := "3"
		var deviceIds []string
		resp, err := client.CreateDevice(platform1, userId1)
		if err != nil {
			t.Fatal(err)
		}
		DeviceResp := new(CreateDeviceResp)
		err = json.NewDecoder(resp.Body).Decode(&DeviceResp)
		if err != nil {
			t.Fatal(err)
		}
		deviceIds = append(deviceIds, DeviceResp.DeviceId)
		resp, err = client.CreateDevice(platform2, userId2)
		if err != nil {
			t.Fatal(err)
		}
		DeviceResp = new(CreateDeviceResp)
		err = json.NewDecoder(resp.Body).Decode(&DeviceResp)
		if err != nil {
			t.Fatal(err)
		}
		deviceIds = append(deviceIds, DeviceResp.DeviceId)

		resp, err = client.CreateDevice(platform3, userId3)
		if err != nil {
			t.Fatal(err)
		}
		DeviceResp = new(CreateDeviceResp)
		err = json.NewDecoder(resp.Body).Decode(&DeviceResp)
		if err != nil {
			t.Fatal(err)
		}
		deviceToDel := DeviceResp.DeviceId

		resp, err = client.RemoveDevice(deviceToDel)
		if err != nil {
			t.Fatal(err)
		}

		// List devices
		resp, err = client.ListDevices("1", "15")
		if err != nil {
			t.Fatal(err)
		}
		DeviceListResp := new(ListDevicesResp)
		err = json.NewDecoder(resp.Body).Decode(&DeviceListResp)
		if err != nil {
			t.Fatal(err)
		}

		// Check that deleted device not in list
		for _, device := range DeviceListResp.Items {
			deviceId := device.Id
			if deviceId == deviceToDel {
				t.Fatalf("Deleted device is in list: %v", deviceToDel)
			}
		}

		// Check that devices are in list
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
}
