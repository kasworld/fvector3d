// Copyright 2015 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package fvector3d

import (
	"math/rand"
)

func (h *Cube) MakeCubeBy8Driect(center Vt, direct8 int) *Cube {
	rtn := Vt{}
	for i := 0; i < 3; i++ {
		if direct8&(1<<uint(i)) != 0 {
			rtn[i] = h.Min[i]
		} else {
			rtn[i] = h.Max[i]
		}
	}
	return NewHyperRect(center, rtn)
}

type Cube struct {
	Min, Max Vt
}

func (h *Cube) Center() Vt {
	return h.Min.Add(h.Max).Idiv(2)
}

func (h *Cube) DiagLen() float64 {
	return h.Min.LenTo(h.Max)
}

func (h *Cube) SizeVector() Vt {
	return h.Max.Sub(h.Min)
}

func (h *Cube) IsContact(c Vt, r float64) bool {
	hc := h.Center()
	hl := h.DiagLen()
	return hl/2+r >= hc.LenTo(c)
}

func NewHyperRectByCR(c Vt, r float64) *Cube {
	return &Cube{
		Vt{c[0] - r, c[1] - r, c[2] - r},
		Vt{c[0] + r, c[1] + r, c[2] + r},
	}
}

func (h *Cube) RandVector() Vt {
	return Vt{
		rand.Float64()*(h.Max[0]-h.Min[0]) + h.Min[0],
		rand.Float64()*(h.Max[1]-h.Min[1]) + h.Min[1],
		rand.Float64()*(h.Max[2]-h.Min[2]) + h.Min[2],
	}
}

func (h *Cube) Move(v Vt) *Cube {
	return &Cube{
		Min: h.Min.Add(v),
		Max: h.Max.Add(v),
	}
}

func (h *Cube) IMul(i float64) *Cube {
	hs := h.SizeVector().Imul(i / 2)
	hc := h.Center()
	return &Cube{
		Min: hc.Sub(hs),
		Max: hc.Add(hs),
	}
}

// make normalized hyperrect , if not need use Cube{Min: , Max:}
func NewHyperRect(v1 Vt, v2 Vt) *Cube {
	rtn := Cube{
		Min: Vt{},
		Max: Vt{},
	}
	for i := 0; i < 3; i++ {
		if v1[i] > v2[i] {
			rtn.Max[i] = v1[i]
			rtn.Min[i] = v2[i]
		} else {
			rtn.Max[i] = v2[i]
			rtn.Min[i] = v1[i]
		}
	}
	return &rtn
}

func (h1 *Cube) IsOverlap(h2 *Cube) bool {
	return !((h1.Min[0] > h2.Max[0] || h1.Max[0] < h2.Min[0]) ||
		(h1.Min[1] > h2.Max[1] || h1.Max[1] < h2.Min[1]) ||
		(h1.Min[2] > h2.Max[2] || h1.Max[2] < h2.Min[2]))

	// for i := 0; i < 3; i++ {
	// 	if !between(h1.Min[i], h1.Max[i], h2.Min[i]) && !between(h1.Min[i], h1.Max[i], h2.Max[i]) {
	// 		return false
	// 	}
	// }
	// return true
}

func (h1 *Cube) IsIn(h2 *Cube) bool {
	for i := 0; i < 3; i++ {
		if h1.Min[i] < h2.Min[i] || h1.Max[i] > h2.Max[i] {
			return false
		}
	}
	return true
}

func (p Vt) IsIn(hr *Cube) bool {
	return hr.Min[0] <= p[0] && p[0] <= hr.Max[0] &&
		hr.Min[1] <= p[1] && p[1] <= hr.Max[1] &&
		hr.Min[2] <= p[2] && p[2] <= hr.Max[2]
	// for i := 0; i < 3; i++ {
	// 	if hr.Min[i] > p[i] || hr.Max[i] < p[i] {
	// 		return false
	// 	}
	// }
	// return true
}

func (p *Vt) MakeIn(hr *Cube) int {
	changed := 0
	var i uint
	for i = 0; i < 3; i++ {
		if p[i] > hr.Max[i] {
			p[i] = hr.Max[i]
			changed += 1 << (i*2 + 1)
		}
		if p[i] < hr.Min[i] {
			p[i] = hr.Min[i]
			changed += 1 << (i * 2)
		}
	}
	return changed
}
