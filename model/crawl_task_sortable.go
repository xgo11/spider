package model

import "sort"

type CrawlTaskSlice []*CrawlTask

func SortTasks(tasks CrawlTaskSlice) CrawlTaskSlice {
	sort.Sort(tasks)
	return tasks
}

func (cts CrawlTaskSlice) Len() int {
	return len(cts)
}

func (cts CrawlTaskSlice) Less(i, j int) bool {
	if diffP := cts[i].Schedule.Priority - cts[j].Schedule.Priority; diffP > 0 {
		return true
	} else if diffP < 0 {
		return false
	} else { // 优先级一样的情况下比较创建时间
		if diffT := cts[i].CreateTime - cts[j].CreateTime; diffT > 0 {
			return false
		} else if diffT < 0 {
			return true
		}
	}
	return false
}

func (cts CrawlTaskSlice) Swap(i, j int) {
	tmp := cts[i]
	cts[i] = cts[j]
	cts[j] = tmp
}
