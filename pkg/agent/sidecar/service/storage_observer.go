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

// Package service
package service

import (
	"context"
	"reflect"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file/watch"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
)

type StorageObserver interface {
	Start(ctx context.Context) (<-chan error, error)
}

type observer struct {
	w                    watch.Watcher
	dirs                 []string
	eg                   errgroup.Group
	checkDuration        time.Duration
	longestCheckDuration time.Duration
}

func New(opts ...Option) (so StorageObserver, err error) {
	o := new(observer)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(o); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	o.w, err = watch.New(
		watch.WithDirs(o.dirs...),
		watch.WithErrGroup(o.eg),
		watch.WithOnWrite(func(ctx context.Context, name string) error {
			ctx, span := trace.StartSpan(ctx, "vald/agent-sidecar/service/StorageObserver.watcher.OnWrite")
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			return o.backup(ctx)
		}),
		watch.WithOnCreate(func(ctx context.Context, name string) error {
			ctx, span := trace.StartSpan(ctx, "vald/agent-sidecar/service/StorageObserver.watcher.OnCreate")
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			return o.backup(ctx)
		}),
	)
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (o *observer) Start(ctx context.Context) (<-chan error, error) {
	wech, err := o.w.Start(ctx)
	if err != nil {
		return nil, err
	}
	ech := make(chan error, 100)
	o.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		ct := time.NewTicker(o.checkDuration)
		defer ct.Stop()
		lct := time.NewTicker(o.longestCheckDuration)
		defer lct.Stop()
		finalize := func() (err error) {
			err = ctx.Err()
			if err != nil && err != context.Canceled {
				return err
			}
			return nil
		}
		for {
			select {
			case <-ctx.Done():
				return finalize()
			case <-ct.C:
				err = o.backup(ctx)
				if err != nil {
					ech <- err
					log.Error(err)
					err = nil
				}
			case <-lct.C:
				err = o.backup(ctx)
				if err != nil {
					ech <- err
					log.Error(err)
					err = nil
				}
			case err = <-wech:
			}
			if err != nil {
				log.Error(err)
				select {
				case <-ctx.Done():
					return finalize()
				case ech <- err:
				}
			}
		}
	}))
	return ech, nil
}

func (o *observer) backup(ctx context.Context) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/agent-sidecar/service/StorageObserver.backup")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	// TODO implement backup logic here

	return nil
}
