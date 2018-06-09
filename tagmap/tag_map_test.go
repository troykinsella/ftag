package tagmap_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/troykinsella/ftag/tagmap"
	. "github.com/onsi/gomega"
)

var _ = Describe("TagMap", func() {

	Describe("New", func() {

		It("should not return nil", func() {
			tm := tagmap.New()
			Expect(tm).ToNot(BeNil())
		})

		It("should populate fields", func() {
			tm := tagmap.New()
			Expect(tm.TagToFile).ToNot(BeNil())
			Expect(tm.FileToTag).ToNot(BeNil())
		})

	})

	Describe("ListFiles", func() {

		It("should return an empty list when no entries", func() {
			tm := tagmap.New()
			l := tm.ListFiles()
			Expect(l).To(Equal([]string{}))
		})

		It("should return a list of existing file names", func() {
			tm := tagmap.New()
			tm.Add("foo", "tag1")
			tm.Add("bar", "tag2")
			l := tm.ListFiles()
			Expect(l).ToNot(BeNil())
			Expect(l).To(HaveLen(2))
			Expect(l).To(ContainElement("foo"))
			Expect(l).To(ContainElement("bar"))
		})

	})

	Describe("FilesFor", func() {

		It("should return an empty slice for no matches", func() {
			tm := tagmap.New()
			f := tm.FilesFor("foo")
			Expect(f).ToNot(BeNil())
			Expect(f).To(Equal([]string{}))
		})

		It("should return files having any given tag", func() {
			tm := tagmap.New()
			tm.Add("foo", "tag1")
			tm.Add("bar", "tag2")
			tm.Add("baz", "tag1")
			tm.Add("biz", "tag3")

			f := tm.FilesFor("tag1", "tag2")
			Expect(f).ToNot(BeNil())
			Expect(f).To(HaveLen(3))
			Expect(f).To(ContainElement("foo"))
			Expect(f).To(ContainElement("bar"))
			Expect(f).To(ContainElement("baz"))
		})

	})

	Describe("Add", func() {

		It("should create bi-directional mappings", func() {
			tm := tagmap.New()
			tm.Add("foo", "tag1")
			tm.Add("bar", "tag2")
			tm.Add("baz", "tag1")

			Expect(tm.FileToTag).To(HaveLen(3))
			Expect(tm.FileToTag["foo"]).To(Equal([]string{"tag1"}))
			Expect(tm.FileToTag["bar"]).To(Equal([]string{"tag2"}))
			Expect(tm.FileToTag["baz"]).To(Equal([]string{"tag1"}))

			Expect(tm.TagToFile).To(HaveLen(2))
			Expect(tm.TagToFile["tag1"]).To(HaveLen(2))
			Expect(tm.TagToFile["tag1"]).To(ContainElement("foo"))
			Expect(tm.TagToFile["tag1"]).To(ContainElement("baz"))
			Expect(tm.TagToFile["tag1"]).ToNot(ContainElement("bar"))

			Expect(tm.TagToFile["tag2"]).To(Equal([]string{"bar"}))
		})

	})

	Describe("Remove", func() {

		It("should return false for non-existent entries", func() {
			tm := tagmap.New()
			r := tm.Remove("foo", "tag")
			Expect(r).To(BeFalse())
		})

		It("should remove bi-directional mappings", func() {
			tm := tagmap.New()
			tm.Add("foo", "tag1")
			tm.Add("bar", "tag2")
			tm.Add("baz", "tag1")

			r := tm.Remove("foo", "tag1")
			Expect(r).To(BeTrue())

			Expect(tm.FileToTag).To(HaveLen(2))
			Expect(tm.FileToTag["foo"]).To(BeNil())
			Expect(tm.FileToTag["bar"]).To(Equal([]string{"tag2"}))
			Expect(tm.FileToTag["baz"]).To(Equal([]string{"tag1"}))

			Expect(tm.TagToFile).To(HaveLen(2))
			Expect(tm.TagToFile["tag1"]).ToNot(ContainElement("foo"))
			Expect(tm.TagToFile["tag1"]).To(ContainElement("baz"))
			Expect(tm.TagToFile["tag1"]).ToNot(ContainElement("bar"))

			Expect(tm.TagToFile["tag2"]).To(Equal([]string{"bar"}))
		})

	})

	Describe("Clear", func() {

		It("should not panic for non-existent file", func() {
			tm := tagmap.New()
			tm.Clear("foo")
		})

		It("should remove all tags for file", func() {
			tm := tagmap.New()
			tm.Add("foo", "tag1")
			tm.Add("bar", "tag1")
			tm.Clear("foo")
			Expect(tm.FileToTag).ToNot(ContainElement("foo"))
			Expect(tm.TagToFile["tag1"]).To(Equal([]string{"bar"}))
		})

	})

	Describe("Normalize", func() {

		It("should set the version", func() {
			tm := tagmap.New()
			tm.Version = "some shit"
			tm.Normalize()
			Expect(tm.Version).To(Equal(tagmap.TM_VERSION))
		})

		It("should sort tag lists", func() {
			tm := tagmap.New()
			tm.FileToTag["foo"] = []string{"b", "a", "c"}
			tm.Normalize()
			Expect(tm.FileToTag["foo"]).To(Equal([]string{"a", "b", "c"}))
		})
	})

})
