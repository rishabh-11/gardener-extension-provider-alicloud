// Copyright (c) 2018 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helper_test

import (
	api "github.com/gardener/gardener-extension-provider-alicloud/pkg/apis/alicloud"
	. "github.com/gardener/gardener-extension-provider-alicloud/pkg/apis/alicloud/helper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/utils/pointer"
)

const profileImageID = "id-1235"

var _ = Describe("Helper", func() {
	var (
		purpose      api.Purpose = "foo"
		purposeWrong api.Purpose = "baz"
	)
	DescribeTable("#FindVSwitchForPurposeAndZone",
		func(vswitches []api.VSwitch, purpose api.Purpose, zone string, expectedVSwitch *api.VSwitch, expectErr bool) {
			subnet, err := FindVSwitchForPurposeAndZone(vswitches, purpose, zone)
			expectResults(subnet, expectedVSwitch, err, expectErr)
		},

		Entry("list is nil", nil, purpose, "europe", nil, true),
		Entry("empty list", []api.VSwitch{}, purpose, "europe", nil, true),
		Entry("entry not found (no purpose)", []api.VSwitch{{ID: "bar", Purpose: purposeWrong, Zone: "europe"}}, purpose, "europe", nil, true),
		Entry("entry not found (no zone)", []api.VSwitch{{ID: "bar", Purpose: purposeWrong, Zone: "europe"}}, purpose, "asia", nil, true),
		Entry("entry exists", []api.VSwitch{{ID: "bar", Purpose: purposeWrong, Zone: "europe"}}, purposeWrong, "europe", &api.VSwitch{ID: "bar", Purpose: purposeWrong, Zone: "europe"}, false),
	)

	DescribeTable("#FindVSwitchForPurpose",
		func(vswitches []api.VSwitch, purpose api.Purpose, expectedVSwitch *api.VSwitch, expectErr bool) {
			subnet, err := FindVSwitchForPurpose(vswitches, purpose)
			expectResults(subnet, expectedVSwitch, err, expectErr)
		},

		Entry("list is nil", nil, purpose, nil, true),
		Entry("empty list", []api.VSwitch{}, purpose, nil, true),
		Entry("entry not found (no purpose)", []api.VSwitch{{ID: "bar", Purpose: purposeWrong, Zone: "europe"}}, purpose, nil, true),
		Entry("entry exists", []api.VSwitch{{ID: "bar", Purpose: purposeWrong, Zone: "europe"}}, purposeWrong, &api.VSwitch{ID: "bar", Purpose: purposeWrong, Zone: "europe"}, false),
	)

	DescribeTable("#FindSecurityGroupByPurpose",
		func(securityGroups []api.SecurityGroup, purpose api.Purpose, expectedSecurityGroup *api.SecurityGroup, expectErr bool) {
			securityGroup, err := FindSecurityGroupByPurpose(securityGroups, purpose)
			expectResults(securityGroup, expectedSecurityGroup, err, expectErr)
		},

		Entry("list is nil", nil, purpose, nil, true),
		Entry("empty list", []api.SecurityGroup{}, purpose, nil, true),
		Entry("entry not found", []api.SecurityGroup{{ID: "bar", Purpose: purposeWrong}}, purpose, nil, true),
		Entry("entry exists", []api.SecurityGroup{{ID: "bar", Purpose: purpose}}, purpose, &api.SecurityGroup{ID: "bar", Purpose: purpose}, false),
	)

	DescribeTable("#FindMachineImage",
		func(machineImage []api.MachineImage, name, version string, encrypted bool, expectedMachineImage *api.MachineImage, expectErr bool) {
			found, err := FindMachineImage(machineImage, name, version, encrypted)
			expectResults(found, expectedMachineImage, err, expectErr)
		},

		Entry("list is nil", nil, "foo", "1.2.3", true, nil, true),
		Entry("empty list", []api.MachineImage{}, "foo", "1.2.3", true, nil, true),
		Entry("entry not found (no name)", []api.MachineImage{{Name: "bar", Version: "1.2.3", ID: "id123"}}, "foo", "1.2.3", true, nil, true),
		Entry("entry not found (no version)", []api.MachineImage{{Name: "bar", Version: "1.2.3", ID: "id123"}}, "foo", "1.2.4", true, nil, true),
		Entry("entry not found (empty encrypted)", []api.MachineImage{{Name: "bar", Version: "1.2.3", ID: "id123"}}, "bar", "1.2.3", true, nil, true),
		Entry("entry not found (false encrypted)", []api.MachineImage{{Name: "bar", Version: "1.2.3", ID: "id123", Encrypted: pointer.BoolPtr(false)}}, "bar", "1.2.3", true, nil, true),

		Entry("entry exists (encrypted value exists)", []api.MachineImage{{Name: "bar", Version: "1.2.3", ID: "id123", Encrypted: pointer.BoolPtr(true)}}, "bar", "1.2.3", true, &api.MachineImage{Name: "bar", Version: "1.2.3", ID: "id123", Encrypted: pointer.BoolPtr(true)}, false),
		Entry("entry exists (empty encrypted value)", []api.MachineImage{{Name: "bar", Version: "1.2.3", ID: "id123"}}, "bar", "1.2.3", false, &api.MachineImage{Name: "bar", Version: "1.2.3", ID: "id123"}, false),
	)

	Describe("#AppendMachineImage",
		func() {

			It("should append a non-existing image", func() {
				existingImages := []api.MachineImage{{Name: "bar", Version: "1.2.3", ID: "id123", Encrypted: pointer.BoolPtr(true)}}
				imageToInsert := api.MachineImage{Name: "bar", Version: "1.2.4", ID: "id123"}
				existingImages = AppendMachineImage(existingImages, imageToInsert)
				Expect(len(existingImages)).To(Equal(2))
				Expect(existingImages, ContainElement(imageToInsert))
			})

			It("should not append the image", func() {
				imageToInsert := api.MachineImage{Name: "bar", Version: "1.2.3", ID: "id123", Encrypted: pointer.BoolPtr(false)}
				imageExisting := api.MachineImage{Name: "bar", Version: "1.2.3", ID: "id123"}
				existingImages := []api.MachineImage{imageExisting}
				existingImages = AppendMachineImage(existingImages, imageToInsert)
				Expect(len(existingImages)).To(Equal(1))
				Expect(existingImages[0]).To(Equal(imageExisting))
			})
		})

	DescribeTable("#FindImageForRegion",
		func(profileImages []api.MachineImages, imageName, version, region string, expectedImage string) {
			cfg := &api.CloudProfileConfig{}
			cfg.MachineImages = profileImages
			image, err := FindImageForRegionFromCloudProfile(cfg, imageName, version, region)

			Expect(image).To(Equal(expectedImage))
			if expectedImage != "" {
				Expect(err).NotTo(HaveOccurred())
			} else {
				Expect(err).To(HaveOccurred())
			}
		},

		Entry("list is nil", nil, "ubuntu", "1", "china", ""),

		Entry("profile empty list", []api.MachineImages{}, "ubuntu", "1", "china", ""),
		Entry("profile entry not found (image does not exist)", makeProfileMachineImages("debian", "1", "china"), "ubuntu", "1", "china", ""),
		Entry("profile entry not found (version does not exist)", makeProfileMachineImages("ubuntu", "2", "china"), "ubuntu", "1", "china", ""),
		Entry("profile entry", makeProfileMachineImages("ubuntu", "1", "china"), "ubuntu", "1", "china", profileImageID),
		Entry("profile non matching region", makeProfileMachineImages("ubuntu", "1", "china"), "ubuntu", "1", "eu", ""),
	)
})

func makeProfileMachineImages(name, version, region string) []api.MachineImages {
	versions := []api.MachineImageVersion{
		{
			Version: version,
			Regions: []api.RegionIDMapping{
				{
					Name: region,
					ID:   profileImageID,
				},
			},
		},
	}

	return []api.MachineImages{
		{
			Name:     name,
			Versions: versions,
		},
	}
}

func expectResults(result, expected interface{}, err error, expectErr bool) {
	if !expectErr {
		Expect(result).To(Equal(expected))
		Expect(err).NotTo(HaveOccurred())
	} else {
		Expect(result).To(BeNil())
		Expect(err).To(HaveOccurred())
	}
}
