package quadtree

import (
	"fmt"

	vec "github.com/machinbrol/vecmath"
)

type Circle struct {
	C vec.Vec2 // rayon
	R float64  // centre
}

//NewCircle retourne un nouveau cercle
func NewCircle(center vec.Vec2, radius float64) *Circle {
	circ := &Circle{}
	circ.C = center
	circ.R = radius
	return circ
}

func (c *Circle) Intersect(other *Circle) bool {
	return c.R+other.R > other.C.Sub(c.C).Length()
}

func (c *Circle) Center() vec.Vec2 {
	return c.C
}

func (c *Circle) Width() float64 {
	return c.R
}

func (c *Circle) Height() float64 {
	return c.R
}

//Rectangle  Représente un rectange centré sur 'Center'.
//W et H sont les demi-largeur et demi-hauteur
type Rectangle struct {
	C  vec.Vec2 //center
	TL vec.Vec2 //top left
	W  float64  //demi-largeur
	H  float64  // demi-hauteur
}

func (r *Rectangle) Center() vec.Vec2 {
	return r.C
}

func (r *Rectangle) Width() float64 {
	return r.W
}

func (r *Rectangle) Height() float64 {
	return r.H
}

func (r *Rectangle) Intersect(o Centered) bool {
	if r == o {
		return false
	}
	return !(r.C.X+r.W < o.Center().X-o.Width() ||
		r.C.X-r.W > o.Center().X+o.Width() ||
		r.C.Y+r.H < o.Center().Y-o.Height() ||
		r.C.Y-r.H > o.Center().Y+o.Height())
}

func (r Rectangle) String() string {
	return fmt.Sprintf("{%v %v %v %v}", r.TL.X, r.TL.Y, r.W*2, r.H*2)
}

//NewRectangle Retourne un rectange dont le point top left se trouve à topLeft
// de longuer width et de hauteur height
func NewRectangle(topLeft vec.Vec2, width float64, height float64) *Rectangle {
	rect := &Rectangle{}
	rect.TL = topLeft
	rect.W = width / 2
	rect.H = height / 2
	rect.C = vec.Vec2{X: topLeft.X + width/2, Y: topLeft.Y + height/2}
	return rect
}

func NewRectangleCentered(center vec.Vec2, width float64, height float64) *Rectangle {
	rect := &Rectangle{}
	rect.C = center
	rect.TL = vec.Vec2{X: rect.C.X - width, Y: rect.C.Y - height}
	rect.W = width
	rect.H = height

	return rect
}
