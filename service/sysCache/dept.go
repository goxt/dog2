package sysCache

import (
	"config/models"
	"github.com/goxt/dog2/util"
	"sync"
)

const DeptRootId = 0

var deptSync sync.Mutex // 部门的更新锁
var deptList []Dept     // 缓存的部门列表
var deptTree []Dept     // 缓存的部门树

type Dept struct {
	DeptId    uint64
	DeptPid   uint64
	DeptType  uint8
	DeptPtype uint8
	DeptName  string
	Sort      uint64
	models.Base
	models.DeptExt
	Children []Dept
}

func ResetDeptCache() {

	deptSync.Lock()
	defer deptSync.Unlock()

	db := util.OpenDbConnection()
	begin := false
	defer util.CloseDbConnection(db, &begin)

	var datum []Dept
	rst := db.Table(models.Dept{}.TableName()).Order("sort ASC").Scan(&datum)
	util.IsEmpty(rst)

	deptList = make([]Dept, len(datum), cap(datum))
	copy(deptList, datum)

	deptTree = DeptTree(datum)
}

/**
 * 从缓存中获取部门列表 - 数组形式
 * @param	all		是否获取所有部门（包括已删除的部门），默认值为false
 */
func DeptList(all ...bool) []Dept {

	if deptList == nil {
		ResetDeptCache()
	}

	allFlag := false
	if len(all) > 0 {
		allFlag = all[0]
	}

	if allFlag {
		datum := make([]Dept, len(deptList), cap(deptList))
		copy(datum, deptList)
		return datum
	} else {
		var tmp []Dept
		for _, v := range deptList {
			if v.DeletedAt != nil {
				continue
			}
			tmp = append(tmp, v)
		}
		datum := make([]Dept, len(tmp), cap(tmp))
		copy(datum, tmp)
		return datum
	}
}

/**
 * 从缓存中获取指定父级ID下的部门列表 - 数组形式
 * @param	pid		父级ID
 * @param	all		是否获取所有部门（包括已删除的部门），默认值为false
 */
func DeptListByPid(pid uint64, all ...bool) []Dept {
	if deptList == nil {
		ResetDeptCache()
	}

	allFlag := false
	if len(all) > 0 {
		allFlag = all[0]
	}

	temp := []Dept{}

	for _,v := range deptList {
		if v.DeptPid != pid {
			continue
		}
		if v.DeletedAt != nil && !allFlag {
			continue
		}
		temp = append(temp, v)
	}

	datum := make([]Dept, len(temp), cap(temp))
	copy(datum, temp)
	return datum
}

/**
 * 从缓存中获取部门列表 - map形式
 * @param	all		是否获取所有部门（包括已删除的部门），默认值为false
 */
func DeptMap(all ...bool) map[uint64]*Dept {

	allFlag := false
	if len(all) > 0 {
		allFlag = all[0]
	}
	arr := DeptList(allFlag)

	var data = map[uint64]*Dept{}

	for k, v := range arr {
		data[v.DeptId] = &arr[k]
	}

	return data
}

/**
 * 将部门列表转成树形结构
 * @param	list	部门列表
 */
func DeptTree(list []Dept) []Dept {

	if list == nil || len(list) == 0 {
		return nil
	}

	// 初始化变量
	var tree []Dept
	sp := map[uint64]uint64{}
	ps := map[uint64][]uint64{}
	var m = map[uint64]*Dept{}
	var root []uint64
	var ids []uint64

	// 初始化数据
	for k, v := range list {
		v.Children = []Dept{}
		m[v.DeptId] = &list[k]
		sp[v.DeptId] = v.DeptPid
		ids = append(ids, v.DeptId)
		ps[v.DeptPid] = append(ps[v.DeptPid], v.DeptId)
	}

	// 找出根节点
	for _, v := range list {
		id := v.DeptId
		pid := sp[id]
		if !util.InArrayUint64(pid, ids) {
			root = append(root, id)
		}
	}

	// 列出push顺序
	arr := recursionPushArr(root, ps)

	// 处理节点
	for i := len(arr) - 1; i >= 0; i-- {
		id := arr[i]
		if util.InArrayUint64(id, root) {
			continue
		}

		pid := sp[id]
		m[pid].Children = append(m[pid].Children, *(m[id]))
	}

	// 将根节点复制到最终树形结构即可
	for _, v := range root {
		tree = append(tree, *(m[v]))
	}

	return tree
}

/**
 * 获取缓存中的树形结构（所有部门）
 */
func DeptTreeDefault() []Dept {
	datum := make([]Dept, len(deptTree), cap(deptTree))
	copy(datum, deptTree)
	return datum
}

func recursionPushArr(arr []uint64, ps map[uint64][]uint64) []uint64 {

	var datum []uint64
	for i := len(arr) - 1; i >= 0; i-- {
		v := arr[i]
		datum = append(datum, v)
		tmp := recursionPushArr(ps[v], ps)
		for _, vv := range tmp {
			datum = append(datum, vv)
		}
	}

	return datum
}
