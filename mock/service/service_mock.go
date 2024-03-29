// Code generated by mockery v2.42.0. DO NOT EDIT.

package service_mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/Hymiside/lamoda-api/pkg/models"

	uuid "github.com/google/uuid"
)

// service is an autogenerated mock type for the service type
type ServiceMock struct {
	mock.Mock
}

// AvailabilityProductsByWarehouseID provides a mock function with given fields: ctx, warehouseID
func (_m *ServiceMock) AvailabilityProductsByWarehouseID(ctx context.Context, warehouseID int) ([]models.AvailabilityProducts, error) {
	ret := _m.Called(ctx, warehouseID)

	if len(ret) == 0 {
		panic("no return value specified for AvailabilityProductsByWarehouseID")
	}

	var r0 []models.AvailabilityProducts
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) ([]models.AvailabilityProducts, error)); ok {
		return rf(ctx, warehouseID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) []models.AvailabilityProducts); ok {
		r0 = rf(ctx, warehouseID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.AvailabilityProducts)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, warehouseID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ConfirmOrCancelReservedProducts provides a mock function with given fields: ctx, status, req
func (_m *ServiceMock) ConfirmOrCancelReservedProducts(ctx context.Context, status int, req models.CancelORConfirmProductsRequest) error {
	ret := _m.Called(ctx, status, req)

	if len(ret) == 0 {
		panic("no return value specified for ConfirmOrCancelReservedProducts")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, models.CancelORConfirmProductsRequest) error); ok {
		r0 = rf(ctx, status, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Products provides a mock function with given fields: ctx
func (_m *ServiceMock) Products(ctx context.Context) ([]models.Product, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Products")
	}

	var r0 []models.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]models.Product, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []models.Product); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReservationProducts provides a mock function with given fields: ctx, data
func (_m *ServiceMock) ReservationProducts(ctx context.Context, data models.ReservationProductsRequest) (uuid.UUID, error) {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for ReservationProducts")
	}

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.ReservationProductsRequest) (uuid.UUID, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.ReservationProductsRequest) uuid.UUID); ok {
		r0 = rf(ctx, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.ReservationProductsRequest) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}