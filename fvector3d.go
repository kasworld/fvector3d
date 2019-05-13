// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package fvector3d : float 3d vector
package fvector3d

import (
	"fmt"
	"math"
	"math/rand"
)

type Vt [3]float64

func (v Vt) String() string {
	return fmt.Sprintf("[%5.2f,%5.2f,%5.2f]", v[0], v[1], v[2])
}

var Zero = Vt{0, 0, 0}
var UnitX = Vt{1, 0, 0}
var UnitY = Vt{0, 1, 0}
var UnitZ = Vt{0, 0, 1}

// func (p Vt) Copy() Vt {
// 	return Vt{p[0], p[1], p[2]}
// }
func (p Vt) Eq(other Vt) bool {
	return p == other
	//return p[0] == other[0] && p[1] == other[1] && p[2] == other[2]
}
func (p Vt) Ne(other Vt) bool {
	return !p.Eq(other)
}
func (p Vt) IsZero() bool {
	return p.Eq(Zero)
}
func (p Vt) Add(other Vt) Vt {
	return Vt{p[0] + other[0], p[1] + other[1], p[2] + other[2]}
}
func (p Vt) Neg() Vt {
	return Vt{-p[0], -p[1], -p[2]}
}
func (p Vt) Sub(other Vt) Vt {
	return Vt{p[0] - other[0], p[1] - other[1], p[2] - other[2]}
}
func (p Vt) Mul(other Vt) Vt {
	return Vt{p[0] * other[0], p[1] * other[1], p[2] * other[2]}
}
func (p Vt) Imul(other float64) Vt {
	return Vt{p[0] * other, p[1] * other, p[2] * other}
}
func (p Vt) Idiv(other float64) Vt {
	return Vt{p[0] / other, p[1] / other, p[2] / other}
}
func (p Vt) Abs() float64 {
	return math.Sqrt(p[0]*p[0] + p[1]*p[1] + p[2]*p[2])
}
func (p Vt) Sqd(q Vt) float64 {
	var sum float64
	for dim, pCoord := range p {
		d := pCoord - q[dim]
		sum += d * d
	}
	return sum
}

func (p Vt) LenTo(other Vt) float64 {
	return math.Sqrt(p.Sqd(other))
}

func (p *Vt) Normalize() {
	d := p.Abs()
	if d > 0 {
		p[0] /= d
		p[1] /= d
		p[2] /= d
	}
}
func (p Vt) Normalized() Vt {
	d := p.Abs()
	if d > 0 {
		return p.Idiv(d)
	}
	return p
}
func (p Vt) NormalizedTo(l float64) Vt {
	d := p.Abs() / l
	if d != 0 {
		return p.Idiv(d)
	}
	return p
}
func (p Vt) Dot(other Vt) float64 {
	return p[0]*other[0] + p[1]*other[1] + p[2]*other[2]
}
func (p Vt) Cross(other Vt) Vt {
	return Vt{
		p[1]*other[2] - p[2]*other[1],
		-p[0]*other[2] + p[2]*other[0],
		p[0]*other[1] - p[1]*other[0],
	}
}

// reflect plane( == normal vector )
func (p Vt) Reflect(normal Vt) Vt {
	d := 2 * (p[0]*normal[0] + p[1]*normal[1] + p[2]*normal[2])
	return Vt{p[0] - d*normal[0], p[1] - d*normal[1], p[2] - d*normal[2]}
}
func (p Vt) RotateAround(axis Vt, theta float64) Vt {
	// Return the vector rotated around axis through angle theta. Right hand rule applies
	// Adapted from equations published by Glenn Murray.
	// http://inside.mines.edu/~gmurray/ArbitraryAxisRotation/ArbitraryAxisRotation.html
	x, y, z := p[0], p[1], p[2]
	u, v, w := axis[0], axis[1], axis[2]

	// Extracted common factors for simplicity and efficiency
	r2 := u*u + v*v + w*w
	r := math.Sqrt(r2)
	ct := math.Cos(theta)
	st := math.Sin(theta) / r
	dt := (u*x + v*y + w*z) * (1 - ct) / r2
	return Vt{
		(u*dt + x*ct + (-w*y+v*z)*st),
		(v*dt + y*ct + (w*x-u*z)*st),
		(w*dt + z*ct + (-v*x+u*y)*st),
	}
}
func (p Vt) Angle(other Vt) float64 {
	// Return the angle to the vector other
	dot := p.Dot(other)
	l := (p.Abs() * other.Abs())
	rtn := math.Acos(dot / l)
	// log.Info("%v %v %v %v %v", p, other, dot, l, rtn)
	return rtn
}
func (p Vt) Project(other Vt) Vt {
	// Return one vector projected on the vector other
	n := other.Normalized()
	return n.Imul(p.Dot(n))
}

// for aim ahead target with projectile
// return time dur
func (srcpos Vt) CalcAimAheadDur(dstpos Vt, dstmv Vt, bulletspeed float64) float64 {
	totargetvt := dstpos.Sub(srcpos)
	a := dstmv.Dot(dstmv) - bulletspeed*bulletspeed
	b := 2 * dstmv.Dot(totargetvt)
	c := totargetvt.Dot(totargetvt)
	p := -b / (2 * a)
	q := math.Sqrt((b*b)-4*a*c) / (2 * a)
	t1 := p - q
	t2 := p + q

	var rtn float64
	if t1 > t2 && t2 > 0 {
		rtn = t2
	} else {
		rtn = t1
	}
	if rtn < 0 || math.IsNaN(rtn) {
		return math.Inf(1)
	}
	return rtn
}

// for serialize
func (v Vt) NewInt32Vector() [3]int32 {
	return [3]int32{int32(v[0]), int32(v[1]), int32(v[2])}
}

func FromInt32Vector(s [3]int32) Vt {
	return Vt{float64(s[0]), float64(s[1]), float64(s[2])}
}

func RandVt(st, end float64) Vt {
	return Vt{
		rand.Float64()*(end-st) + st,
		rand.Float64()*(end-st) + st,
		rand.Float64()*(end-st) + st,
	}
}

func RandVector(st, end Vt) Vt {
	return Vt{
		rand.Float64()*(end[0]-st[0]) + st[0],
		rand.Float64()*(end[1]-st[1]) + st[1],
		rand.Float64()*(end[2]-st[2]) + st[2],
	}
}

func (center Vt) To8Direct(v2 Vt) int {
	rtn := 0
	for i := 0; i < 3; i++ {
		if center[i] > v2[i] {
			rtn += 1 << uint(i)
		}
	}
	return rtn
}
