// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pivotal-cf/kiln/internal/commands"
	"github.com/pivotal-cf/kiln/pkg/cargo"
)

type HistoricKilnfileLock struct {
	Stub        func(*git.Repository, plumbing.Hash, string) (cargo.KilnfileLock, error)
	mutex       sync.RWMutex
	argsForCall []struct {
		arg1 *git.Repository
		arg2 plumbing.Hash
		arg3 string
	}
	returns struct {
		result1 cargo.KilnfileLock
		result2 error
	}
	returnsOnCall map[int]struct {
		result1 cargo.KilnfileLock
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *HistoricKilnfileLock) Spy(arg1 *git.Repository, arg2 plumbing.Hash, arg3 string) (cargo.KilnfileLock, error) {
	fake.mutex.Lock()
	ret, specificReturn := fake.returnsOnCall[len(fake.argsForCall)]
	fake.argsForCall = append(fake.argsForCall, struct {
		arg1 *git.Repository
		arg2 plumbing.Hash
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.Stub
	returns := fake.returns
	fake.recordInvocation("HistoricKilnfileLockFunc", []interface{}{arg1, arg2, arg3})
	fake.mutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return returns.result1, returns.result2
}

func (fake *HistoricKilnfileLock) CallCount() int {
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	return len(fake.argsForCall)
}

func (fake *HistoricKilnfileLock) Calls(stub func(*git.Repository, plumbing.Hash, string) (cargo.KilnfileLock, error)) {
	fake.mutex.Lock()
	defer fake.mutex.Unlock()
	fake.Stub = stub
}

func (fake *HistoricKilnfileLock) ArgsForCall(i int) (*git.Repository, plumbing.Hash, string) {
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	return fake.argsForCall[i].arg1, fake.argsForCall[i].arg2, fake.argsForCall[i].arg3
}

func (fake *HistoricKilnfileLock) Returns(result1 cargo.KilnfileLock, result2 error) {
	fake.mutex.Lock()
	defer fake.mutex.Unlock()
	fake.Stub = nil
	fake.returns = struct {
		result1 cargo.KilnfileLock
		result2 error
	}{result1, result2}
}

func (fake *HistoricKilnfileLock) ReturnsOnCall(i int, result1 cargo.KilnfileLock, result2 error) {
	fake.mutex.Lock()
	defer fake.mutex.Unlock()
	fake.Stub = nil
	if fake.returnsOnCall == nil {
		fake.returnsOnCall = make(map[int]struct {
			result1 cargo.KilnfileLock
			result2 error
		})
	}
	fake.returnsOnCall[i] = struct {
		result1 cargo.KilnfileLock
		result2 error
	}{result1, result2}
}

func (fake *HistoricKilnfileLock) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.mutex.RLock()
	defer fake.mutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *HistoricKilnfileLock) recordInvocation(key string, args []interface{}) {
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

var _ commands.HistoricKilnfileLockFunc = new(HistoricKilnfileLock).Spy
