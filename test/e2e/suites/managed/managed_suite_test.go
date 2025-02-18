//go:build e2e
// +build e2e

/*
Copyright 2020 The Kubernetes Authors.

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

package managed

import (
	"flag"
	"testing"

	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"

	ekscontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/eks/api/v1beta2"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/shared"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	expclusterv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"
	"sigs.k8s.io/cluster-api/test/framework"
)

var (
	e2eCtx           *shared.E2EContext
	skipUpgradeTests bool
	skipGeneralTests bool
	skipLegacyTests  bool
)

func init() {
	e2eCtx = shared.NewE2EContext(shared.WithManaged(), shared.WithSchemeInit(initScheme))

	shared.CreateDefaultFlags(e2eCtx)
	flag.BoolVar(&skipGeneralTests, "skip-eks-general-tests", false, "if true, the general EKS tests will be skipped")
	flag.BoolVar(&skipUpgradeTests, "skip-eks-upgrade-tests", false, "if true, the EKS upgrade tests will be skipped")
	flag.BoolVar(&skipLegacyTests, "skip-eks-legacy-tests", false, "if true, the EKS legacy tests will be skipped")
}

func TestE2E(t *testing.T) {
	RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecsWithDefaultAndCustomReporters(t, "capa-eks-e2e", []ginkgo.Reporter{framework.CreateJUnitReporterForProw(e2eCtx.Settings.ArtifactFolder)})
}

var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {
	return shared.Node1BeforeSuite(e2eCtx)
}, func(data []byte) {
	shared.AllNodesBeforeSuite(e2eCtx, data)
})

var _ = ginkgo.SynchronizedAfterSuite(
	func() {
		shared.AllNodesAfterSuite(e2eCtx)
	},
	func() {
		shared.Node1AfterSuite(e2eCtx)
	},
)

func runGeneralTests() bool {
	return !skipGeneralTests
}

func runUpgradeTests() bool {
	return !skipUpgradeTests
}

func runLegacyTests() bool {
	return !skipLegacyTests
}

func initScheme() *runtime.Scheme {
	sc := shared.DefaultScheme()
	_ = expinfrav1.AddToScheme(sc)
	_ = clusterv1.AddToScheme(sc)
	_ = ekscontrolplanev1.AddToScheme(sc)
	_ = expclusterv1.AddToScheme(sc)

	return sc
}
