package google

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"golang.org/x/oauth2/google"
)

const testFakeCredentialsPath = "./test-fixtures/fake_account.json"
const testOauthScope = "https://www.googleapis.com/auth/compute"

func TestConfigLoadAndValidate_accountFilePath(t *testing.T) {
	config := Config{
		Credentials: testFakeCredentialsPath,
		Project:     "my-gce-project",
		Region:      "us-central1",
	}

	err := config.loadAndValidate()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestConfigLoadAndValidate_accountFileJSON(t *testing.T) {
	contents, err := ioutil.ReadFile(testFakeCredentialsPath)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	config := Config{
		Credentials: string(contents),
		Project:     "my-gce-project",
		Region:      "us-central1",
	}

	err = config.loadAndValidate()
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestConfigLoadAndValidate_accountFileJSONInvalid(t *testing.T) {
	config := Config{
		Credentials: "{this is not json}",
		Project:     "my-gce-project",
		Region:      "us-central1",
	}

	if config.loadAndValidate() == nil {
		t.Fatalf("expected error, but got nil")
	}
}

func TestAccConfigLoadValidate_credentials(t *testing.T) {
	if os.Getenv(resource.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Network access not allowed; use %s=1 to enable", resource.TestEnvVar))
	}

	creds := getTestCredsFromEnv()
	proj := getTestProjectFromEnv()

	config := Config{
		Credentials: creds,
		Project:     proj,
		Region:      "us-central1",
	}

	err := config.loadAndValidate()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	_, err = config.clientCompute.Zones.Get(proj, "us-central1-a").Do()
	if err != nil {
		t.Fatalf("expected call with loaded config client to work, got error: %s", err)
	}
}

func TestAccConfigLoadValidate_accessToken(t *testing.T) {
	if os.Getenv(resource.TestEnvVar) == "" {
		t.Skip(fmt.Sprintf("Network access not allowed; use %s=1 to enable", resource.TestEnvVar))
	}

	creds := getTestCredsFromEnv()
	proj := getTestProjectFromEnv()

	c, err := google.CredentialsFromJSON(context.Background(), []byte(creds), testOauthScope)
	if err != nil {
		t.Fatalf("invalid test credentials: %s", err)
	}

	token, err := c.TokenSource.Token()
	if err != nil {
		t.Fatalf("Unable to generate test access token: %s", err)
	}

	config := Config{
		AccessToken: token.AccessToken,
		Project:     proj,
		Region:      "us-central1",
	}

	err = config.loadAndValidate()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	_, err = config.clientCompute.Zones.Get(proj, "us-central1-a").Do()
	if err != nil {
		t.Fatalf("expected API call with loaded config to work, got error: %s", err)
	}
}
