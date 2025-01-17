// Code generated by mockery v2.22.1. DO NOT EDIT.

package synchronizer

import (
	context "context"
	big "math/big"

	mock "github.com/stretchr/testify/mock"

	types "github.com/0xPolygon/cdk-validium-node/jsonrpc/types"
)

// zkEVMClientMock is an autogenerated mock type for the zkEVMClientInterface type
type zkEVMClientMock struct {
	mock.Mock
}

// BatchByNumber provides a mock function with given fields: ctx, number
func (_m *zkEVMClientMock) BatchByNumber(ctx context.Context, number *big.Int) (*types.Batch, error) {
	ret := _m.Called(ctx, number)

	var r0 *types.Batch
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int) (*types.Batch, error)); ok {
		return rf(ctx, number)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *big.Int) *types.Batch); ok {
		r0 = rf(ctx, number)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Batch)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *big.Int) error); ok {
		r1 = rf(ctx, number)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// BatchNumber provides a mock function with given fields: ctx
func (_m *zkEVMClientMock) BatchNumber(ctx context.Context) (uint64, error) {
	ret := _m.Called(ctx)

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (uint64, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) uint64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewZkEVMClientMock interface {
	mock.TestingT
	Cleanup(func())
}

// newZkEVMClientMock creates a new instance of zkEVMClientMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newZkEVMClientMock(t mockConstructorTestingTnewZkEVMClientMock) *zkEVMClientMock {
	mock := &zkEVMClientMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
