package e2e

import (
	"fmt"
	"os"
	"testing"
)

func TestE2E(t *testing.T) {
	if os.Getenv("PA_E2E_TEST_RUN") != "ON" {
		msg := `E2E test skip.
If you run E2E test, Set below environment variables.

- PA_E2E_TEST_RUN=ON
- PA_USERNAME=<pixela-username-for-testing>
- PA_FIRST_TOKEN=<pixela-token-for-testing>
- PA_SECOND_TOKEN=<pixela-token-for-testing>`
		fmt.Println(msg)
		return
	}

	testE2EUserCreate(t)
	testE2EUserUpdate(t)

	testE2EUserProfileUpdate(t)

	testE2EGraphCreate(t)
	testE2EGraphGetAll(t)
	testE2EGraphGetSVG(t)
	testE2EGraphStats(t)
	testE2EGraphUpdate(t)
	testE2EGraphGetPixelDates(t)
	testE2EGraphStopwatch(t)

	testE2EPixelCreate(t)
	testE2EPixelIncrement(t)
	testE2EPixelDecrement(t)
	testE2EPixelGet(t)
	testE2EPixelUpdate(t)
	testE2EPixelDelete(t)

	testE2EWebhookCreate(t)
	testE2EWebhookGetAll(t)
	testE2EWebhookInvoke(t)
	testE2EWebhookDelete(t)

	testE2EGraphDelete(t)
	testE2EUserDelete(t)
}
