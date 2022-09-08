// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"context"
	"log"
	"sync"

	"github.com/pivotal-cf/kiln/internal/component"
)

type ReleaseSource struct {
	ConfigurationErrorsStub        func() []error
	configurationErrorsMutex       sync.RWMutex
	configurationErrorsArgsForCall []struct {
	}
	configurationErrorsReturns struct {
		result1 []error
	}
	configurationErrorsReturnsOnCall map[int]struct {
		result1 []error
	}
	DownloadReleaseStub        func(context.Context, *log.Logger, string, component.Lock) (component.Local, error)
	downloadReleaseMutex       sync.RWMutex
	downloadReleaseArgsForCall []struct {
		arg1 context.Context
		arg2 *log.Logger
		arg3 string
		arg4 component.Lock
	}
	downloadReleaseReturns struct {
		result1 component.Local
		result2 error
	}
	downloadReleaseReturnsOnCall map[int]struct {
		result1 component.Local
		result2 error
	}
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
	GetMatchedReleaseStub        func(context.Context, *log.Logger, component.Spec) (component.Lock, error)
	getMatchedReleaseMutex       sync.RWMutex
	getMatchedReleaseArgsForCall []struct {
		arg1 context.Context
		arg2 *log.Logger
		arg3 component.Spec
	}
	getMatchedReleaseReturns struct {
		result1 component.Lock
		result2 error
	}
	getMatchedReleaseReturnsOnCall map[int]struct {
		result1 component.Lock
		result2 error
	}
	IDStub        func() string
	iDMutex       sync.RWMutex
	iDArgsForCall []struct {
	}
	iDReturns struct {
		result1 string
	}
	iDReturnsOnCall map[int]struct {
		result1 string
	}
	IsPublishableStub        func() bool
	isPublishableMutex       sync.RWMutex
	isPublishableArgsForCall []struct {
	}
	isPublishableReturns struct {
		result1 bool
	}
	isPublishableReturnsOnCall map[int]struct {
		result1 bool
	}
	TypeStub        func() string
	typeMutex       sync.RWMutex
	typeArgsForCall []struct {
	}
	typeReturns struct {
		result1 string
	}
	typeReturnsOnCall map[int]struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ReleaseSource) ConfigurationErrors() []error {
	fake.configurationErrorsMutex.Lock()
	ret, specificReturn := fake.configurationErrorsReturnsOnCall[len(fake.configurationErrorsArgsForCall)]
	fake.configurationErrorsArgsForCall = append(fake.configurationErrorsArgsForCall, struct {
	}{})
	stub := fake.ConfigurationErrorsStub
	fakeReturns := fake.configurationErrorsReturns
	fake.recordInvocation("ConfigurationErrors", []interface{}{})
	fake.configurationErrorsMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *ReleaseSource) ConfigurationErrorsCallCount() int {
	fake.configurationErrorsMutex.RLock()
	defer fake.configurationErrorsMutex.RUnlock()
	return len(fake.configurationErrorsArgsForCall)
}

func (fake *ReleaseSource) ConfigurationErrorsCalls(stub func() []error) {
	fake.configurationErrorsMutex.Lock()
	defer fake.configurationErrorsMutex.Unlock()
	fake.ConfigurationErrorsStub = stub
}

func (fake *ReleaseSource) ConfigurationErrorsReturns(result1 []error) {
	fake.configurationErrorsMutex.Lock()
	defer fake.configurationErrorsMutex.Unlock()
	fake.ConfigurationErrorsStub = nil
	fake.configurationErrorsReturns = struct {
		result1 []error
	}{result1}
}

func (fake *ReleaseSource) ConfigurationErrorsReturnsOnCall(i int, result1 []error) {
	fake.configurationErrorsMutex.Lock()
	defer fake.configurationErrorsMutex.Unlock()
	fake.ConfigurationErrorsStub = nil
	if fake.configurationErrorsReturnsOnCall == nil {
		fake.configurationErrorsReturnsOnCall = make(map[int]struct {
			result1 []error
		})
	}
	fake.configurationErrorsReturnsOnCall[i] = struct {
		result1 []error
	}{result1}
}

