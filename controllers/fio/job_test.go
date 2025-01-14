/*
Copyright 2019 The xridge kubestone contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package fio

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"

	perfv1alpha1 "github.com/xridge/kubestone/api/v1alpha1"
)

var _ = Describe("fio job", func() {
	Describe("cr with cmd args", func() {
		var cr perfv1alpha1.Fio
		var job *batchv1.Job

		BeforeEach(func() {
			cr = perfv1alpha1.Fio{
				Spec: perfv1alpha1.FioSpec{
					Image: perfv1alpha1.ImageSpec{
						Name:       "xridge/fio:test",
						PullPolicy: "Always",
						PullSecret: "a-pull-secret",
					},
					CmdLineArgs: "--name=randwrite --iodepth=1 --rw=randwrite --bs=4m --direct=1 --size=256M --numjobs=1",
				},
			}
			job = NewJob(&cr)
		})

		Context("with Image details specified", func() {
			It("should match on Image.Name", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Image).To(
					Equal(cr.Spec.Image.Name))
			})
			It("should match on Image.PullPolicy", func() {
				Expect(job.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(
					Equal(corev1.PullPolicy(cr.Spec.Image.PullPolicy)))
			})
			It("should match on Image.PullSecret", func() {
				Expect(job.Spec.Template.Spec.ImagePullSecrets[0].Name).To(
					Equal(cr.Spec.Image.PullSecret))
			})
		})

		Context("with command line args specified", func() {
			It("should have the same args", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("--name=randwrite"))
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("--iodepth=1"))
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("--size=256M"))
			})
		})
	})

	Describe("cr with builtin job files", func() {
		var cr perfv1alpha1.Fio
		var job *batchv1.Job

		BeforeEach(func() {
			cr = perfv1alpha1.Fio{
				Spec: perfv1alpha1.FioSpec{
					Image: perfv1alpha1.ImageSpec{
						Name:       "xridge/fio:test",
						PullPolicy: "IfNotPresent",
					},
					BuiltinJobFiles: []string{"/jobs/rand-read.fio"},
				},
			}
			job = NewJob(&cr)
		})

		Context("with Image details specified", func() {
			It("should match on Image.Name", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Image).To(
					Equal(cr.Spec.Image.Name))
			})
			It("should match on Image.PullPolicy", func() {
				Expect(job.Spec.Template.Spec.Containers[0].ImagePullPolicy).To(
					Equal(corev1.PullPolicy(cr.Spec.Image.PullPolicy)))
			})
		})

		Context("with builtin job files specified", func() {
			It("should have those files", func() {
				Expect(job.Spec.Template.Spec.Containers[0].Args).To(
					ContainElement("/jobs/rand-read.fio"))
			})
		})
	})
})
