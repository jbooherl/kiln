// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"context"
	"log"
	"sync"

	"github.com/pivotal-cf/kiln/internal/component"
)

type ReleaseVersionFinder struct {
	FindReleaseVersionStub        func(context.Context, *log.Logger, component.Spec) (component.Lock, error)
	findReleaseVersionMutex       sync.RWMutex
	findReleaseVersionArgsForCall []struct {
		arg1 context.Context
		arg2 *log.Logger
		arg3 component.Spec
	}
	findReleaseVersionReturns struct {
		result1 component.Lock
		result2 error
	}
	findReleaseVersionReturnsOnCall map[int]struct {
		result1 component.Lock
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ReleaseVersionFinder) FindReleaseVersion(arg1 context.Context, arg2 *log.Logger, arg3 component.Spec) (component.Lock, error) {
	fake.findReleaseVersionMutex.Lock()
	ret, specificReturn := fake.findReleaseVersionReturnsOnCall[len(fake.findReleaseVersionArgsForCall)]
	fake.findReleaseVersionArgsForCall = append(fake.findReleaseVersionArgsForCall, struct {
		arg1 context.Context
		arg2 *log.Logger
		arg3 component.Spec
	}{arg1, arg2, arg3})
	stub := fake.FindReleaseVersionStub
	fakeReturns := fake.findReleaseVersionReturns
	fake.recordInvocation("FindReleaseVersion", []interface{}{arg1, arg2, arg3})
	fake.findReleaseVersionMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ReleaseVersionFinder) FindReleaseVersionCallCount() int {
	fake.findReleaseVersionMutex.RLock()
	defer fake.findReleaseVersionMutex.RUnlock()
	return len(fake.findReleaseVersionArgsForCall)
}

func (fake *ReleaseVersionFinder) FindReleaseVersionCalls(stub func(context.Context, *log.Logger, component.Spec) (component.Lock, error)) {
	fake.findReleaseVersionMutex.Lock()
	defer fake.findReleaseVersionMutex.Unlock()
	fake.FindReleaseVersionStub = stub
}

func (fake *ReleaseVersionFinder) FindReleaseVersionArgsForCall(i int) (context.Context, *log.Logger, component.Spec) {
	fake.findReleaseVersionMutex.RLock()
	defer fake.findReleaseVersionMutex.RUnlock()
	argsForCall := fake.findReleaseVersionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *ReleaseVersionFinder) FindReleaseVersionReturns(result1 component.Lock, result2 error) {
	fake.findReleaseVersionMutex.Lock()
	defer fake.findReleaseVersionMutex.Unlock()
	fake.FindReleaseVersionStub = nil
	fake.findReleaseVersionReturns = struct {
		result1 component.Lock
		result2 error
	}{result1, result2}
}

func (fake *ReleaseVersionFinder) FindReleaseVersionReturnsOnCall(i int, result1 component.Lock, result2 error) {
	fake.findReleaseVersionMutex.Lock()
	defer fake.findReleaseVersionMutex.Unlock()
	fake.FindReleaseVersionStub = nil
	if fake.findReleaseVersionReturnsOnCall == nil {
		fake.findReleaseVersionReturnsOnCall = make(map[int]struct {
			result1 component.Lock
			result2 error
		})
	}
	fake.findReleaseVersionReturnsOnCall[i] = struct {
		result1 component.Lock
		result2 error
	}{result1, result2}
}

func (fake *ReleaseVersionFinder) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.findReleaseVersionMutex.RLock()
	defer fake.findReleaseVersionMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ReleaseVersionFinder) recordInvocation(key string, args []interface{}) {
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

var _ component.ReleaseVersionFinder = new(ReleaseVersionFinder)
