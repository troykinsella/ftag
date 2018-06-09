package tagmap_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/troykinsella/ftag/tagmap"
)

var _ = Describe("StringListMap", func() {

	Describe("Keys", func() {

		It("should be empty with no entries", func() {
			m := make(tagmap.StringListMap)
			keys := m.Keys()
			Expect(keys).ToNot(BeNil())
			Expect(keys).To(Equal([]string{}))
		})

		It("should list entry keys", func() {
			m := make(tagmap.StringListMap)
			m["foo"] = []string{"foo1", "foo2"}
			m["bar"] = []string{"bar1", "bar2"}

			keys := m.Keys()
			Expect(keys).ToNot(BeNil())
			Expect(keys).To(HaveLen(2))
			Expect(keys).To(ContainElement("foo"))
			Expect(keys).To(ContainElement("bar"))
		})
	})

	Describe("GetOrCreate", func() {

		It("creates a new entry", func() {
			m := make(tagmap.StringListMap)
			list := m.GetOrCreate("new")
			Expect(list).ToNot(BeNil())
			Expect(list).To(Equal([]string{}))
		})

		It("gets an existing entry", func() {
			m := make(tagmap.StringListMap)
			m["foo"] = []string{"foo1", "foo2"}
			list := m.GetOrCreate("foo")
			Expect(list).ToNot(BeNil())
			Expect(list).To(Equal([]string{"foo1", "foo2"}))
		})

	})

	Describe("HasValue", func() {

		It("returns false for a non-existent entry", func() {
			m := make(tagmap.StringListMap)
			hv := m.HasValue("foo", "foo1")
			Expect(hv).To(BeFalse())
		})

		It("returns false for a non-existent list value", func() {
			m := make(tagmap.StringListMap)
			m["foo"] = []string{"foo1", "foo2"}
			hv := m.HasValue("foo", "nofoo")
			Expect(hv).To(BeFalse())
		})

		It("returns true for an existing key and value", func() {
			m := make(tagmap.StringListMap)
			m["foo"] = []string{"foo1", "foo2"}
			hv := m.HasValue("foo", "foo1")
			Expect(hv).To(BeTrue())
		})

	})

	Describe("Add", func() {

		It("should add a new entry", func() {
			m := make(tagmap.StringListMap)
			m.Add("foo", "bar")
			Expect(m).To(Equal(tagmap.StringListMap{
				"foo": []string{"bar"},
			}))
		})

		It("should add to an existing entry", func() {
			m := make(tagmap.StringListMap)
			m["foo"] = []string{"bar1"}
			m.Add("foo", "bar2")
			list := m["foo"]
			Expect(list).To(HaveLen(2))
			Expect(list).To(ContainElement("bar1"))
			Expect(list).To(ContainElement("bar2"))
		})

	})

	Describe("AddUnique", func() {

		It("should add a new entry", func() {
			m := make(tagmap.StringListMap)
			m.AddUnique("foo", "bar")
			Expect(m).To(Equal(tagmap.StringListMap{
				"foo": []string{"bar"},
			}))
		})

		It("should add to an existing entry", func() {
			m := make(tagmap.StringListMap)
			m["foo"] = []string{"bar1"}
			m.AddUnique("foo", "bar2")
			list := m["foo"]
			Expect(list).To(HaveLen(2))
			Expect(list).To(ContainElement("bar1"))
			Expect(list).To(ContainElement("bar2"))
		})

		It("should not add a duplicate entry", func() {
			m := make(tagmap.StringListMap)
			m["foo"] = []string{"bar"}
			m.AddUnique("foo", "bar")
			Expect(m).To(Equal(tagmap.StringListMap{
				"foo": []string{"bar"},
			}))
		})
	})

	Describe("IndexOf", func() {

		It("should return -1 for a non-existent key", func() {
			m := make(tagmap.StringListMap)
			i := m.IndexOf("foo", "bar")
			Expect(i).To(Equal(-1))
		})

		It("should return -1 for a non-existent list entry", func() {
			m := make(tagmap.StringListMap)
			m["foo"] = []string{"bar"}
			i := m.IndexOf("foo", "baz")
			Expect(i).To(Equal(-1))
		})

		It("should return the index of an existing entry", func() {
			m := make(tagmap.StringListMap)
			m["foo"] = []string{"bar1", "bar2"}
			i := m.IndexOf("foo", "bar1")
			Expect(i).To(Equal(0))
			i = m.IndexOf("foo", "bar2")
			Expect(i).To(Equal(1))
		})
	})

	Describe("Remove", func() {

		It("should succeed for non-existent entry", func() {
			m := make(tagmap.StringListMap)
			err := m.Remove("foo", 0)
			Expect(err).To(BeNil())
		})

		It("should error for index below zero", func() {
			m := make(tagmap.StringListMap)
			m["foo"] = []string{"bar1", "bar2"}
			err := m.Remove("foo", -1)
			Expect(err).To(Equal(errors.New("index out of bounds: -1")))
		})

		It("should error for index above range", func() {
			m := make(tagmap.StringListMap)
			m["foo"] = []string{"bar1", "bar2"}
			err := m.Remove("foo", 2)
			Expect(err).To(Equal(errors.New("index out of bounds: 2")))
		})

		It("should remove an existing entry", func() {
			m := make(tagmap.StringListMap)
			m["foo"] = []string{"bar1", "bar2"}
			err := m.Remove("foo", 1)
			Expect(err).To(BeNil())
			Expect(m).To(Equal(tagmap.StringListMap{
				"foo": []string{"bar1"},
			}))
		})

		It("should delete empty lists", func() {
			m := make(tagmap.StringListMap)
			m["foo"] = []string{"bar"}
			err := m.Remove("foo", 0)
			Expect(err).To(BeNil())
			Expect(m["foo"]).To(BeNil())
		})
	})

	Describe("RemoveFirst", func() {

		It("should return false for non-existent key", func() {
			m := make(tagmap.StringListMap)
			r := m.RemoveFirst("foo", "bar")
			Expect(r).To(BeFalse())
		})

		It("should return false for non-existent list value", func() {
			m := make(tagmap.StringListMap)
			m["foo"] = []string{"bar"}
			r := m.RemoveFirst("foo", "baz")
			Expect(r).To(BeFalse())
			Expect(m).To(Equal(tagmap.StringListMap{
				"foo": []string{"bar"},
			}))
		})

		It("should remove an existing list value", func() {
			m := make(tagmap.StringListMap)
			m["foo"] = []string{"bar", "baz"}
			r := m.RemoveFirst("foo", "bar")
			Expect(r).To(BeTrue())
			Expect(m).To(Equal(tagmap.StringListMap{
				"foo": []string{"baz"},
			}))
		})

		It("should delete empty lists", func() {
			m := make(tagmap.StringListMap)
			m["foo"] = []string{"bar"}
			r := m.RemoveFirst("foo", "bar")
			Expect(r).To(BeTrue())
			Expect(m["foo"]).To(BeNil())
		})

	})

})
