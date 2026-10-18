//go:build conformance

// The conformance tests run against the live Loopia API and need envtest
// binaries, so they are behind the "conformance" build tag. Run them with `make test`.

package main

import (
	"os"
	"strconv"
	"testing"
	"time"

	acmetest "github.com/cert-manager/cert-manager/test/acme"
)

var (
	// Environment variable holding the name of the zone to test, ex: example.com however this needs to be a zone you have control over in Loopia.
	// This needs to be set before running the test.
	zone = os.Getenv("TEST_ZONE_NAME")

	// Environment variable for enabling the strict mode testing.
	strictmodeenv = os.Getenv("TEST_STRICT_MODE")
)

func TestRunsSuite(t *testing.T) {
	// The conformance suite talks to the live Loopia API, so it only runs when a test zone is given.
	if zone == "" {
		t.Skip("TEST_ZONE_NAME not set, skipping conformance tests against the live Loopia API")
	}

	// The manifest path should contain a file named config.json that is a sniplet of valid configuration that should be included on the ChallengeRequest passed as part of the test cases.
	// The test fixture also requires etcd, kube-apiserver and kubectl, found via the TEST_ASSET_ETCD, TEST_ASSET_KUBE_APISERVER and TEST_ASSET_KUBECTL environment variables or on PATH; `make test` downloads them with setup-envtest.

	var strictmode, err = strconv.ParseBool(strictmodeenv)
	if err != nil {
		strictmode = false
	}

	solver := &loopiaDNSProviderSolver{}
	fixture := acmetest.NewFixture(solver,
		acmetest.SetStrict(strictmode),
		acmetest.SetResolvedZone(zone),
		acmetest.SetAllowAmbientCredentials(false),
		acmetest.SetManifestPath("testdata/loopia"),
		acmetest.SetPollInterval(time.Second*60),
		acmetest.SetPropagationLimit(time.Minute*30),
	)

	fixture.RunConformance(t)
}
