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

// Package tensorflow provides implementation of Go API for extract data to vector
package tensorflow

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/vdaas/vald/internal/errors"
)

type SessionOptions = tf.SessionOptions
type Operation = tf.Operation

type TF interface {
	GetVector(inputs ...string) ([]float64, error)
	GetValue(inputs ...string) (interface{}, error)
	GetValues(inputs ...string) (values []interface{}, err error)
	Close() error
}

type tensorflow struct {
	exportDir     string
	tags          []string
	feeds         []OutputSpec
	fetches       []OutputSpec
	operations    []*Operation
	sessionTarget string
	sessionConfig []byte
	options       *SessionOptions
	graph         *tf.Graph
	session       *tf.Session
	ndim          uint8
}

type OutputSpec struct {
	operationName string
	outputIndex   int
}

const (
	TwoDim uint8 = iota + 2
	ThreeDim
)

func New(opts ...Option) (TF, error) {
	t := new(tensorflow)
	for _, opt := range append(defaultOpts, opts...) {
		opt(t)
	}

	if t.options == nil && (len(t.sessionTarget) != 0 || t.sessionConfig != nil) {
		t.options = &tf.SessionOptions{
			Target: t.sessionTarget,
			Config: t.sessionConfig,
		}
	}

	model, err := tf.LoadSavedModel(t.exportDir, t.tags, t.options)
	if err != nil {
		return nil, err
	}
	t.graph = model.Graph
	t.session = model.Session
	return t, nil
}

func (t *tensorflow) Close() error {
	return t.session.Close()
}

func (t *tensorflow) run(inputs ...string) ([]*tf.Tensor, error) {
	if len(inputs) != len(t.feeds) {
		return nil, errors.ErrInputLength(len(inputs), len(t.feeds))
	}

	feeds := make(map[tf.Output]*tf.Tensor, len(inputs))
	for i, val := range inputs {
		inputTensor, err := tf.NewTensor(val)
		if err != nil {
			return nil, err
		}
		feeds[t.graph.Operation(t.feeds[i].operationName).Output(t.feeds[i].outputIndex)] = inputTensor
	}

	fetches := make([]tf.Output, 0, len(t.fetches))
	for _, fetch := range t.fetches {
		fetches = append(fetches, t.graph.Operation(fetch.operationName).Output(fetch.outputIndex))
	}

	return t.session.Run(feeds, fetches, t.operations)
}

func (t *tensorflow) GetVector(inputs ...string) ([]float64, error) {
	tensors, err := t.run(inputs...)
	if err != nil {
		return nil, err
	}
	if tensors == nil || tensors[0] == nil || tensors[0].Value() == nil {
		return nil, errors.ErrNilTensorTF(tensors)
	}

	switch t.ndim {
	case TwoDim:
		value, ok := tensors[0].Value().([][]float64)
		if ok {
			if value == nil {
				return nil, errors.ErrNilTensorValueTF(value)
			}
			return value[0], nil
		} else {
			return nil, errors.ErrFailedToCastTF(tensors[0].Value())
		}
	case ThreeDim:
		value, ok := tensors[0].Value().([][][]float64)
		if ok {
			if value == nil || value[0] == nil {
				return nil, errors.ErrNilTensorValueTF(value)
			}
			return value[0][0], nil
		} else {
			return nil, errors.ErrFailedToCastTF(tensors[0].Value())
		}
	default:
		value, ok := tensors[0].Value().([]float64)
		if ok {
			return value, nil
		} else {
			return nil, errors.ErrFailedToCastTF(tensors[0].Value())
		}
	}
}

func (t *tensorflow) GetValue(inputs ...string) (interface{}, error) {
	tensors, err := t.run(inputs...)
	if err != nil {
		return nil, err
	}
	if tensors == nil || tensors[0] == nil {
		return nil, errors.ErrNilTensorTF(tensors)
	}
	return tensors[0].Value(), nil
}

func (t *tensorflow) GetValues(inputs ...string) (values []interface{}, err error) {
	tensors, err := t.run(inputs...)
	if err != nil {
		return nil, err
	}
	values = make([]interface{}, 0, len(tensors))
	for _, tensor := range tensors {
		values = append(values, tensor.Value())
	}
	return values, nil
}
