package quadtree

import (
	vec "github.com/etic4/vecmath"
	rl "github.com/gen2brain/raylib-go/raylib"
)

//Centered ...
type Centered interface {
	Center() vec.Vec2
	Width() float64
	Height() float64
	Intersect(Centered) bool
}

//Quadtree ...
type Quadtree struct {
	points    []Centered
	maxPoints int
	divided   bool
	rect      *Rectangle
	ne        *Quadtree
	nw        *Quadtree
	se        *Quadtree
	sw        *Quadtree
	root      *Quadtree
	pointsMap map[vec.Vec2][]*Quadtree // stocke la liste des quadtrees dans lesquels se trouvent
	// chaque forme
}

//NewQuadtree retourne un nouveau Quadtree.
func NewQuadtree(rect *Rectangle, maxPoints int, root *Quadtree) *Quadtree {
	q := &Quadtree{}
	q.rect = rect
	q.maxPoints = maxPoints
	q.root = root
	if root == nil {
		q.root = q
		q.pointsMap = map[vec.Vec2][]*Quadtree{}
	}
	return q
}

//Size ...
func (q *Quadtree) Size() int {
	sze := 0
	sze += len(q.points)
	if q.divided {
		sze += q.ne.Size()
		sze += q.nw.Size()
		sze += q.se.Size()
		sze += q.sw.Size()
	}
	return sze
}

//GetQuadtreesFor retourne la liste des Quadtrees auxquels appartient 'c'
func (q *Quadtree) GetQuadtreesFor(v vec.Vec2) []*Quadtree {
	return q.pointsMap[v]
}

//Insert ...
func (q *Quadtree) Insert(c Centered) {
	if !q.Intersect(c) {
		return
	}

	if len(q.points) < q.maxPoints {
		q.points = append(q.points, c)
		_, ok := q.root.pointsMap[c.Center()]
		if !ok {
			q.root.pointsMap[c.Center()] = []*Quadtree{}
		}
		q.root.pointsMap[c.Center()] = append(q.root.pointsMap[c.Center()], q)
	} else {
		if !q.divided {
			W2 := q.rect.W / 2
			H2 := q.rect.H / 2

			rect := NewRectangleCentered(q.rect.C.Add(vec.Vec2{X: -W2, Y: -H2}), W2, H2)
			q.ne = NewQuadtree(rect, 4, q.root)

			rect = NewRectangleCentered(q.rect.C.Add(vec.Vec2{X: W2, Y: -H2}), W2, H2)
			q.nw = NewQuadtree(rect, 4, q.root)

			rect = NewRectangleCentered(q.rect.C.Add(vec.Vec2{X: W2, Y: H2}), W2, H2)
			q.se = NewQuadtree(rect, 4, q.root)

			rect = NewRectangleCentered(q.rect.C.Add(vec.Vec2{X: -W2, Y: H2}), W2, H2)
			q.sw = NewQuadtree(rect, 4, q.root)

			q.divided = true
		}
		q.ne.Insert(c)
		q.nw.Insert(c)
		q.se.Insert(c)
		q.sw.Insert(c)
	}
}

//Remove ...
func (q *Quadtree) Remove(c Centered) {
	if qtrees, ok := q.root.pointsMap[c.Center()]; ok {
		for _, qtree := range qtrees {
			i := 0
			length := len(qtree.points)
			for i < length && qtree.points[i].Center() != c.Center() {
				i++
			}
			if i < length {
				qtree.points[i] = qtree.points[length-1]
				qtree.points[length-1] = nil
				qtree.points = qtree.points[:length-1]
			}
		}
		delete(q.pointsMap, c.Center())
	}
}

//Clear ...
func (q *Quadtree) Clear() {
	q.points = nil
	q.points = []Centered{}
	if q.divided {
		q.ne.Clear()
		q.nw.Clear()
		q.se.Clear()
		q.sw.Clear()
	}
	q.root.pointsMap = nil
	q.root.pointsMap = map[vec.Vec2][]*Quadtree{}
}

//Intersect ...
func (q *Quadtree) Intersect(r Centered) bool {
	rect := q.rect
	return !(rect.Center().X+rect.W < r.Center().X-r.Width() ||
		rect.Center().X-rect.W > r.Center().X+r.Width() ||
		rect.Center().Y+rect.H < r.Center().Y-r.Height() ||
		rect.Center().Y-rect.H > r.Center().Y+r.Height())
}

//QueryRange ...
func (q *Quadtree) QueryRange(rect Centered) []Centered {
	res := []Centered{}
	if q.Intersect(rect) {
		res = append(res, q.points...) //inutile je pense de checker s'ils intersectent avec rect

		if q.divided {
			res = append(res, q.ne.QueryRange(rect)...)
			res = append(res, q.nw.QueryRange(rect)...)
			res = append(res, q.se.QueryRange(rect)...)
			res = append(res, q.sw.QueryRange(rect)...)
		}
	}
	return res
}

//Draw ...
func (q *Quadtree) Draw() {
	pos := q.rect.Center().Sub(vec.Vec2{X: q.rect.W, Y: q.rect.H})
	rl.DrawRectangleLines(int32(pos.X), int32(pos.Y), int32(q.rect.W*2), int32(q.rect.H*2), rl.White)

	if q.divided {
		q.ne.Draw()
		q.nw.Draw()
		q.se.Draw()
		q.sw.Draw()
	}
}

//DrawOne ...
func (q *Quadtree) DrawOne() {
	pos := q.rect.Center().Sub(vec.Vec2{X: q.rect.W, Y: q.rect.H})
	rl.DrawRectangleLines(int32(pos.X), int32(pos.Y), int32(q.rect.W*2), int32(q.rect.H*2), rl.White)
}