func (fake *ReleaseSource) DownloadRelease(arg1 context.Context, arg2 *log.Logger, arg3 string, arg4 component.Lock) (component.Local, error) {
	fake.downloadReleaseMutex.Lock()
	ret, specificReturn := fake.downloadReleaseReturnsOnCall[len(fake.downloadReleaseArgsForCall)]
	fake.downloadReleaseArgsForCall = append(fake.downloadReleaseArgsForCall, struct {
		arg1 context.Context
		arg2 *log.Logger
		arg3 string
		arg4 component.Lock
	}{arg1, arg2, arg3, arg4})
	stub := fake.DownloadReleaseStub
	fakeReturns := fake.downloadReleaseReturns
	fake.recordInvocation("DownloadRelease", []interface{}{arg1, arg2, arg3, arg4})
	fake.downloadReleaseMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ReleaseSource) DownloadReleaseCallCount() int {
	fake.downloadReleaseMutex.RLock()
	defer fake.downloadReleaseMutex.RUnlock()
	return len(fake.downloadReleaseArgsForCall)
}

func (fake *ReleaseSource) DownloadReleaseCalls(stub func(context.Context, *log.Logger, string, component.Lock) (component.Local, error)) {
	fake.downloadReleaseMutex.Lock()
	defer fake.downloadReleaseMutex.Unlock()
	fake.DownloadReleaseStub = stub
}

func (fake *ReleaseSource) DownloadReleaseArgsForCall(i int) (context.Context, *log.Logger, string, component.Lock) {
	fake.downloadReleaseMutex.RLock()
	defer fake.downloadReleaseMutex.RUnlock()
	argsForCall := fake.downloadReleaseArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *ReleaseSource) DownloadReleaseReturns(result1 component.Local, result2 error) {
	fake.downloadReleaseMutex.Lock()
	defer fake.downloadReleaseMutex.Unlock()
	fake.DownloadReleaseStub = nil
	fake.downloadReleaseReturns = struct {
		result1 component.Local
		result2 error
	}{result1, result2}
}

func (fake *ReleaseSource) DownloadReleaseReturnsOnCall(i int, result1 component.Local, result2 error) {
	fake.downloadReleaseMutex.Lock()
	defer fake.downloadReleaseMutex.Unlock()
	fake.DownloadReleaseStub = nil
	if fake.downloadReleaseReturnsOnCall == nil {
		fake.downloadReleaseReturnsOnCall = make(map[int]struct {
			result1 component.Local
			result2 error
		})
	}
	fake.downloadReleaseReturnsOnCall[i] = struct {
		result1 component.Local
		result2 error
	}{result1, result2}
}

