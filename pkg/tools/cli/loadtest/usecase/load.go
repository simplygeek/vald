//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

package usecase

import (
	"context"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/pkg/agent/core/ngt/config"
)

type run struct {
	eg  errgroup.Group
	cfg *config.Data
}

func New(cfg *config.Data) (r runner.Runner, err error) {
	return nil, nil
}

func (r *run) PreStart(ctx context.Context) error {
	return nil
}

func (r *run) Start(ctx context.Context) (<-chan error, error) {
	return nil, nil
}

func (r *run) PreStop(ctx context.Context) error {
	return nil
}

func (r *run) Stop(ctx context.Context) error {
	return nil
}

func (r *run) PostStop(ctx context.Context) error {
	return nil
}
