// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"
)

type Logger struct {
	PrintlnStub        func(...interface{})
	printlnMutex       sync.RWMutex
	printlnArgsForCall []struct {
		arg1 []interface{}
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Logger) Println(arg1 ...interface{}) {
	fake.printlnMutex.Lock()
	fake.printlnArgsForCall = append(fake.printlnArgsForCall, struct {
		arg1 []interface{}
	}{arg1})
	fake.recordInvocation("Println", []interface{}{arg1})
	fake.printlnMutex.Unlock()
	if fake.PrintlnStub != nil {
		fake.PrintlnStub(arg1...)
	}
}

func (fake *Logger) PrintlnCallCount() int {
	fake.printlnMutex.RLock()
	defer fake.printlnMutex.RUnlock()
	return len(fake.printlnArgsForCall)
}

func (fake *Logger) PrintlnCalls(stub func(...interface{})) {
	fake.printlnMutex.Lock()
	defer fake.printlnMutex.Unlock()
	fake.PrintlnStub = stub
}

func (fake *Logger) PrintlnArgsForCall(i int) []interface{} {
	fake.printlnMutex.RLock()
	defer fake.printlnMutex.RUnlock()
	argsForCall := fake.printlnArgsForCall[i]
	return argsForCall.arg1
}

func (fake *Logger) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.printlnMutex.RLock()
	defer fake.printlnMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Logger) recordInvocation(key string, args []interface{}) {
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