func (fake *ReleaseSource) FindReleaseVersion(arg1 context.Context, arg2 *log.Logger, arg3 component.Spec) (component.Lock, error) {
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

func (fake *ReleaseSource) FindReleaseVersionCallCount() int {
	fake.findReleaseVersionMutex.RLock()
	defer fake.findReleaseVersionMutex.RUnlock()
	return len(fake.findReleaseVersionArgsForCall)
}

func (fake *ReleaseSource) FindReleaseVersionCalls(stub func(context.Context, *log.Logger, component.Spec) (component.Lock, error)) {
	fake.findReleaseVersionMutex.Lock()
	defer fake.findReleaseVersionMutex.Unlock()
	fake.FindReleaseVersionStub = stub
}

func (fake *ReleaseSource) FindReleaseVersionArgsForCall(i int) (context.Context, *log.Logger, component.Spec) {
	fake.findReleaseVersionMutex.RLock()
	defer fake.findReleaseVersionMutex.RUnlock()
	argsForCall := fake.findReleaseVersionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *ReleaseSource) FindReleaseVersionReturns(result1 component.Lock, result2 error) {
	fake.findReleaseVersionMutex.Lock()
	defer fake.findReleaseVersionMutex.Unlock()
	fake.FindReleaseVersionStub = nil
	fake.findReleaseVersionReturns = struct {
		result1 component.Lock
		result2 error
	}{result1, result2}
}

func (fake *ReleaseSource) FindReleaseVersionReturnsOnCall(i int, result1 component.Lock, result2 error) {
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

func (fake *ReleaseSource) GetMatchedRelease(arg1 context.Context, arg2 *log.Logger, arg3 component.Spec) (component.Lock, error) {
	fake.getMatchedReleaseMutex.Lock()
	ret, specificReturn := fake.getMatchedReleaseReturnsOnCall[len(fake.getMatchedReleaseArgsForCall)]
	fake.getMatchedReleaseArgsForCall = append(fake.getMatchedReleaseArgsForCall, struct {
		arg1 context.Context
		arg2 *log.Logger
		arg3 component.Spec
	}{arg1, arg2, arg3})
	stub := fake.GetMatchedReleaseStub
	fakeReturns := fake.getMatchedReleaseReturns
	fake.recordInvocation("GetMatchedRelease", []interface{}{arg1, arg2, arg3})
	fake.getMatchedReleaseMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *ReleaseSource) GetMatchedReleaseCallCount() int {
	fake.getMatchedReleaseMutex.RLock()
	defer fake.getMatchedReleaseMutex.RUnlock()
	return len(fake.getMatchedReleaseArgsForCall)
}

func (fake *ReleaseSource) GetMatchedReleaseCalls(stub func(context.Context, *log.Logger, component.Spec) (component.Lock, error)) {
	fake.getMatchedReleaseMutex.Lock()
	defer fake.getMatchedReleaseMutex.Unlock()
	fake.GetMatchedReleaseStub = stub
}

func (fake *ReleaseSource) GetMatchedReleaseArgsForCall(i int) (context.Context, *log.Logger, component.Spec) {
	fake.getMatchedReleaseMutex.RLock()
	defer fake.getMatchedReleaseMutex.RUnlock()
	argsForCall := fake.getMatchedReleaseArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *ReleaseSource) GetMatchedReleaseReturns(result1 component.Lock, result2 error) {
	fake.getMatchedReleaseMutex.Lock()
	defer fake.getMatchedReleaseMutex.Unlock()
	fake.GetMatchedReleaseStub = nil
	fake.getMatchedReleaseReturns = struct {
		result1 component.Lock
		result2 error
	}{result1, result2}
}

func (fake *ReleaseSource) GetMatchedReleaseReturnsOnCall(i int, result1 component.Lock, result2 error) {
	fake.getMatchedReleaseMutex.Lock()
	defer fake.getMatchedReleaseMutex.Unlock()
	fake.GetMatchedReleaseStub = nil
	if fake.getMatchedReleaseReturnsOnCall == nil {
		fake.getMatchedReleaseReturnsOnCall = make(map[int]struct {
			result1 component.Lock
			result2 error
		})
	}
	fake.getMatchedReleaseReturnsOnCall[i] = struct {
		result1 component.Lock
		result2 error
	}{result1, result2}
}

func (fake *ReleaseSource) ID() string {
	fake.iDMutex.Lock()
	ret, specificReturn := fake.iDReturnsOnCall[len(fake.iDArgsForCall)]
	fake.iDArgsForCall = append(fake.iDArgsForCall, struct {
	}{})
	stub := fake.IDStub
	fakeReturns := fake.iDReturns
	fake.recordInvocation("ID", []interface{}{})
	fake.iDMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *ReleaseSource) IDCallCount() int {
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	return len(fake.iDArgsForCall)
}

func (fake *ReleaseSource) IDCalls(stub func() string) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = stub
}

func (fake *ReleaseSource) IDReturns(result1 string) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = nil
	fake.iDReturns = struct {
		result1 string
	}{result1}
}

