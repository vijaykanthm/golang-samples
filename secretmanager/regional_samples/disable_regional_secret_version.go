// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package regional_secretmanager

// [START secretmanager_disable_regional_secret_version]
import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
)

// disableSecretVersion disables the given secret version. Future requests will
// throw an error until the secret version is enabled. Other secrets versions
// are unaffected.
func DisableRegionalSecretVersion(projectId, locationId, secretId, versionId string) error {
	// name := "projects/my-project/locations/my-location/secrets/my-secret/versions/5"

	// Create the client.
	ctx := context.Background()
	//Endpoint to send the request to regional server
	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationId)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))

	if err != nil {
		return fmt.Errorf("failed to create regional secretmanager client: %w", err)
	}
	defer client.Close()

	name := fmt.Sprintf("projects/%s/locations/%s/secrets/%s/versions/%s", projectId, locationId, secretId, versionId)
	// Build the request.
	req := &secretmanagerpb.DisableSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	if _, err := client.DisableSecretVersion(ctx, req); err != nil {
		return fmt.Errorf("failed to disable regional secret version: %w", err)
	}
	return nil
}

// [END secretmanager_disable_regional_secret_version]
