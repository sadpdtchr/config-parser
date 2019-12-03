/*
Copyright 2019 HAProxy Technologies

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package actions

import (
	"fmt"
)

type TrackSc1 struct {
	Key   string
	Table string
}

func (f *TrackSc1) Parse(parts []string) error {

	if len(parts) <= 1 {
		return fmt.Errorf("not enough params")
	}

	f.Key = parts[1]

	if len(parts) == 3 {
		f.Table = parts[2]
	}

	return nil
}

func (f *TrackSc1) String() string {

	if f.Table != "" {
		return fmt.Sprintf("track-sc1 %s %s", f.Key, f.Table)
	}

	return fmt.Sprintf("track-sc1 %s", f.Key)
}
