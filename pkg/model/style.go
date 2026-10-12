/*
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

package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
)

// Style holds the three colors used to style 'good', 'ok' and 'bad' values.
// The first color is used for 'good', the second for 'ok' and the third for
// 'bad'. The gradient displayed in the progress bars runs from 'bad' to 'good'.
type Style struct {
	green    func(str string) string
	yellow   func(str string) string
	red      func(str string) string
	gradient progress.Option
}

// ParseStyle parses a comma separated list of three colors, e.g. "#04B575,#FFFF00,#FF0000".
func ParseStyle(style string) (*Style, error) {
	colors := strings.Split(style, ",")
	if len(colors) != 3 {
		return nil, fmt.Errorf("three colors must be provided for the style, found %d (%q)", len(colors), style)
	}
	s := &Style{}
	s.green = lipgloss.NewStyle().Foreground(lipgloss.Color(colors[0])).Render
	s.yellow = lipgloss.NewStyle().Foreground(lipgloss.Color(colors[1])).Render
	s.red = lipgloss.NewStyle().Foreground(lipgloss.Color(colors[2])).Render

	s.gradient = progress.WithGradient(colors[2], colors[0])
	return s, nil
}
