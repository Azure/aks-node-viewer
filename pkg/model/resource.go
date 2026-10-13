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
	"strconv"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// formatResource renders a resource quantity in the most readable unit for the
// resource instead of the canonical Kubernetes representation. CPU is shown in
// cores (millicores only below a single core) and memory/storage in binary
// units (Ki/Mi/Gi/Ti/Pi). Other resources keep their canonical form.
func formatResource(res v1.ResourceName, q resource.Quantity) string {
	switch res {
	case v1.ResourceCPU:
		return formatCPU(q)
	case v1.ResourceMemory, v1.ResourceEphemeralStorage, v1.ResourceStorage:
		return formatBytes(q)
	default:
		return q.String()
	}
}

// formatCPU renders CPU in cores, e.g. "4", "1.5" or "4,250". Values below a
// single core keep the familiar millicore form, e.g. "500m".
func formatCPU(q resource.Quantity) string {
	milli := q.MilliValue()
	switch {
	case milli == 0:
		return "0"
	case milli < 1000:
		return fmt.Sprintf("%dm", milli)
	default:
		return groupThousands(trimFloat(float64(milli)/1000, 3))
	}
}

// formatBytes renders a byte quantity using binary units, e.g. "16Gi" or
// "1.5Ti".
func formatBytes(q resource.Quantity) string {
	units := []struct {
		suffix string
		scale  float64
	}{
		{"Pi", 1 << 50},
		{"Ti", 1 << 40},
		{"Gi", 1 << 30},
		{"Mi", 1 << 20},
		{"Ki", 1 << 10},
	}
	b := q.AsApproximateFloat64()
	for _, u := range units {
		if b >= u.scale {
			return groupThousands(trimFloat(b/u.scale, 2)) + u.suffix
		}
	}
	return groupThousands(strconv.FormatInt(q.Value(), 10))
}

// formatCount renders an integer with thousands separators, e.g. "1234"
// becomes "1,234".
func formatCount(n int) string {
	return groupThousands(strconv.FormatInt(int64(n), 10))
}

// groupThousands inserts ',' separators into the integer part of a decimal
// number string, e.g. "4250" becomes "4,250" and "12345.5" becomes "12,345.5".
func groupThousands(s string) string {
	neg := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")
	intPart, frac, hasFrac := strings.Cut(s, ".")
	var b strings.Builder
	for i := 0; i < len(intPart); i++ {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteByte(intPart[i])
	}
	out := b.String()
	if hasFrac {
		out += "." + frac
	}
	if neg {
		return "-" + out
	}
	return out
}

// trimFloat formats v with at most prec decimals and drops any trailing zeros,
// so 16.00 becomes "16" and 1.50 becomes "1.5".
func trimFloat(v float64, prec int) string {
	s := strconv.FormatFloat(v, 'f', prec, 64)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}
