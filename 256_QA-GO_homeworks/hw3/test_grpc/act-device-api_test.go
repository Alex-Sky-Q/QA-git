//go:build e2e_test_noallure

package test_grpc

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	act_device_api "gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api/gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api"
	"testing"
	"time"
)

// Basic device lifecycle - create, update, get and remove device
func TestDeviceLifecycle(t *testing.T) {
	client := NewGrpcClient()
	//ctx := context.Background()
	ctx := SetTimeoutCtx(5)
	t.Run("Test device lifecycle (CRUD)", func(t *testing.T) {
		platform := "ios"
		var userId uint64 = 1
		newPlatform := "android"
		var newUserId uint64 = 2

		// Create a device
		reqCreate := act_device_api.CreateDeviceV1Request{
			Platform: platform,
			UserId:   userId,
		}
		respCreate, err := client.CreateDeviceV1(ctx, &reqCreate)
		deviceId := respCreate.DeviceId
		require.NoError(t, err)
		require.NotEmptyf(t, deviceId, "Device ID is empty")

		// Check that device was created
		reqDesc := act_device_api.DescribeDeviceV1Request{DeviceId: deviceId}
		respDesc, err := client.DescribeDeviceV1(ctx, &reqDesc)
		require.NoError(t, err)
		assert.Equalf(t, respDesc.Value.Platform, platform, "Platform is not correct. Got: %v. Want: %v",
			respDesc.Value.Platform, platform)
		assert.Equalf(t, respDesc.Value.UserId, userId, "UserId is not correct. Got: %v. Want: %v",
			respDesc.Value.UserId, userId)
		assert.GreaterOrEqualf(t, time.Now().Unix(), respDesc.Value.EnteredAt.Seconds, "Time is not correct. "+
			"%v is not greater or equal than %v", time.Now().Unix(), respDesc.Value.EnteredAt.Seconds)

		// Update device
		reqUpd := act_device_api.UpdateDeviceV1Request{
			DeviceId: deviceId,
			Platform: newPlatform,
			UserId:   newUserId,
		}
		respUpd, err := client.UpdateDeviceV1(ctx, &reqUpd)
		require.NoError(t, err)
		assert.Truef(t, respUpd.Success, "Success is not True")

		// Check that device was updated
		respDesc, err = client.DescribeDeviceV1(ctx, &reqDesc)
		require.NoError(t, err)
		assert.Equalf(t, respDesc.Value.Platform, newPlatform, "Platform is not correct. Got: %v. Want: %v",
			respDesc.Value.Platform, newPlatform)
		assert.Equalf(t, respDesc.Value.UserId, newUserId, "UserId is not correct. Got: %v. Want: %v",
			respDesc.Value.UserId, newUserId)
		assert.GreaterOrEqualf(t, time.Now().Unix(), respDesc.Value.EnteredAt.Seconds, "Time is not correct. "+
			"%v is not greater or equal than %v", time.Now().Unix(), respDesc.Value.EnteredAt.Seconds)

		// Remove device
		reqRem := act_device_api.RemoveDeviceV1Request{DeviceId: deviceId}
		respRem, err := client.RemoveDeviceV1(ctx, &reqRem)
		require.NoError(t, err)
		assert.Truef(t, respRem.Found, "Found is not True")

		// Check that device was removed
		respDesc, err = client.DescribeDeviceV1(ctx, &reqDesc)
		assert.Equalf(t, err.Error(), "rpc error: code = NotFound desc = device not found", "Error code "+
			"is not correct")
	})
}

// Test devices list
func TestDevicesList(t *testing.T) {
	client := NewGrpcClient()
	ctx := SetTimeoutCtx(5)
	t.Run("Test devices list", func(t *testing.T) {
		platform1 := "ios"
		var userId1 uint64 = 1
		platform2 := "android"
		var userId2 uint64 = 2
		platform3 := "android"
		var userId3 uint64 = 3
		var deviceIds []uint64

		// Create devices
		reqCreate := act_device_api.CreateDeviceV1Request{
			Platform: platform1,
			UserId:   userId1,
		}
		respCreate, err := client.CreateDeviceV1(ctx, &reqCreate)
		deviceId := respCreate.DeviceId
		require.NoError(t, err)
		require.NotEmptyf(t, deviceId, "Device ID is empty")
		deviceIds = append(deviceIds, deviceId)

		reqCreate = act_device_api.CreateDeviceV1Request{
			Platform: platform2,
			UserId:   userId2,
		}
		respCreate, err = client.CreateDeviceV1(ctx, &reqCreate)
		deviceId = respCreate.DeviceId
		require.NoError(t, err)
		require.NotEmptyf(t, deviceId, "Device ID is empty")
		deviceIds = append(deviceIds, deviceId)

		reqCreate = act_device_api.CreateDeviceV1Request{
			Platform: platform3,
			UserId:   userId3,
		}
		respCreate, err = client.CreateDeviceV1(ctx, &reqCreate)
		deviceToDel := respCreate.DeviceId
		require.NoError(t, err)
		require.NotEmptyf(t, deviceId, "Device ID is empty")

		reqRem := act_device_api.RemoveDeviceV1Request{DeviceId: deviceToDel}
		respRem, err := client.RemoveDeviceV1(ctx, &reqRem)
		require.NoError(t, err)
		assert.Truef(t, respRem.Found, "Found is not True")

		// List devices
		reqList := act_device_api.ListDevicesV1Request{
			Page:    1,
			PerPage: 15,
		}
		respList, err := client.ListDevicesV1(ctx, &reqList)
		require.NoError(t, err)

		// Check that deleted device not in list
		for _, device := range respList.Items {
			deviceId := device.Id
			if deviceId == deviceToDel {
				t.Fatalf("Deleted device is in list: %v", deviceToDel)
			}
		}

		// Check that devices are in list
		for _, x := range deviceIds {
			found := false
			for _, device := range respList.Items {
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
