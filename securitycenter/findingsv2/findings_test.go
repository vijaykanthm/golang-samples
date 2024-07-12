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

package findingsv2

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	securitycenter "cloud.google.com/go/securitycenter/apiv2"
	"cloud.google.com/go/securitycenter/apiv2/securitycenterpb"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
	"github.com/golang/protobuf/ptypes"
)

var orgID = ""
var sourceName = ""
var findingName = ""
var untouchedFindingName = ""

func createTestFinding(ctx context.Context, client *securitycenter.Client, findingID string, category string) (*securitycenterpb.Finding, error) {
	eventTime, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		return nil, fmt.Errorf("TimestampProto: %w", err)
	}

	req := &securitycenterpb.CreateFindingRequest{
		Parent:    sourceName,
		FindingId: findingID,
		Finding: &securitycenterpb.Finding{
			State: securitycenterpb.Finding_ACTIVE,
			// Resource the finding is associated with.  This is an
			// example any resource identifier can be used.
			ResourceName: "//cloudresourcemanager.googleapis.com/organizations/11232",
			// A free-form category.
			Category: category,
			// The time associated with discovering the issue.
			EventTime: eventTime,
		},
	}
	return client.CreateFinding(ctx, req)
}

// setupEntities initializes variables in this file with entityNames to
// use for testing.
func setupEntities() error {
	orgID = os.Getenv("GCLOUD_ORGANIZATION")
	if orgID == "" {
		// Each test checks for GCLOUD_ORGANIZATION. Return nil so we see every skip.
		return nil
	}
	ctx := context.Background()
	client, err := securitycenter.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("securitycenter.NewClient: %w", err)
	}
	defer client.Close() // Closing the client safely cleans up background resources.

	source, err := client.CreateSource(ctx, &securitycenterpb.CreateSourceRequest{
		Source: &securitycenterpb.Source{
			DisplayName: "Customized Display Name",
			Description: "A new custom source that does X",
		},
		Parent: fmt.Sprintf("organizations/%s", orgID),
	})

	if err != nil {
		return fmt.Errorf("CreateSource: %w", err)
	}
	sourceName = source.Name
	finding, err := createTestFinding(ctx, client, "updated", "MEDIUM_RISK_ONE")
	if err != nil {
		return fmt.Errorf("createTestFinding: %w", err)
	}
	findingName = finding.Name
	finding, err = createTestFinding(ctx, client, "untouched", "XSS")
	if err != nil {
		return fmt.Errorf("createTestFinding: %w", err)
	}
	untouchedFindingName = finding.Name
	return nil
}

func setup(t *testing.T) string {
	if orgID == "" {
		t.Skip("GCLOUD_ORGANIZATION not set")
	}
	if sourceName == "" {
		t.Errorf("sourceName not set")
		os.Exit(1)
	}
	return orgID
}

func TestMain(m *testing.M) {
	if err := setupEntities(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialize findings test environment: %v\n", err)
		return
	}
	code := m.Run()
	os.Exit(code)
}

func TestCreateSource(t *testing.T) {
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		orgID := setup(t)
		buf := new(bytes.Buffer)

		err := createSource(buf, orgID)

		if err != nil {
			r.Errorf("createSource(%s) had error: %v", orgID, err)
			return
		}

		got := buf.String()
		if want := "New source created"; !strings.Contains(got, want) {
			r.Errorf("createSource(%s) got: %s want %s", orgID, got, want)
		}
		if !strings.Contains(got, orgID) {
			r.Errorf("createSource(%s) got: %s want %s", orgID, got, orgID)
		}
	})
}

func TestGetSource(t *testing.T) {
	setup(t)
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		err := getSource(buf, sourceName)

		if err != nil {
			r.Errorf("getSource(%s) had error: %v", sourceName, err)
			return
		}

		got := buf.String()
		if want := "Description: A new custom source that does X"; !strings.Contains(got, want) {
			r.Errorf("getSource(%s) got: %s want %s", sourceName, got, want)
		}
		if !strings.Contains(got, orgID) {
			r.Errorf("getSource(%s) got: %s want %s", sourceName, got, sourceName)
		}
	})
}

func TestListSources(t *testing.T) {
	setup(t)
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		err := listSources(buf, orgID)

		if err != nil {
			r.Errorf("listSource(%s) had error: %v", orgID, err)
			return
		}

		got := buf.String()
		if want := "Description: A new custom source that does X"; !strings.Contains(got, want) {
			r.Errorf("listSource(%s) got: %s want %s", orgID, got, want)
		}
		if !strings.Contains(got, sourceName) {
			r.Errorf("listSource(%s) got: %s want %s", orgID, got, sourceName)
		}
	})
}

