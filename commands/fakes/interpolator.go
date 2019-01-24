// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/pivotal-cf/kiln/builder"
)

type Interpolator struct {
	InterpolateStub        func(builder.InterpolateInput, []byte) ([]byte, error)
	interpolateMutex       sync.RWMutex
	interpolateArgsForCall []struct {
		arg1 builder.InterpolateInput
		arg2 []byte
	}
	interpolateReturns struct {
		result1 []byte
		result2 error
	}
	interpolateReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Interpolator) Interpolate(arg1 builder.InterpolateInput, arg2 []byte) ([]byte, error) {
	var arg2Copy []byte
	if arg2 != nil {
		arg2Copy = make([]byte, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.interpolateMutex.Lock()
	ret, specificReturn := fake.interpolateReturnsOnCall[len(fake.interpolateArgsForCall)]
	fake.interpolateArgsForCall = append(fake.interpolateArgsForCall, struct {
		arg1 builder.InterpolateInput
		arg2 []byte
	}{arg1, arg2Copy})
	fake.recordInvocation("Interpolate", []interface{}{arg1, arg2Copy})
	fake.interpolateMutex.Unlock()
	if fake.InterpolateStub != nil {
		return fake.InterpolateStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.interpolateReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *Interpolator) InterpolateCallCount() int {
	fake.interpolateMutex.RLock()
	defer fake.interpolateMutex.RUnlock()
	return len(fake.interpolateArgsForCall)
}

func (fake *Interpolator) InterpolateCalls(stub func(builder.InterpolateInput, []byte) ([]byte, error)) {
	fake.interpolateMutex.Lock()
	defer fake.interpolateMutex.Unlock()
	fake.InterpolateStub = stub
}

func (fake *Interpolator) InterpolateArgsForCall(i int) (builder.InterpolateInput, []byte) {
	fake.interpolateMutex.RLock()
	defer fake.interpolateMutex.RUnlock()
	argsForCall := fake.interpolateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *Interpolator) InterpolateReturns(result1 []byte, result2 error) {
	fake.interpolateMutex.Lock()
	defer fake.interpolateMutex.Unlock()
	fake.InterpolateStub = nil
	fake.interpolateReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Interpolator) InterpolateReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.interpolateMutex.Lock()
	defer fake.interpolateMutex.Unlock()
	fake.InterpolateStub = nil
	if fake.interpolateReturnsOnCall == nil {
		fake.interpolateReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.interpolateReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Interpolator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.interpolateMutex.RLock()
	defer fake.interpolateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Interpolator) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
