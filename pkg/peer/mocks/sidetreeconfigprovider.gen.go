// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	"github.com/trustbloc/sidetree-fabric/pkg/config"
)

type SidetreeConfigProvider struct {
	ForChannelStub        func(channelID string) config.SidetreeService
	forChannelMutex       sync.RWMutex
	forChannelArgsForCall []struct {
		channelID string
	}
	forChannelReturns struct {
		result1 config.SidetreeService
	}
	forChannelReturnsOnCall map[int]struct {
		result1 config.SidetreeService
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *SidetreeConfigProvider) ForChannel(channelID string) config.SidetreeService {
	fake.forChannelMutex.Lock()
	ret, specificReturn := fake.forChannelReturnsOnCall[len(fake.forChannelArgsForCall)]
	fake.forChannelArgsForCall = append(fake.forChannelArgsForCall, struct {
		channelID string
	}{channelID})
	fake.recordInvocation("ForChannel", []interface{}{channelID})
	fake.forChannelMutex.Unlock()
	if fake.ForChannelStub != nil {
		return fake.ForChannelStub(channelID)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.forChannelReturns.result1
}

func (fake *SidetreeConfigProvider) ForChannelCallCount() int {
	fake.forChannelMutex.RLock()
	defer fake.forChannelMutex.RUnlock()
	return len(fake.forChannelArgsForCall)
}

func (fake *SidetreeConfigProvider) ForChannelArgsForCall(i int) string {
	fake.forChannelMutex.RLock()
	defer fake.forChannelMutex.RUnlock()
	return fake.forChannelArgsForCall[i].channelID
}

func (fake *SidetreeConfigProvider) ForChannelReturns(result1 config.SidetreeService) {
	fake.ForChannelStub = nil
	fake.forChannelReturns = struct {
		result1 config.SidetreeService
	}{result1}
}

func (fake *SidetreeConfigProvider) ForChannelReturnsOnCall(i int, result1 config.SidetreeService) {
	fake.ForChannelStub = nil
	if fake.forChannelReturnsOnCall == nil {
		fake.forChannelReturnsOnCall = make(map[int]struct {
			result1 config.SidetreeService
		})
	}
	fake.forChannelReturnsOnCall[i] = struct {
		result1 config.SidetreeService
	}{result1}
}

func (fake *SidetreeConfigProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.forChannelMutex.RLock()
	defer fake.forChannelMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *SidetreeConfigProvider) recordInvocation(key string, args []interface{}) {
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
