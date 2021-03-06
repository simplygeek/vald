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

// Package config providers configuration type and load configuration logic
package config

type AgentSidecar struct {
	// Name string `yaml:"name" json:"name"`

	// WatchPaths represents watch path list for backup
	WatchPaths []string `yaml:"watch_paths" json:"watch_paths"`

	// AutoBackupDurationLimit represents auto backup duration limit
	AutoBackupDurationLimit string `yaml:"auto_backup_duration_limit" json:"auto_backup_duration_limit"`

	// AutoBackupDuration represent checking loop duration for auto backup execution
	AutoBackupDuration string `yaml:"auto_backup_duration" json:"auto_backup_duration"`
}

func (s *AgentSidecar) Bind() *AgentSidecar {
	s.WatchPaths = GetActualValues(s.WatchPaths)
	s.AutoBackupDuration = GetActualValue(s.AutoBackupDuration)
	s.AutoBackupDurationLimit = GetActualValue(s.AutoBackupDurationLimit)
	return s
}
