package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("the cache-init image", func() {
	It("can do addition", func() {
		Expect(1+1).To(Equal(2))
	})
})
