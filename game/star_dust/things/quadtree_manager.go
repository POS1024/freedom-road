package things

import (
	"github.com/google/uuid"
	"math"
	"star_dust/inter"
	"sync"
)

type TreeRoute struct {
}

type Rect struct {
	X, Y, Width, Height int
}

type Quadtree struct {
	uuid            string
	level           int
	maxForks        int
	maxObjects      int
	maxLevels       int
	gameObjectsNums int
	bounds          Rect
	manager         *QuadtreeManager
	nodes           []*Quadtree
	gameObjects     map[string]inter.GameObject
}

func (q *Quadtree) GetUuid() string {
	return q.uuid
}

func (q *Quadtree) QuadrantMatching(obj inter.GameObject) (beMixed bool) {
	x, y := obj.GetXY()
	width := obj.GetWidth()
	height := obj.GetHeight()
	beMixed = true
	if q.bounds.X+q.bounds.Width <= x || x+width <= q.bounds.X {
		beMixed = false
	}
	if q.bounds.Y+q.bounds.Height <= y || y+height <= q.bounds.Y {
		beMixed = false
	}
	return
}

func (q *Quadtree) AddGameObject(obj inter.GameObject) (success bool) {
	// 判断是否存在于当前象限
	if _, ok := q.manager.gameObjectCoverTreeRoute[obj.GetUuid()][q.uuid]; ok {
		return true
	}
	// 判断是否符合当前象限
	if !q.QuadrantMatching(obj) {
		return false
	}
	// 判断是否允许继续分叉
	if q.level < q.maxLevels {
		// 判断是否分叉
		if q.nodes == nil {
			// 判断是否需要分叉
			if q.gameObjectsNums < q.maxObjects {
				q.gameObjects[obj.GetUuid()] = obj
				if treeRoute, ok := q.manager.gameObjectCoverTreeRoute[obj.GetUuid()]; ok {
					treeRoute[q.uuid] = q
					q.manager.gameObjectCoverTreeRoute[obj.GetUuid()] = treeRoute
					success = true
					q.gameObjectsNums++
				}
			} else {
				// 进行分叉
				squareRootFloat := math.Sqrt(float64(q.maxForks))
				squareRoot := int(squareRootFloat)
				unitWidth := q.bounds.Width / squareRoot
				unitHeight := q.bounds.Height / squareRoot
				for i := 0; i < squareRoot; i++ {
					for j := 0; j < squareRoot; j++ {
						childTreeX := q.bounds.X + unitWidth*i
						childTreeY := q.bounds.Y + unitHeight*j
						childTreeWidth := unitWidth
						childTreeHeight := unitHeight
						if i == squareRoot {
							childTreeWidth = q.bounds.X + q.bounds.Width - childTreeX
							childTreeHeight = q.bounds.Y + q.bounds.Height - childTreeY
						}
						q.nodes = append(q.nodes, NewQuadtree(q.maxForks, q.maxObjects, q.maxLevels, q.level+1, Rect{
							X:      childTreeX,
							Y:      childTreeY,
							Width:  childTreeWidth,
							Height: childTreeHeight,
						}, q.manager))
					}

				}

				// 旧对象取消本路由树
				for _, oldObj := range q.gameObjects {
					if treeRoute, ok := q.manager.gameObjectCoverTreeRoute[oldObj.GetUuid()]; ok {
						delete(treeRoute, q.uuid)
					}
				}

				// 提取全部需要重新分布的对象
				q.gameObjects[obj.GetUuid()] = obj
				middlewareGameObjects := q.gameObjects
				q.gameObjects = nil
				q.gameObjectsNums = 0

				// 对象推送重分布
				for _, gameObj := range middlewareGameObjects {
					innerGameObj := gameObj
					innerOk := q.AddGameObject(innerGameObj)
					if innerGameObj.GetUuid() == obj.GetUuid() && innerOk {
						success = true
					}
				}
			}
		} else {
			// 向下递归
			for _, tree := range q.nodes {
				if tree.AddGameObject(obj) {
					success = true
				}
			}
		}
	} else {
		q.gameObjects[obj.GetUuid()] = obj
		if treeRoute, ok := q.manager.gameObjectCoverTreeRoute[obj.GetUuid()]; ok {
			treeRoute[q.uuid] = q
			q.manager.gameObjectCoverTreeRoute[obj.GetUuid()] = treeRoute
			success = true
			q.gameObjectsNums++
		}
	}
	return
}

func NewQuadtree(maxForks, maxObjects, maxLevels, level int, bounds Rect, manager *QuadtreeManager) *Quadtree {
	return &Quadtree{
		maxForks:    maxForks,
		maxObjects:  maxObjects,
		maxLevels:   maxLevels,
		level:       level,
		manager:     manager,
		uuid:        uuid.New().String(),
		bounds:      bounds,
		gameObjects: make(map[string]inter.GameObject),
	}
}

type QuadtreeManager struct {
	maxForks                 int
	maxObjects               int
	maxLevels                int
	qt                       *Quadtree
	gameObjects              map[string]inter.GameObject
	gameObjectCoverTreeRoute map[string]map[string]*Quadtree
	lock                     sync.Mutex
}

