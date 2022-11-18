//go:build e2e_test

package test_grpc

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	_ "github.com/lib/pq"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/256_homeworks/hw5"
	act_device_api "gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api/gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api"
	"testing"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "docker"
	password = "docker"
	dbname   = "act_device_api"
	sslmode  = "disable"
)

func TestMain(m *testing.M) {
	startTime := time.Now().UTC()

	defer func() {
		db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
			host, port, dbname, user, password, sslmode))
		if err != nil {
			panic(err)
		}
		// Clean tables
		tablesToClean := []string{"notification_events", "devices_events", "devices"}
		for _, t := range tablesToClean {
			query := fmt.Sprintf("DELETE FROM %s where created_at >= $1", t)
			_, err = db.Exec(query, startTime)
			if err != nil {
				panic(fmt.Errorf("DB cleaning error: %v", err))
			}
		}
		db.Close()
	}()

	m.Run()
}

func TestDeviceLifecycleGRPC(t *testing.T) {
	client := NewGrpcClient()
	//ctx := context.Background()
	ctx := SetTimeoutCtx(5)

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		host, port, dbname, user, password, sslmode))
	require.NoError(t, err)
	db.Mapper = reflectx.NewMapper("db")
	defer db.Close()

	runner.Run(t, "Test device lifecycle (CRUD)", func(t provider.T) {
		t.Title("gRPC - Test device lifecycle (CRUD)")
		t.Description("Basic device lifecycle - create, update, get and remove device")
		platform := "ios"
		var userId uint64 = 1
		newPlatform := "android"
		var newUserId uint64 = 2
		var deviceId uint64

		// Create a device
		t.WithNewStep("Create a device", func(sCtx provider.StepCtx) {
			reqCreate := act_device_api.CreateDeviceV1Request{
				Platform: platform,
				UserId:   userId,
			}
			respCreate, err := client.CreateDeviceV1(ctx, &reqCreate)
			sCtx.Require().NoError(err)
			deviceId = respCreate.DeviceId
			sCtx.Require().NotEmpty(deviceId, "Device ID should not be empty")
			sCtx.WithNewAttachment("Response", allure.Text, []byte(fmt.Sprint(respCreate)))
		})

		// Check that device was created
		t.WithNewStep("Check that device was created", func(sCtx provider.StepCtx) {
			reqDesc := act_device_api.DescribeDeviceV1Request{DeviceId: deviceId}
			respDesc, err := client.DescribeDeviceV1(ctx, &reqDesc)
			sCtx.Require().NoError(err)
			sCtx.Assert().Equal(platform, respDesc.Value.Platform, "Platform should be correct. Got: %v. Want: %v",
				respDesc.Value.Platform, platform)
			sCtx.Assert().Equal(userId, respDesc.Value.UserId, "UserId should be correct. Got: %v. Want: %v",
				respDesc.Value.UserId, userId)
			sCtx.Assert().GreaterOrEqual(time.Now().Unix(), respDesc.Value.EnteredAt.Seconds, "Time should be correct. "+
				"%v should be greater or equal than %v", time.Now().Unix(), respDesc.Value.EnteredAt.Seconds)
			sCtx.WithNewAttachment("Response", allure.Text, []byte(fmt.Sprint(respDesc)))
			// Check DB
			var res hw5.DevicesTable
			query, err := db.Preparex("select * from devices where id=$1")
			sCtx.Require().NoError(err)
			err = query.Get(&res, deviceId)
			sCtx.Require().NoError(err)
			sCtx.Assert().Equal(platform, res.Platform, "Platforms should match. Got: %d. Want: %d",
				res.Platform, platform)
			sCtx.Assert().Equal(userId, res.UserId, "User IDs should match. Got: %d. Want: %d",
				res.UserId, userId)
			sCtx.Assert().False(res.Removed, "Removed should be false")
			sCtx.Assert().GreaterOrEqual(time.Now().UTC(), res.EnteredAt, "EnteredAt time should be correct. "+
				"%v should be greater or equal than %v", time.Now().UTC(), res.EnteredAt)
			sCtx.Assert().GreaterOrEqual(time.Now().UTC(), res.CreatedAt, "CreatedAt time should be correct. "+
				"%v should be greater or equal than %v", time.Now().UTC(), res.CreatedAt)
			sCtx.Assert().GreaterOrEqual(time.Now().UTC(), res.UpdatedAt, "UpdatedAt time should be correct. "+
				"%v should be greater or equal than %v", time.Now().UTC(), res.UpdatedAt)
		})

		// Update device
		t.WithNewStep("Update device", func(sCtx provider.StepCtx) {
			reqUpd := act_device_api.UpdateDeviceV1Request{
				DeviceId: deviceId,
				Platform: newPlatform,
				UserId:   newUserId,
			}
			respUpd, err := client.UpdateDeviceV1(ctx, &reqUpd)
			sCtx.Require().NoError(err)
			sCtx.Assert().True(respUpd.Success, "Success should be True")
		})

		// Check that device was updated
		t.WithNewStep("Check that device was updated", func(sCtx provider.StepCtx) {
			reqDesc := act_device_api.DescribeDeviceV1Request{DeviceId: deviceId}
			respDesc, err := client.DescribeDeviceV1(ctx, &reqDesc)
			sCtx.Require().NoError(err)
			sCtx.Assert().Equal(newPlatform, respDesc.Value.Platform, "Platform should be correct. Got: %v. Want: %v",
				respDesc.Value.Platform, newPlatform)
			sCtx.Assert().Equal(newUserId, respDesc.Value.UserId, "UserId should be correct. Got: %v. Want: %v",
				respDesc.Value.UserId, newUserId)
			sCtx.Assert().GreaterOrEqual(time.Now().Unix(), respDesc.Value.EnteredAt.Seconds, "Time should be correct. "+
				"%v should be greater or equal than %v", time.Now().Unix(), respDesc.Value.EnteredAt.Seconds)
			sCtx.WithNewAttachment("Response", allure.Text, []byte(fmt.Sprint(respDesc)))
			// Check DB
			var res hw5.DevicesTable
			query, err := db.Preparex("select * from devices where id=$1")
			sCtx.Require().NoError(err)
			err = query.Get(&res, deviceId)
			sCtx.Require().NoError(err)
			sCtx.Assert().Equal(newPlatform, res.Platform, "Platforms should match. Got: %d. Want: %d",
				res.Platform, newPlatform)
			sCtx.Assert().Equal(newUserId, res.UserId, "User IDs should match. Got: %d. Want: %d",
				res.UserId, newUserId)
			sCtx.Assert().False(res.Removed, "Removed should be false")
			sCtx.Assert().GreaterOrEqual(time.Now().UTC(), res.EnteredAt, "EnteredAt time should be correct. "+
				"%v should be greater or equal than %v", time.Now().UTC(), res.EnteredAt)
			sCtx.Assert().GreaterOrEqual(time.Now().UTC(), res.CreatedAt, "CreatedAt time should be correct. "+
				"%v should be greater or equal than %v", time.Now().UTC(), res.CreatedAt)
			sCtx.Assert().GreaterOrEqual(time.Now().UTC(), res.UpdatedAt, "UpdatedAt time should be correct. "+
				"%v should be greater or equal than %v", time.Now().UTC(), res.UpdatedAt)
		})

		// Remove device
		t.WithNewStep("Remove device", func(sCtx provider.StepCtx) {
			reqRem := act_device_api.RemoveDeviceV1Request{DeviceId: deviceId}
			respRem, err := client.RemoveDeviceV1(ctx, &reqRem)
			sCtx.Require().NoError(err)
			sCtx.Assert().True(respRem.Found, "Found should be True")
		})

		// Check that device was removed
		t.WithNewStep("Check that device was removed", func(sCtx provider.StepCtx) {
			reqDesc := act_device_api.DescribeDeviceV1Request{DeviceId: deviceId}
			_, err := client.DescribeDeviceV1(ctx, &reqDesc)
			sCtx.Assert().Equal("rpc error: code = NotFound desc = device not found", err.Error(),
				"Error code should be correct")
			sCtx.WithNewAttachment("Response", allure.Text, []byte(fmt.Sprint(err)))
			// Check DB
			var res hw5.DevicesTable
			query, err := db.Preparex("select * from devices where id=$1")
			sCtx.Require().NoError(err)
			err = query.Get(&res, deviceId)
			sCtx.Require().NoError(err)
			sCtx.Assert().Equal(newPlatform, res.Platform, "Platforms should match. Got: %d. Want: %d",
				res.Platform, newPlatform)
			sCtx.Assert().Equal(newUserId, res.UserId, "User IDs should match. Got: %d. Want: %d",
				res.UserId, newUserId)
			sCtx.Assert().True(res.Removed, "Removed should be true")
			sCtx.Assert().GreaterOrEqual(time.Now().UTC(), res.EnteredAt, "EnteredAt time should be correct. "+
				"%v should be greater or equal than %v", time.Now().UTC(), res.EnteredAt)
			sCtx.Assert().GreaterOrEqual(time.Now().UTC(), res.CreatedAt, "CreatedAt time should be correct. "+
				"%v should be greater or equal than %v", time.Now().UTC(), res.CreatedAt)
			sCtx.Assert().GreaterOrEqual(time.Now().UTC(), res.UpdatedAt, "UpdatedAt time should be correct. "+
				"%v should be greater or equal than %v", time.Now().UTC(), res.UpdatedAt)
		})
	})
}

