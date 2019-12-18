//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package runner

import "context"

type runnerMock struct {
	PreStartFunc func(ctx context.Context) error
	StartFunc    func(ctx context.Context) <-chan error
	PreStopFunc  func(ctx context.Context) error
	StopFunc     func(ctx context.Context) error
	PostStopFunc func(ctx context.Context) error
}

func (r *runnerMock) PreStart(ctx context.Context) error {
	if r.PreStartFunc != nil {
		return r.PreStartFunc(ctx)
	}
	return nil
}

func (r *runnerMock) Start(ctx context.Context) <-chan error {
	if r.StartFunc != nil {
		return r.StartFunc(ctx)
	}
	return nil
}

func (r *runnerMock) PreStop(ctx context.Context) error {
	if r.PreStopFunc != nil {
		return r.PreStopFunc(ctx)
	}
	return nil
}

func (r *runnerMock) Stop(ctx context.Context) error {
	if r.StopFunc != nil {
		return r.StopFunc(ctx)
	}
	return nil
}

func (r *runnerMock) PostStop(ctx context.Context) error {
	if r.PostStopFunc != nil {
		return r.PostStopFunc(ctx)
	}
	return nil
}
