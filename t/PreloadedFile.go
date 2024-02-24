package t

import "sync"

type PreloadedFile struct {
	File    string
	TmpPath string
	IsReady bool
}

type TPreloadedFiles struct {
	sync.Mutex
	List []*PreloadedFile
}

func (j *TPreloadedFiles) Append(v *PreloadedFile) {
	j.Lock()
	j.List = append(j.List, v)
	j.Unlock()
}

func (j *TPreloadedFiles) Remove(a string) {
	j.Lock()
	var i = -1
	for ii, b := range j.List {
		if b.File == a {
			i = ii
			break
		}
	}
	if i == -1 {
		return
	}
	// replace "to be deleted" with last element
	j.List[i] = j.List[len(j.List)-1]
	// return while array excluding the last element (that now sits on the to be replaced index)
	j.List = j.List[:len(j.List)-1]
	j.Unlock()
}

func (j *TPreloadedFiles) Exists(a string) bool {
	j.Lock()
	defer j.Unlock()
	for _, b := range j.List {
		if b.File == a {
			return true
		}
	}
	return false
}

func (j *TPreloadedFiles) Get() []*PreloadedFile {
	return j.List
}
