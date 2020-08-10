// Package util defines some utility/helper types and functions.
package util

// Pow2 returns the first power-of-two value >= to n.
// This can be used to create suitable texture dimensions.
func Pow2(n int) int {
	x := uint32(n - 1)
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	return int(x + 1)
}

// Min returns a iff a < b. otherwise returns b.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns a iff a > b. otherwise returns b.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Clampi returns v, clamped to the range [min, max].
func Clampi(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// Mat4 defines a 4x4 matrix.
type Mat4 [16]float32

// Copy returns a copy of m.
func (m *Mat4) Copy() *Mat4 {
	var n Mat4
	copy(n[:], (*m)[:])
	return &n
}

// Identity sets m to the identity matrix.
func (m *Mat4) Identity() {
	mm := *m
	mm[0] = 1
	mm[1] = 0
	mm[2] = 0
	mm[3] = 0
	mm[4] = 0
	mm[5] = 1
	mm[6] = 0
	mm[7] = 0
	mm[8] = 0
	mm[9] = 0
	mm[10] = 1
	mm[11] = 0
	mm[12] = 0
	mm[13] = 0
	mm[14] = 0
	mm[15] = 1
	*m = mm
}

// Mul sets ma to the multiplication of ma with mb.
func (ma *Mat4) Mul(mb *Mat4) {
	a, b := *ma, *mb
	a0 := a[0]
	a1 := a[1]
	a2 := a[2]
	a3 := a[3]
	a4 := a[4]
	a5 := a[5]
	a6 := a[6]
	a7 := a[7]
	a8 := a[8]
	a9 := a[9]
	a10 := a[10]
	a11 := a[11]
	a12 := a[12]
	a13 := a[13]
	a14 := a[14]
	a15 := a[15]

	a[0] = a0*b[0] + a4*b[1] + a8*b[2] + a12*b[3]
	a[1] = a1*b[0] + a5*b[1] + a9*b[2] + a13*b[3]
	a[2] = a2*b[0] + a6*b[1] + a10*b[2] + a14*b[3]
	a[3] = a3*b[0] + a7*b[1] + a11*b[2] + a15*b[3]

	a[4] = a0*b[4] + a4*b[5] + a8*b[6] + a12*b[7]
	a[5] = a1*b[4] + a5*b[5] + a9*b[6] + a13*b[7]
	a[6] = a2*b[4] + a6*b[5] + a10*b[6] + a14*b[7]
	a[7] = a3*b[4] + a7*b[5] + a11*b[6] + a15*b[7]

	a[8] = a0*b[8] + a4*b[9] + a8*b[10] + a12*b[11]
	a[9] = a1*b[8] + a5*b[9] + a9*b[10] + a13*b[11]
	a[10] = a2*b[8] + a6*b[9] + a10*b[10] + a14*b[11]
	a[11] = a3*b[8] + a7*b[9] + a11*b[10] + a15*b[11]

	a[12] = a0*b[12] + a4*b[13] + a8*b[14] + a12*b[15]
	a[13] = a1*b[12] + a5*b[13] + a9*b[14] + a13*b[15]
	a[14] = a2*b[12] + a6*b[13] + a10*b[14] + a14*b[15]
	a[15] = a3*b[12] + a7*b[13] + a11*b[14] + a15*b[15]

	*ma = a
}

// Mat4Ortho returns the orthographic projection for the given viewport.
func Mat4Ortho(left, right, top, bottom, znear, zfar float32) *Mat4 {
	var s Mat4
	s.Identity()

	rml := right - left
	rpl := right + left
	tmb := top - bottom
	tpb := top + bottom
	fmn := zfar - znear
	fpn := zfar + znear

	s[0] = 2.0 / rml
	s[5] = 2.0 / tmb
	s[10] = -2.0 / fmn
	s[12] = -rpl / rml
	s[13] = -tpb / tmb
	s[14] = -fpn / fmn
	s[15] = 1.0
	return &s
}

// Mat4Scale returns the scale matrix for dimensions x/y/z.
func Mat4Scale(x, y, z float32) *Mat4 {
	var s Mat4
	s.Identity()
	s[0] = x
	s[5] = y
	s[10] = z
	s[15] = 1
	return &s
}

// Translate returns the translation matrix for coordinates x/y/z.
func Mat4Translate(x, y, z float32) *Mat4 {
	var s Mat4
	s[0] = 1
	s[5] = 1
	s[10] = 1
	s[12] = x
	s[13] = y
	s[14] = z
	s[15] = 1
	return &s
}
