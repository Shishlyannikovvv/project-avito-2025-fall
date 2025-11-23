package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080"

func TestFullFlow(t *testing.T) {
	// 1. Создаём команду
	team := createTeam(t, `{"name":"gophers"}`)
	teamID := int(team["id"].(float64))

	// 2. Создаём пользователей
	alice := createUser(t, `{"username":"alice","team_id":`+strconv.Itoa(teamID)+`}`)
	bob := createUser(t, `{"username":"bob","team_id":`+strconv.Itoa(teamID)+`}`)
	charlie := createUser(t, `{"username":"charlie","team_id":`+strconv.Itoa(teamID)+`}`)

	aliceID := int(alice["id"].(float64))
	bobID := int(bob["id"].(float64))
	charlieID := int(charlie["id"].(float64))

	// 3. Создаём PR
	pr := createPR(t, `{"title":"fix bug","author_id":`+strconv.Itoa(aliceID)+`}`)
	prMap := pr // уже map[string]any

	// Безопасно достаём ревьюверов
	rev1 := getReviewerID(t, prMap, "reviewer1_id")
	rev2 := getReviewerID(t, prMap, "reviewer2_id")

	assert.Contains(t, []int{rev1, rev2}, bobID)
	assert.Contains(t, []int{rev1, rev2}, charlieID)

	// 4. Reassign
	newPR := reassignReviewer(t, int(prMap["id"].(float64)), rev1)
	newRev1 := getReviewerID(t, newPR, "reviewer1_id")
	assert.NotEqual(t, rev1, newRev1)

	// 5. Merge
	merged := mergePR(t, int(prMap["id"].(float64)))
	assert.Equal(t, "MERGED", merged["status"])
}

// Вспомогательные функции
func createTeam(t *testing.T, body string) map[string]any {
	return postJSON(t, "/teams", body)
}

func createUser(t *testing.T, body string) map[string]any {
	return postJSON(t, "/users", body)
}

func createPR(t *testing.T, body string) map[string]any {
	return postJSON(t, "/pull-requests", body)
}

func reassignReviewer(t *testing.T, prID, reviewerID int) map[string]any {
	payload := map[string]int{"reviewer_id": reviewerID}
	jsonBody, _ := json.Marshal(payload)
	resp := patch(t, "/pull-requests/"+strconv.Itoa(prID)+"/reassign", bytes.NewBuffer(jsonBody))
	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}

func mergePR(t *testing.T, prID int) map[string]any {
	resp := patch(t, "/pull-requests/"+strconv.Itoa(prID)+"/merge", nil)
	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	return result
}

func postJSON(t *testing.T, path, body string) map[string]any {
	resp := post(t, path, body)
	var result map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	return result
}

func post(t *testing.T, path, body string) *http.Response {
	resp, err := http.Post(baseURL+path, "application/json", bytes.NewBufferString(body))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	return resp
}

func patch(t *testing.T, path string, body *bytes.Buffer) *http.Response {
	req, _ := http.NewRequest("PATCH", baseURL+path, body)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	return resp
}

// Безопасно достаём ID ревьювера (может быть null)
func getReviewerID(t *testing.T, m map[string]any, key string) int {
	if val, ok := m[key]; ok && val != nil {
		return int(val.(float64))
	}
	return 0 // не назначен
}
