package proofing_test

import (
	"github.com/pivotal-cf/kiln/internal/proofing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CollectionPropertyInput", func() {
	var collectionPropertyInput proofing.CollectionPropertyInput

	BeforeEach(func() {
		productTemplate, err := proofing.Parse("fixtures/form_types.yml")
		Expect(err).NotTo(HaveOccurred())

		var ok bool
		collectionPropertyInput, ok = productTemplate.FormTypes[0].PropertyInputs[1].(proofing.CollectionPropertyInput)
		Expect(ok).To(BeTrue())
	})

	It("parses their structure", func() {
		Expect(collectionPropertyInput.Description).To(Equal("some-description"))
		Expect(collectionPropertyInput.Label).To(Equal("some-label"))
		Expect(collectionPropertyInput.Placeholder).To(Equal("some-placeholder"))
		Expect(collectionPropertyInput.Reference).To(Equal("some-reference"))

		Expect(collectionPropertyInput.PropertyInputs).To(HaveLen(1))
	})

	Context("property_inputs", func() {
		It("parses their structure", func() {
			propertyInput := collectionPropertyInput.PropertyInputs[0]

			Expect(propertyInput.Description).To(Equal("some-description"))
			Expect(propertyInput.Label).To(Equal("some-label"))
			Expect(propertyInput.Placeholder).To(Equal("some-placeholder"))
			Expect(propertyInput.Reference).To(Equal("some-reference"))
			Expect(propertyInput.Slug).To(BeTrue())
		})
	})
})
