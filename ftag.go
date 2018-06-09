package main

import (
	"fmt"
	"github.com/troykinsella/ftag/tagmap"
	"os"
	"sort"
)

type FTag struct {
	tagMapStore tagmap.Store

	tagMap *tagmap.TM
}

func New(tagMapStore tagmap.Store) *FTag {
	if tagMapStore == nil {
		panic("tagMapStore required")
	}

	return &FTag{
		tagMapStore: tagMapStore,
	}
}

func (ft *FTag) LoadTagMap() error {
	var err error
	ft.tagMap, err = ft.tagMapStore.Load()
	if err != nil {
		return err
	}

	return err
}

func (ft *FTag) StoreTagMap() error {
	err := ft.tagMapStore.Put(ft.tagMap)
	if err != nil {
		return err
	}
	ft.tagMap = nil

	return nil
}

func (ft *FTag) Add(file string, tags ...string) error {

	if _, err := os.Stat(file); err != nil {
		return err
	}

	for _, tag := range tags {
		ft.tagMap.Add(file, tag)
	}

	return nil
}

func (ft *FTag) Clear(files ...string) {
	for _, file := range files {
		ft.tagMap.Clear(file)
	}
}

func (ft *FTag) Find(tags ...string) []string {

	// FilesFor files having any of the given tags
	all_files := ft.tagMap.FilesFor(tags...)

	result := make([]string, 0, len(all_files))

	// Retain files that have all the given tags
file_loop:
	for _, file := range all_files {

		for _, tag := range tags {
			if !ft.tagMap.FileToTag.HasValue(file, tag) {
				continue file_loop
			}
		}

		result = append(result, file)
	}

	return result
}

func (ft *FTag) Remove(file string, tags ...string) {
	for _, tag := range tags {
		ft.tagMap.Remove(file, tag)
	}
}

func (ft *FTag) List(files []string) []string {
	if len(files) == 0 {
		files = ft.tagMap.ListFiles()
	}

	tagSet := make(map[string]bool)

	for _, file := range files {
		tags, ok := ft.tagMap.FileToTag[file]
		if !ok {
			continue
		}
		for _, tag := range tags {
			tagSet[tag] = true
		}
	}

	count := len(tagSet)
	tagList := make([]string, count)
	i := 0
	for tag := range tagSet {
		tagList[i] = tag
		i++
	}

	sort.Strings(tagList)
	return tagList
}

func (ft *FTag) Check() []error {
	files := ft.tagMap.ListFiles()

	result := make([]error, 0)

	for _, file := range files {
		if _, err := os.Stat(file); err != nil {
			result = append(result, err)
		}
	}

	return result
}

func (ft *FTag) Move(from, to string) error {

	tags, ok := ft.tagMap.FileToTag[from]
	if !ok {
		return fmt.Errorf("tag mapping for file not found: %s", from)
	}

	// Re-map FileToTag
	delete(ft.tagMap.FileToTag, from)
	ft.tagMap.FileToTag.AddUnique(to, tags...)

	// Re-map TagToFile
	for _, tag := range tags {
		ft.tagMap.TagToFile.RemoveFirst(tag, from)
		ft.tagMap.TagToFile.AddUnique(tag, to)
	}

	return nil
}