func TestDevicesListGRPC(t *testing.T) {
	client := NewGrpcClient()
	ctx := SetTimeoutCtx(5)
	runner.Run(t, "Test devices list", func(t provider.T) {
		t.Title("gRPC - Test devices list")
		t.Description("gRPC - Test ListDevicesV1 handle")
		platform1 := "ios"
		var userId1 uint64 = 1
		platform2 := "android"
		var userId2 uint64 = 2
		platform3 := "android"
		var userId3 uint64 = 3
		var deviceIds []uint64
		var deviceToDel uint64

		// Create devices
		t.WithNewStep("Create devices", func(sCtx provider.StepCtx) {
			reqCreate := act_device_api.CreateDeviceV1Request{
				Platform: platform1,
				UserId:   userId1,
			}
			respCreate, err := client.CreateDeviceV1(ctx, &reqCreate)
			deviceId := respCreate.DeviceId
			sCtx.Require().NoError(err)
			sCtx.Assert().NotEmpty(deviceId, "Device ID should not be empty")
			deviceIds = append(deviceIds, deviceId)

			reqCreate = act_device_api.CreateDeviceV1Request{
				Platform: platform2,
				UserId:   userId2,
			}
			respCreate, err = client.CreateDeviceV1(ctx, &reqCreate)
			deviceId = respCreate.DeviceId
			sCtx.Require().NoError(err)
			sCtx.Assert().NotEmpty(deviceId, "Device ID should not be empty")
			deviceIds = append(deviceIds, deviceId)

			reqCreate = act_device_api.CreateDeviceV1Request{
				Platform: platform3,
				UserId:   userId3,
			}
			respCreate, err = client.CreateDeviceV1(ctx, &reqCreate)
			deviceToDel = respCreate.DeviceId
			sCtx.Require().NoError(err)
			sCtx.Assert().NotEmpty(deviceId, "Device ID should not be empty")

			reqRem := act_device_api.RemoveDeviceV1Request{DeviceId: deviceToDel}
			respRem, err := client.RemoveDeviceV1(ctx, &reqRem)
			sCtx.Require().NoError(err)
			sCtx.Assert().True(respRem.Found, "Found should be True")
		})

		// List devices
		var respList *act_device_api.ListDevicesV1Response
		t.WithNewStep("List devices", func(sCtx provider.StepCtx) {
			reqList := act_device_api.ListDevicesV1Request{
				Page:    1,
				PerPage: 15,
			}
			var err error
			respList, err = client.ListDevicesV1(ctx, &reqList)
			sCtx.Require().NoError(err)
			sCtx.WithNewAttachment("Response", allure.Text, []byte(fmt.Sprint(respList)))
		})

		// Check that deleted device not in list
		t.WithNewStep("Check deleted device is not in list", func(sCtx provider.StepCtx) {
			for _, device := range respList.Items {
				deviceId := device.Id
				//sCtx.Assert().NotEqual(deviceToDel, deviceId, "Deleted device should not be in list: %v", deviceToDel)
				if deviceId == deviceToDel {
					t.Fatalf("Deleted device is in list: %v", deviceToDel)
				}
			}
		})

		// Check that devices are in list
		t.WithNewStep("Check that devices are in list", func(sCtx provider.StepCtx) {
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
	})
}
