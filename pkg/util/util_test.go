package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing RKE2 to Kubernetes Version conversion", func() {

	machineVersion := "v1.24.6"
	rke2Version := "v1.24.6+rke2r1"
	It("Should match RKE2 and Kubernetes version", func() {
		cpKubeVersion, err := Rke2ToKubeVersion(rke2Version)
		Expect(err).ToNot(HaveOccurred())
		Expect(cpKubeVersion).To(Equal(machineVersion))
	})

})
