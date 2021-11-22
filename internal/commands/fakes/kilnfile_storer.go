// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/pivotal-cf/kiln/internal/commands"
	"github.com/pivotal-cf/kiln/internal/commands/options"
	"github.com/pivotal-cf/kiln/pkg/cargo"
)

type KilnfileStorer struct {
	LoadStub        func(options.StandardOptionsEmbedder) (cargo.Kilnfile, cargo.KilnfileLock, error)
	loadMutex       sync.RWMutex
	loadArgsForCall []struct {
		arg1 options.StandardOptionsEmbedder
	}
	loadReturns struct {
		result1 cargo.Kilnfile
		result2 cargo.KilnfileLock
		result3 error
	}
	loadReturnsOnCall map[int]struct {
		result1 cargo.Kilnfile
		result2 cargo.KilnfileLock
		result3 error
	}
	SaveLockStub        func(string, cargo.KilnfileLock) error
	saveLockMutex       sync.RWMutex
	saveLockArgsForCall []struct {
		arg1 string
		arg2 cargo.KilnfileLock
	}
	saveLockReturns struct {
		result1 error
	}
	saveLockReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *KilnfileStorer) Load(arg1 options.StandardOptionsEmbedder) (cargo.Kilnfile, cargo.KilnfileLock, error) {
	fake.loadMutex.Lock()
	ret, specificReturn := fake.loadReturnsOnCall[len(fake.loadArgsForCall)]
	fake.loadArgsForCall = append(fake.loadArgsForCall, struct {
		arg1 options.StandardOptionsEmbedder
	}{arg1})
	stub := fake.LoadStub
	fakeReturns := fake.loadReturns
	fake.recordInvocation("Load", []interface{}{arg1})
	fake.loadMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *KilnfileStorer) LoadCallCount() int {
	fake.loadMutex.RLock()
	defer fake.loadMutex.RUnlock()
	return len(fake.loadArgsForCall)
}

func (fake *KilnfileStorer) LoadCalls(stub func(options.StandardOptionsEmbedder) (cargo.Kilnfile, cargo.KilnfileLock, error)) {
	fake.loadMutex.Lock()
	defer fake.loadMutex.Unlock()
	fake.LoadStub = stub
}

func (fake *KilnfileStorer) LoadArgsForCall(i int) options.StandardOptionsEmbedder {
	fake.loadMutex.RLock()
	defer fake.loadMutex.RUnlock()
	argsForCall := fake.loadArgsForCall[i]
	return argsForCall.arg1
}

func (fake *KilnfileStorer) LoadReturns(result1 cargo.Kilnfile, result2 cargo.KilnfileLock, result3 error) {
	fake.loadMutex.Lock()
	defer fake.loadMutex.Unlock()
	fake.LoadStub = nil
	fake.loadReturns = struct {
		result1 cargo.Kilnfile
		result2 cargo.KilnfileLock
		result3 error
	}{result1, result2, result3}
}

func (fake *KilnfileStorer) LoadReturnsOnCall(i int, result1 cargo.Kilnfile, result2 cargo.KilnfileLock, result3 error) {
	fake.loadMutex.Lock()
	defer fake.loadMutex.Unlock()
	fake.LoadStub = nil
	if fake.loadReturnsOnCall == nil {
		fake.loadReturnsOnCall = make(map[int]struct {
			result1 cargo.Kilnfile
			result2 cargo.KilnfileLock
			result3 error
		})
	}
	fake.loadReturnsOnCall[i] = struct {
		result1 cargo.Kilnfile
		result2 cargo.KilnfileLock
		result3 error
	}{result1, result2, result3}
}

func (fake *KilnfileStorer) SaveLock(arg1 string, arg2 cargo.KilnfileLock) error {
	fake.saveLockMutex.Lock()
	ret, specificReturn := fake.saveLockReturnsOnCall[len(fake.saveLockArgsForCall)]
	fake.saveLockArgsForCall = append(fake.saveLockArgsForCall, struct {
		arg1 string
		arg2 cargo.KilnfileLock
	}{arg1, arg2})
	stub := fake.SaveLockStub
	fakeReturns := fake.saveLockReturns
	fake.recordInvocation("SaveLock", []interface{}{arg1, arg2})
	fake.saveLockMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *KilnfileStorer) SaveLockCallCount() int {
	fake.saveLockMutex.RLock()
	defer fake.saveLockMutex.RUnlock()
	return len(fake.saveLockArgsForCall)
}

func (fake *KilnfileStorer) SaveLockCalls(stub func(string, cargo.KilnfileLock) error) {
	fake.saveLockMutex.Lock()
	defer fake.saveLockMutex.Unlock()
	fake.SaveLockStub = stub
}

func (fake *KilnfileStorer) SaveLockArgsForCall(i int) (string, cargo.KilnfileLock) {
	fake.saveLockMutex.RLock()
	defer fake.saveLockMutex.RUnlock()
	argsForCall := fake.saveLockArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *KilnfileStorer) SaveLockReturns(result1 error) {
	fake.saveLockMutex.Lock()
	defer fake.saveLockMutex.Unlock()
	fake.SaveLockStub = nil
	fake.saveLockReturns = struct {
		result1 error
	}{result1}
}

func (fake *KilnfileStorer) SaveLockReturnsOnCall(i int, result1 error) {
	fake.saveLockMutex.Lock()
	defer fake.saveLockMutex.Unlock()
	fake.SaveLockStub = nil
	if fake.saveLockReturnsOnCall == nil {
		fake.saveLockReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.saveLockReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *KilnfileStorer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.loadMutex.RLock()
	defer fake.loadMutex.RUnlock()
	fake.saveLockMutex.RLock()
	defer fake.saveLockMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *KilnfileStorer) recordInvocation(key string, args []interface{}) {
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

var _ commands.KilnfileStorer = new(KilnfileStorer)
