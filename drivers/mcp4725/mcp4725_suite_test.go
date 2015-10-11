package mcp4725

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMcp4725(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mcp4725 Suite")
}

var _ = Describe("mcp4725 conversion", func() {

	It("converts 0xfff to the right value", func() {
		b := toBytes(0xfff)
		Expect(b).To(Equal([]byte{0xff, 0xf0}))
	})

	It("converts 0x4fff to the right value", func() {
		b := toBytes(0x4fff)
		Expect(b).To(Equal([]byte{0xff, 0xf0}))
	})

	It("converts 0x000 to the right value", func() {
		b := toBytes(0x000)
		Expect(b).To(Equal([]byte{0x00, 0x00}))
	})

	It("converts 0x1000 to the right value", func() {
		b := toBytes(0x1000)
		Expect(b).To(Equal([]byte{0x00, 0x00}))
	})

	It("converts 0x800 to the right value", func() {
		b := toBytes(0x800)
		Expect(b).To(Equal([]byte{0x80, 0x00}))
	})

})
