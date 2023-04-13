// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"io"
	"sync"

	"github.com/pivotal-cf/kiln/pkg/planitest/internal"
)

type RenderService struct {
	RenderManifestStub        func(io.Reader, io.Reader) (string, error)
	renderManifestMutex       sync.RWMutex
	renderManifestArgsForCall []struct {
		arg1 io.Reader
		arg2 io.Reader
	}
	renderManifestReturns struct {
		result1 string
		result2 error
	}
	renderManifestReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *RenderService) RenderManifest(arg1 io.Reader, arg2 io.Reader) (string, error) {
	fake.renderManifestMutex.Lock()
	ret, specificReturn := fake.renderManifestReturnsOnCall[len(fake.renderManifestArgsForCall)]
	fake.renderManifestArgsForCall = append(fake.renderManifestArgsForCall, struct {
		arg1 io.Reader
		arg2 io.Reader
	}{arg1, arg2})
	fake.recordInvocation("RenderManifest", []interface{}{arg1, arg2})
	fake.renderManifestMutex.Unlock()
	if fake.RenderManifestStub != nil {
		return fake.RenderManifestStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.renderManifestReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *RenderService) RenderManifestCallCount() int {
	fake.renderManifestMutex.RLock()
	defer fake.renderManifestMutex.RUnlock()
	return len(fake.renderManifestArgsForCall)
}

func (fake *RenderService) RenderManifestCalls(stub func(io.Reader, io.Reader) (string, error)) {
	fake.renderManifestMutex.Lock()
	defer fake.renderManifestMutex.Unlock()
	fake.RenderManifestStub = stub
}

func (fake *RenderService) RenderManifestArgsForCall(i int) (io.Reader, io.Reader) {
	fake.renderManifestMutex.RLock()
	defer fake.renderManifestMutex.RUnlock()
	argsForCall := fake.renderManifestArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *RenderService) RenderManifestReturns(result1 string, result2 error) {
	fake.renderManifestMutex.Lock()
	defer fake.renderManifestMutex.Unlock()
	fake.RenderManifestStub = nil
	fake.renderManifestReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *RenderService) RenderManifestReturnsOnCall(i int, result1 string, result2 error) {
	fake.renderManifestMutex.Lock()
	defer fake.renderManifestMutex.Unlock()
	fake.RenderManifestStub = nil
	if fake.renderManifestReturnsOnCall == nil {
		fake.renderManifestReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.renderManifestReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *RenderService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.renderManifestMutex.RLock()
	defer fake.renderManifestMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *RenderService) recordInvocation(key string, args []interface{}) {
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

var _ internal.RenderService = new(RenderService)