func (q *QuadtreeManager) AddGameObject(obj inter.GameObject, isLock bool) bool {
	if isLock {
		q.lock.Lock()
		defer q.lock.Unlock()
	}

	// 判断是否存在
	if _, ok := q.gameObjectCoverTreeRoute[obj.GetUuid()]; ok {
		return false
	}
	q.gameObjectCoverTreeRoute[obj.GetUuid()] = make(map[string]*Quadtree, 0)

	q.gameObjects[obj.GetUuid()] = obj

	// 递归将对象放入n叉树
	if success := q.qt.AddGameObject(obj); success {
		return success
	}
	delete(q.gameObjectCoverTreeRoute, obj.GetUuid())
	delete(q.gameObjects, obj.GetUuid())
	return false
}

func (q *QuadtreeManager) RemoveGameObject(obj inter.GameObject) bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	// 判断是否存在
	if treeRouteMap, ok := q.gameObjectCoverTreeRoute[obj.GetUuid()]; ok {
		for _, tree := range treeRouteMap {
			delete(tree.gameObjects, obj.GetUuid())
			tree.gameObjectsNums = tree.gameObjectsNums - 1
		}
		delete(q.gameObjectCoverTreeRoute, obj.GetUuid())
		delete(q.gameObjects, obj.GetUuid())
		return true
	}
	return false

}

func (q *QuadtreeManager) MoveGameObject(obj inter.GameObject) bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	// 判断是否存在
	if treeRouteMap, ok := q.gameObjectCoverTreeRoute[obj.GetUuid()]; ok {
		for _, tree := range treeRouteMap {
			if !tree.QuadrantMatching(obj) {
				delete(tree.gameObjects, obj.GetUuid())
				tree.gameObjectsNums = tree.gameObjectsNums - 1
				delete(treeRouteMap, tree.GetUuid())
			}
		}
	}

	if success := q.qt.AddGameObject(obj); success {
		return success
	}
	return false
}

func (q *QuadtreeManager) SameAreaGameObjects(obj inter.GameObject) (objs map[string]inter.GameObject) {
	q.lock.Lock()
	defer q.lock.Unlock()
	objs = make(map[string]inter.GameObject, 0)

	// 判断是否存在
	if treeRouteMap, ok := q.gameObjectCoverTreeRoute[obj.GetUuid()]; ok {
		// 遍历所有的相关树
		for _, tree := range treeRouteMap {
			for _, otherObj := range tree.gameObjects {
				if otherObj.GetUuid() != obj.GetUuid() {
					objs[otherObj.GetUuid()] = otherObj
				}
			}
		}
	}
	return
}

func (q *QuadtreeManager) MoveAreaGameObjects(offsetX, offsetY int) {

	for _, obj := range q.gameObjects {
		innerObj := obj
		x, y := innerObj.GetXY()
		innerObj.SetXY(x+offsetX, y+offsetY)
		q.MoveGameObject(innerObj)
	}
	return
}

func (q *QuadtreeManager) PixelCollision(obj1, obj2 inter.GameObject) bool {

	obj1X, obj1Y := obj1.GetXY()
	obj1W := obj1.GetWidth()
	obj1Image := obj1.GetImage()
	obj1PictureFrame := obj1.GetPictureFrame()
	obj2X, obj2Y := obj2.GetXY()
	obj2W := obj2.GetWidth()
	obj2Image := obj2.GetImage()
	obj2PictureFrame := obj2.GetPictureFrame()

	w1, h1 := obj1Image.Size()
	w2, h2 := obj2Image.Size()

	for i := 0; i < w1; i++ {
		for j := 0; j < h1; j++ {
			c1 := obj1Image.At(i+obj1PictureFrame*obj1W, j)
			_, _, _, a1 := c1.RGBA()
			if a1 == 0 {
				continue
			}
			x2i := i + obj1X - obj2X
			y2j := j + obj1Y - obj2Y

			if x2i < 0 || x2i >= w2 || y2j < 0 || y2j >= h2 {
				continue
			}

			c2 := obj2Image.At(x2i+obj2PictureFrame*obj2W, y2j)
			_, _, _, a2 := c2.RGBA()
			if a2 != 0 {
				return true
			}
		}
	}
	return false
}

func NewQuadtreeManager(maxForks, maxObjects, maxLevels int, bounds Rect) *QuadtreeManager {
	q := &QuadtreeManager{
		maxForks:                 maxForks,
		maxObjects:               maxObjects,
		maxLevels:                maxLevels,
		gameObjectCoverTreeRoute: make(map[string]map[string]*Quadtree, 0),
		gameObjects:              make(map[string]inter.GameObject, 0),
	}
	q.qt = NewQuadtree(maxForks, maxObjects, maxLevels, 1, bounds, q)
	return q
}

var QM = NewQuadtreeManager(4, 100, 10, Rect{X: 0, Y: 0, Width: inter.ScreenWidth, Height: inter.ScreenHeight})