func TestUpdateSource(t *testing.T) {
	setup(t)
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		err := updateSource(buf, sourceName)

		if err != nil {
			r.Errorf("updateSource(%s) had error: %v", sourceName, err)
			return
		}

		got := buf.String()
		if want := "Display name: New Display Name"; !strings.Contains(got, want) {
			r.Errorf("updateSource(%s) got: %s want %s", sourceName, got, want)
		}
		if !strings.Contains(got, sourceName) {
			r.Errorf("updateSource(%s) got: %s want %s", sourceName, got, sourceName)
		}
		if want := "does X"; !strings.Contains(got, want) {
			r.Errorf("updateSource(%s) got: %s want %s", sourceName, got, want)
		}

	})
}

func TestListAllFindings(t *testing.T) {
	// issue #3260 tracks this test skip.ß
	t.Skip()
	testutil.Retry(t, 5, 20*time.Second, func(r *testutil.R) {
		orgID := setup(t)
		buf := new(bytes.Buffer)

		err := listFindings(buf, orgID)

		if err != nil {
			r.Errorf("listFindings(%s) had error: %v", orgID, err)
			return
		}

		got := buf.String()
		if !strings.Contains(got, findingName) {
			r.Errorf("listFindings(%s) got: %s want %s", orgID, got, findingName)
		}

		if !strings.Contains(got, untouchedFindingName) {
			r.Errorf("listFindings(%s) got: %s want %s", orgID, got, untouchedFindingName)
		}
	})
}

func TestListFilteredFindings(t *testing.T) {
	setup(t)
	testutil.Retry(t, 5, 20*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		err := listFilteredFindings(buf, sourceName)

		if err != nil {
			r.Errorf("listFilteredFindings(%s) had error: %v", sourceName, err)
			return
		}

		got := buf.String()
		if !strings.Contains(got, findingName) {
			r.Errorf("listFilteredFindings(%s) got: %s want %s", sourceName, got, findingName)
		}

		if strings.Contains(got, untouchedFindingName) {
			r.Errorf("listFilteredFindings(%s) got: %s didn't want %s", sourceName, got, untouchedFindingName)
		}
	})
}

func TestAddSecurityMarks(t *testing.T) {
	setup(t)
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		err := addSecurityMarks(buf, findingName)

		if err != nil {
			r.Errorf("addSecurityMarks(%s) adding marks had error: %v", findingName, err)
			return
		}

		got := buf.String()
		if want := "key_a = value_a"; !strings.Contains(got, want) {
			r.Errorf("addSecurityMarks(%s) got: %s want %s", findingName, got, want)
		}

		if want := "key_b = value_b"; !strings.Contains(got, want) {
			r.Errorf("addSecurityMarks(%s) got: %s want %s", findingName, got, want)
		}
	})
}

func TestListFindingsWithMarks(t *testing.T) {
	testutil.Retry(t, 5, 20*time.Second, func(r *testutil.R) {
		orgID := setup(t)
		buf := new(bytes.Buffer)
		// Ensure security marks have been added so filter is effective.
		err := addSecurityMarks(buf, findingName)
		buf.Truncate(0)
		if err != nil {
			r.Errorf("listFindingsWithMark(%s) adding marks had error: %v", findingName, err)
			return
		}

		err = listFindingsWithMarks(buf, sourceName)

		if err != nil {
			r.Errorf("listFindingsWithMark(%s) had error: %v", sourceName, err)
			return
		}

		got := buf.String()
		if !strings.Contains(got, untouchedFindingName) {
			r.Errorf("listFindingWithMarks(%s) got: %s want %s", sourceName, got, untouchedFindingName)
		}

		if strings.Contains(got, findingName) {
			r.Errorf("listFindingWithMarks(%s) got: %s didn't want %s", orgID, got, findingName)
		}

	})
}

func TestGroupFindings(t *testing.T) {
	orgID := setup(t)
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		err := groupFindings(buf, orgID)

		if err != nil {
			r.Errorf("groupFindings(%s) had error: %v", orgID, err)
			return
		}

		got := buf.String()
		if want := "Grouped Findings"; !strings.Contains(got, want) {
			r.Errorf("groupFindings(%s) got: %s want %s", orgID, got, want)
		}
	})
}

func TestGroupFindingsWithFilter(t *testing.T) {
	setup(t)
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		err := groupFindingsWithFilter(buf, sourceName)

		if err != nil {
			r.Errorf("groupFindingsWithFilter(%s) had error: %v", sourceName, err)
			return
		}

		got := buf.String()
		if want := "Grouped Findings"; !strings.Contains(got, want) {
			r.Errorf("groupFindingsWithFilter(%s) got: %s want %s", sourceName, got, want)
		}
	})
}

func TestGroupFindingsByState(t *testing.T) {
	setup(t)
	testutil.Retry(t, 5, 5*time.Second, func(r *testutil.R) {
		buf := new(bytes.Buffer)

		err := groupFindingsByState(buf, sourceName)

		if err != nil {
			r.Errorf("groupFindingsByState(%s) had error: %v", sourceName, err)
			return
		}

		got := buf.String()
		if want := "Grouped Findings by State"; !strings.Contains(got, want) {
			r.Errorf("groupFindingsByState(%s) got: %s want %s", sourceName, got, want)
		}
	})
}
