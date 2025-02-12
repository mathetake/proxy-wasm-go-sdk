// Copyright 2020-2021 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type (
	rootContextState struct {
		context       types.RootContext
		httpCallbacks map[uint32]*httpCallbackAttribute
	}

	httpCallbackAttribute struct {
		callback        func(numHeaders, bodySize, numTrailers int)
		callerContextID uint32
	}
)

type state struct {
	newRootContext func(contextID uint32) types.RootContext
	rootContexts   map[uint32]*rootContextState
	streams        map[uint32]types.StreamContext
	httpStreams    map[uint32]types.HttpContext

	contextIDToRootID map[uint32]uint32
	activeContextID   uint32
}

var currentState = &state{
	rootContexts:      make(map[uint32]*rootContextState),
	httpStreams:       make(map[uint32]types.HttpContext),
	streams:           make(map[uint32]types.StreamContext),
	contextIDToRootID: make(map[uint32]uint32),
}

func SetNewRootContextFn(f func(contextID uint32) types.RootContext) {
	currentState.newRootContext = f
}

func RegisterHttpCallout(calloutID uint32, callback func(numHeaders, bodySize, numTrailers int)) {
	currentState.registerHttpCallOut(calloutID, callback)
}

func (s *state) createRootContext(contextID uint32) {
	var ctx types.RootContext
	if s.newRootContext == nil {
		ctx = &types.DefaultRootContext{}
	} else {
		ctx = s.newRootContext(contextID)
	}

	s.rootContexts[contextID] = &rootContextState{
		context:       ctx,
		httpCallbacks: map[uint32]*httpCallbackAttribute{},
	}

	// NOTE: this is a temporary work around for avoiding nil pointer panic
	// when users make http dispatch(es) on RootContext.
	// See https://github.com/tetratelabs/proxy-wasm-go-sdk/issues/110
	// TODO: refactor
	s.contextIDToRootID[contextID] = contextID
}

func (s *state) createStreamContext(contextID uint32, rootContextID uint32) bool {
	root, ok := s.rootContexts[rootContextID]
	if !ok {
		panic("invalid root context id")
	}

	if _, ok := s.streams[contextID]; ok {
		panic("context id duplicated")
	}

	ctx := root.context.NewStreamContext(contextID)
	if ctx == nil {
		// NewStreamContext is not defined by the user
		return false
	}
	s.contextIDToRootID[contextID] = rootContextID
	s.streams[contextID] = ctx
	return true
}

func (s *state) createHttpContext(contextID uint32, rootContextID uint32) bool {
	root, ok := s.rootContexts[rootContextID]
	if !ok {
		panic("invalid root context id")
	}

	if _, ok := s.httpStreams[contextID]; ok {
		panic("context id duplicated")
	}

	ctx := root.context.NewHttpContext(contextID)
	if ctx == nil {
		// NewHttpContext is not defined by the user
		return false
	}
	s.contextIDToRootID[contextID] = rootContextID
	s.httpStreams[contextID] = ctx
	return true
}

func (s *state) registerHttpCallOut(calloutID uint32, callback func(numHeaders, bodySize, numTrailers int)) {
	r := s.rootContexts[s.contextIDToRootID[s.activeContextID]]
	r.httpCallbacks[calloutID] = &httpCallbackAttribute{callback: callback, callerContextID: s.activeContextID}
}

//go:inline
func (s *state) setActiveContextID(contextID uint32) {
	s.activeContextID = contextID
}
