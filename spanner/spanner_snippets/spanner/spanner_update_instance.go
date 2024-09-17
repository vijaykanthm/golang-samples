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

package spanner

// [START spanner_update_instance]
import (
	"context"
	"fmt"
	"io"

	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/admin/instance/apiv1/instancepb"
	"google.golang.org/genproto/protobuf/field_mask"
)

func updateInstance(w io.Writer, projectID, instanceID string) error {
	// projectID := "my-project-id"
	// instanceID := "my-instance"
	ctx := context.Background()
	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return err
	}
	defer instanceAdmin.Close()

	req := &instancepb.UpdateInstanceRequest{
		Instance: &instancepb.Instance{
			Name: fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID),
			// The edition selected for this instance.
			// Different editions provide different capabilities at different price points.
			// For more information, see https://cloud.google.com/spanner/docs/editions-overview.
			Edition: instancepb.Instance_ENTERPRISE,
		},
		FieldMask: &field_mask.FieldMask{
			Paths: []string{"edition"},
		},
	}
	op, err := instanceAdmin.UpdateInstance(ctx, req)
	if err != nil {
		return fmt.Errorf("could not update instance %s: %w", fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID), err)
	}
	// Wait for the instance update to finish.
	_, err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for instance update to finish failed: %w", err)
	}

	fmt.Fprintf(w, "Updated instance [%s]\n", instanceID)
	return nil
}

// [END spanner_update_instance]
