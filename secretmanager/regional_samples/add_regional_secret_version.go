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

// [START secretmanager_add_regional_secret_version]
import (
	"context"
	"fmt"
	"hash/crc32"
	"io"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
)

// addSecretVersion adds a new secret version to the given secret with the
// provided payload.
func AddRegionalSecretVersion(w io.Writer, projectId, locationId, secretId string) error {
	// parent := "projects/my-project/locations/my-location/secrets/my-secret"

	// Declare the payload to store.
	payload := []byte("my super secret data")
	// Compute checksum, use Castagnoli polynomial. Providing a checksum
	// is optional.
	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(payload, crc32c))

	// Create the client.
	ctx := context.Background()

	//Endpoint to send the request to regional server
	endpoint := fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", locationId)
	client, err := secretmanager.NewClient(ctx, option.WithEndpoint(endpoint))

	if err != nil {
		return fmt.Errorf("failed to create regional secretmanager client: %w", err)
	}
	defer client.Close()

	parent := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", projectId, locationId, secretId)

	// Build the request.
	req := &secretmanagerpb.AddSecretVersionRequest{
		Parent: parent,
		Payload: &secretmanagerpb.SecretPayload{
			Data:       payload,
			DataCrc32C: &checksum,
		},
	}

	// Call the API.
	result, err := client.AddSecretVersion(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to add regional secret version: %w", err)
	}
	fmt.Fprintf(w, "Added regional secret version: %s\n", result.Name)
	return nil
}

// [END secretmanager_add_regional_secret_version]