func (fake *ReleaseSource) IDReturnsOnCall(i int, result1 string) {
	fake.iDMutex.Lock()
	defer fake.iDMutex.Unlock()
	fake.IDStub = nil
	if fake.iDReturnsOnCall == nil {
		fake.iDReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.iDReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *ReleaseSource) IsPublishable() bool {
	fake.isPublishableMutex.Lock()
	ret, specificReturn := fake.isPublishableReturnsOnCall[len(fake.isPublishableArgsForCall)]
	fake.isPublishableArgsForCall = append(fake.isPublishableArgsForCall, struct {
	}{})
	stub := fake.IsPublishableStub
	fakeReturns := fake.isPublishableReturns
	fake.recordInvocation("IsPublishable", []interface{}{})
	fake.isPublishableMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *ReleaseSource) IsPublishableCallCount() int {
	fake.isPublishableMutex.RLock()
	defer fake.isPublishableMutex.RUnlock()
	return len(fake.isPublishableArgsForCall)
}

func (fake *ReleaseSource) IsPublishableCalls(stub func() bool) {
	fake.isPublishableMutex.Lock()
	defer fake.isPublishableMutex.Unlock()
	fake.IsPublishableStub = stub
}

func (fake *ReleaseSource) IsPublishableReturns(result1 bool) {
	fake.isPublishableMutex.Lock()
	defer fake.isPublishableMutex.Unlock()
	fake.IsPublishableStub = nil
	fake.isPublishableReturns = struct {
		result1 bool
	}{result1}
}

func (fake *ReleaseSource) IsPublishableReturnsOnCall(i int, result1 bool) {
	fake.isPublishableMutex.Lock()
	defer fake.isPublishableMutex.Unlock()
	fake.IsPublishableStub = nil
	if fake.isPublishableReturnsOnCall == nil {
		fake.isPublishableReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isPublishableReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *ReleaseSource) Type() string {
	fake.typeMutex.Lock()
	ret, specificReturn := fake.typeReturnsOnCall[len(fake.typeArgsForCall)]
	fake.typeArgsForCall = append(fake.typeArgsForCall, struct {
	}{})
	stub := fake.TypeStub
	fakeReturns := fake.typeReturns
	fake.recordInvocation("Type", []interface{}{})
	fake.typeMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *ReleaseSource) TypeCallCount() int {
	fake.typeMutex.RLock()
	defer fake.typeMutex.RUnlock()
	return len(fake.typeArgsForCall)
}

func (fake *ReleaseSource) TypeCalls(stub func() string) {
	fake.typeMutex.Lock()
	defer fake.typeMutex.Unlock()
	fake.TypeStub = stub
}

func (fake *ReleaseSource) TypeReturns(result1 string) {
	fake.typeMutex.Lock()
	defer fake.typeMutex.Unlock()
	fake.TypeStub = nil
	fake.typeReturns = struct {
		result1 string
	}{result1}
}

func (fake *ReleaseSource) TypeReturnsOnCall(i int, result1 string) {
	fake.typeMutex.Lock()
	defer fake.typeMutex.Unlock()
	fake.TypeStub = nil
	if fake.typeReturnsOnCall == nil {
		fake.typeReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.typeReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *ReleaseSource) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.configurationErrorsMutex.RLock()
	defer fake.configurationErrorsMutex.RUnlock()
	fake.downloadReleaseMutex.RLock()
	defer fake.downloadReleaseMutex.RUnlock()
	fake.findReleaseVersionMutex.RLock()
	defer fake.findReleaseVersionMutex.RUnlock()
	fake.getMatchedReleaseMutex.RLock()
	defer fake.getMatchedReleaseMutex.RUnlock()
	fake.iDMutex.RLock()
	defer fake.iDMutex.RUnlock()
	fake.isPublishableMutex.RLock()
	defer fake.isPublishableMutex.RUnlock()
	fake.typeMutex.RLock()
	defer fake.typeMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ReleaseSource) recordInvocation(key string, args []interface{}) {
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

var _ component.ReleaseSource = new(ReleaseSource)
