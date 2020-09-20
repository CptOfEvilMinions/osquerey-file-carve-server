package auth

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/CptOfEvilMinions/osquery-file-carve-server/pkg/config"
)

type vaultTokenLookupResult struct {
	Error string `json:"error"`
	Data  struct {
		Token         string    `json:"id"`
		ExpireTime    time.Time `json:"expire_time"`
		TokenAccessor string    `json:"accessor"`
		Policies      []string  `json:"identity_policies"`
	}
}

var vaultURL string
var vaultVerifyTLS bool

// InitVault define global variables to use when interacting with Vault
func InitVault(cfg *config.Config) {
	// Create token validatio URL
	// Enforcing HTTPS
	vaultURL = fmt.Sprintf("https://%s:%d/v1/auth/token/lookup-self", cfg.Vault.Hostname, cfg.Vault.Port)
	vaultVerifyTLS = cfg.Vault.VerifyTLS

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: vaultVerifyTLS},
	}
	client := &http.Client{Transport: tr}
	// Create HTTP request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s:%d/v1/sys/health", cfg.Vault.Hostname, cfg.Vault.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

}

func vaultTokenPolicyLookup(vaultPolicyCheck string, policies []string) bool {
	for _, i := range policies {
		if i == vaultPolicyCheck {
			return true
		}
	}
	return false
}

func vaultTokenLookup(token string, vaultPolicyCheck string) error {
	// Init HTTP client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: vaultVerifyTLS},
	}
	client := &http.Client{Transport: tr}

	// Create HTTP request
	req, err := http.NewRequest("GET", vaultURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Vault-Token", token)

	// Do HTTP request
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Unable to retrieve authentication token from vault %q", err)
	} else if res.StatusCode != 200 {
		return fmt.Errorf("Unable to retrieve authentication token from vault (status code %d)", res.StatusCode)
	}

	// Extrat body from payload
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Unable to parse request body: %q", err)
	}

	// Does body contian token
	var vtResuslt vaultTokenLookupResult
	err = json.Unmarshal(body, &vtResuslt)
	if err != nil {
		return fmt.Errorf("Unable to parse request body: %q", err)
	}

	// Make sure there are no auth errors
	if vtResuslt.Error != "" || vtResuslt.Data.Token != token {
		return fmt.Errorf("Could not validate token: %s", vtResuslt.Error)
	}

	// Make sure the token has a specific policy
	if result := vaultTokenPolicyLookup(vaultPolicyCheck, vtResuslt.Data.Policies); result != true {
		return fmt.Errorf("Token does not have the appropriate identity policy")
	}

	// Add valid token and expiration date to map
	mutex.Lock()
	tokenSessionMap[vtResuslt.Data.TokenAccessor] = vtResuslt.Data.ExpireTime
	mutex.Unlock()

	return nil
}
