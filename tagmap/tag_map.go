package tagmap

import (
	"sort"
)

const TM_VERSION = "1"

type TM struct {
	Version   string        `json:"version"`
	FileToTag StringListMap `json:"fileToTag"`
	TagToFile StringListMap `json:"tagToFile"`
}

func New() *TM {
	return &TM{
		Version: TM_VERSION,
		FileToTag: make(StringListMap),
		TagToFile: make(StringListMap),
	}
}

func (tm *TM) ListFiles() []string {
	return tm.FileToTag.Keys()
}

func (tm *TM) FilesFor(tags ...string) []string {
	fileSet := make(map[string]bool)

	for _, tag := range tags {
		files, ok := tm.TagToFile[tag]
		if !ok {
			continue
		}

		for _, f := range files {
			fileSet[f] = true
		}
	}

	count := len(fileSet)
	fileList := make([]string, count)
	i := 0
	for f := range fileSet {
		fileList[i] = f
		i++
	}

	sort.Strings(fileList)
	return fileList
}

func (tm *TM) Add(file, tag string) {
	tm.FileToTag.AddUnique(file, tag)
	tm.TagToFile.AddUnique(tag, file)
}

func (tm *TM) Remove(file, tag string) bool {
	found1 := tm.FileToTag.RemoveFirst(file, tag)
	found2 := tm.TagToFile.RemoveFirst(tag, file)
	return found1 || found2
}

func (tm *TM) Clear(file string) {
	tags, ok := tm.FileToTag[file]
	if !ok {
		return
	}

	delete(tm.FileToTag, file)

	for _, tag := range tags {
		tm.TagToFile.RemoveFirst(tag, file)
	}
}

func (tm *TM) Normalize() *TM {
	tm = &(*tm) // clone

	tm.Version = TM_VERSION

	for _, tags := range tm.FileToTag {
		sort.Strings(tags)
	}

	return tm
}
