package integration

import (
	"avito-intern-test-task-2025/internal/config"
	"avito-intern-test-task-2025/internal/entity/repo"
	"avito-intern-test-task-2025/internal/http/dto"
	"avito-intern-test-task-2025/internal/http/handlers"
	httproutes "avito-intern-test-task-2025/internal/http/routes"
	"avito-intern-test-task-2025/internal/usecase"
	"avito-intern-test-task-2025/pkg/ServiceDependencies"
	"avito-intern-test-task-2025/pkg/db"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRouter initializes a test router with all dependencies
func setupTestRouter(t *testing.T) (*gin.Engine, *ServiceDependencies.ServiceDependencies) {
	// Set test environment
	os.Setenv("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=avito_test port=5432 sslmode=disable")
	os.Setenv("ENV", "test")

	c := config.LoadConfig()
	database, err := db.InitDB(c)
	if err != nil {
		t.Skipf("Skipping test: database not available - %v", err)
	}

	// Run migrations
	database.RunMigrations()

	r := gin.Default()
	s := &ServiceDependencies.ServiceDependencies{
		DB: database,
	}

	// Setup repositories
	tr := repo.NewTeamRepo(s.DB)
	ur := repo.NewUserRepo(s.DB)
	pr := repo.NewPrRepo(s.DB, ur)

	// Setup use cases
	uc := usecase.NewUserUsecase(ur, pr)
	pc := usecase.NewPullRequestUsecase(ur, pr)
	tc := usecase.NewTeamUsecase(tr, ur, pr, pc)
	sc := usecase.NewStatsUsecase(pr)

	// Setup handlers
	uh := handlers.NewUserHandler(uc)
	th := handlers.NewTeamHandler(tc)
	ph := handlers.NewPrHandler(pc)
	sh := handlers.NewStatsHandler(sc)

	// Register routes
	v1 := r.Group("/v1")
	api := v1.Group("/api")
	api.GET("/health", handlers.Health())

	httproutes.RegisterUserRoutes(api, uh)
	httproutes.RegisterTeamRoutes(api, th)
	httproutes.RegisterPullRequestRoutes(api, ph)
	httproutes.RegisterStatsRoutes(api, sh)

	return r, s
}

// TestIntegrationFullWorkflow tests the complete workflow:
// 1. Create a team
// 2. Add users to team
// 3. Create pull requests
// 4. Assign reviewers
// 5. Check stats
// 6. Merge PR
// 7. Bulk deactivate and reassign
func TestIntegrationFullWorkflow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, s := setupTestRouter(t)
	if s == nil {
		return
	}
	defer s.DB.Close()

	// Step 1: Create a team
	teamPayload := dto.TeamDTO{
		TeamName: "test-team-integration",
		Members: []dto.TeamMemberDTO{
			{
				UserId:   "user-int-1",
				Username: "john_dev",
				IsActive: true,
			},
			{
				UserId:   "user-int-2",
				Username: "jane_dev",
				IsActive: true,
			},
			{
				UserId:   "user-int-3",
				Username: "bob_qa",
				IsActive: true,
			},
		},
	}

	bodyBytes, _ := json.Marshal(teamPayload)
	req := httptest.NewRequest(http.MethodPost, "/v1/api/team/add", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Failed to create team")

	// Step 2: Verify team was created
	req = httptest.NewRequest(http.MethodGet, "/v1/api/team/get?team_name=test-team-integration", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Failed to get team")

	// Step 3: Create pull requests
	prPayload := dto.PullRequestDTO{
		PullRequestId:   "pr-int-1",
		PullRequestName: "Fix bug in auth",
		AuthorId:        "user-int-1",
		Status:          "OPEN",
	}

	bodyBytes, _ = json.Marshal(prPayload)
	req = httptest.NewRequest(http.MethodPost, "/v1/api/pullRequest/create", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, "Failed to create PR")

	// Create second PR
	prPayload2 := dto.PullRequestDTO{
		PullRequestId:   "pr-int-2",
		PullRequestName: "Add new feature",
		AuthorId:        "user-int-2",
		Status:          "OPEN",
	}

	bodyBytes, _ = json.Marshal(prPayload2)
	req = httptest.NewRequest(http.MethodPost, "/v1/api/pullRequest/create", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, "Failed to create second PR")

	// Step 4: Assign reviewers to PR 1
	req = httptest.NewRequest(http.MethodGet, "/v1/api/users/getReview?pull_request_id=pr-int-1&team_name=test-team-integration", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Failed to get reviewers")

	// Verify PR has reviewers assigned
	assert.NotEmpty(t, w.Body.String(), "Response should not be empty")

	// Step 5: Get statistics
	req = httptest.NewRequest(http.MethodGet, "/v1/api/stats/assignments", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Failed to get stats")
	assert.Contains(t, w.Body.String(), "stats", "Stats response should contain stats key")

	var statsResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &statsResp)
	assert.NotNil(t, statsResp["stats"], "Stats should not be nil")

	// Step 6: Merge a PR
	mergePayload := dto.PullRequestDTO{
		PullRequestId: "pr-int-1",
	}
	bodyBytes, _ = json.Marshal(mergePayload)
	req = httptest.NewRequest(http.MethodPost, "/v1/api/pullRequest/merge", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Failed to merge PR")

	// Step 7: Bulk deactivate users
	bulkDeactivatePayload := dto.TeamBulkDeactivateDTO{
		TeamName: "test-team-integration",
		UserIds:  []string{"user-int-2"},
	}
	bodyBytes, _ = json.Marshal(bulkDeactivatePayload)
	req = httptest.NewRequest(http.MethodPost, "/v1/api/team/bulkDeactivate", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Failed to bulk deactivate")

	// Verify result contains reassignment information
	var deactivateResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &deactivateResp)
	assert.NotNil(t, deactivateResp["result"], "Result should not be nil")

	// Step 8: Verify stats after bulk deactivation
	req = httptest.NewRequest(http.MethodGet, "/v1/api/stats/assignments", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Failed to get final stats")

	t.Logf("Integration test completed successfully")
}

// TestIntegrationTeamWorkflow tests team creation and management
func TestIntegrationTeamWorkflow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, s := setupTestRouter(t)
	if s == nil {
		return
	}
	defer s.DB.Close()

	// Create multiple teams
	teams := []dto.TeamDTO{
		{
			TeamName: "frontend-team",
			Members: []dto.TeamMemberDTO{
				{UserId: "fe-1", Username: "alice", IsActive: true},
				{UserId: "fe-2", Username: "bob", IsActive: true},
			},
		},
		{
			TeamName: "backend-team",
			Members: []dto.TeamMemberDTO{
				{UserId: "be-1", Username: "charlie", IsActive: true},
				{UserId: "be-2", Username: "david", IsActive: true},
			},
		},
	}

	for _, team := range teams {
		bodyBytes, _ := json.Marshal(team)
		req := httptest.NewRequest(http.MethodPost, "/v1/api/team/add", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code, "Failed to create team: %s", team.TeamName)
	}

	// Verify both teams can be retrieved
	for _, team := range teams {
		req := httptest.NewRequest(http.MethodGet, "/v1/api/team/get?team_name="+team.TeamName, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Failed to retrieve team: %s", team.TeamName)
	}

	t.Log("Team workflow test passed")
}

// TestIntegrationReviewerAssignment tests reviewer assignment and reassignment
func TestIntegrationReviewerAssignment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, s := setupTestRouter(t)
	if s == nil {
		return
	}
	defer s.DB.Close()

	// Setup: Create team and PR
	teamPayload := dto.TeamDTO{
		TeamName: "review-test-team",
		Members: []dto.TeamMemberDTO{
			{UserId: "reviewer-1", Username: "reviewer1", IsActive: true},
			{UserId: "reviewer-2", Username: "reviewer2", IsActive: true},
			{UserId: "reviewer-3", Username: "reviewer3", IsActive: true},
		},
	}

	bodyBytes, _ := json.Marshal(teamPayload)
	req := httptest.NewRequest(http.MethodPost, "/v1/api/team/add", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Create PR
	prPayload := dto.PullRequestDTO{
		PullRequestId:   "pr-review-test",
		PullRequestName: "Review assignment test",
		AuthorId:        "reviewer-1",
		Status:          "OPEN",
	}

	bodyBytes, _ = json.Marshal(prPayload)
	req = httptest.NewRequest(http.MethodPost, "/v1/api/pullRequest/create", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Get reviewer assignment
	req = httptest.NewRequest(http.MethodGet, "/v1/api/users/getReview?pull_request_id=pr-review-test&team_name=review-test-team", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should get reviewers")

	// Reassign reviewer
	reassignPayload := dto.PullRequestReassignDTO{
		PullRequestId: "pr-review-test",
		OldUserId:     "reviewer-1",
	}

	bodyBytes, _ = json.Marshal(reassignPayload)
	req = httptest.NewRequest(http.MethodPost, "/v1/api/pullRequest/reassign", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should reassign reviewer")

	t.Log("Reviewer assignment test passed")
}

// TestIntegrationStatistics tests statistics endpoint
func TestIntegrationStatistics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, s := setupTestRouter(t)
	if s == nil {
		return
	}
	defer s.DB.Close()

	// Setup: Create team with users and multiple PRs
	teamPayload := dto.TeamDTO{
		TeamName: "stats-team",
		Members: []dto.TeamMemberDTO{
			{UserId: "stats-user-1", Username: "user1", IsActive: true},
			{UserId: "stats-user-2", Username: "user2", IsActive: true},
			{UserId: "stats-user-3", Username: "user3", IsActive: true},
		},
	}

	bodyBytes, _ := json.Marshal(teamPayload)
	req := httptest.NewRequest(http.MethodPost, "/v1/api/team/add", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Create multiple PRs
	for i := 1; i <= 3; i++ {
		prPayload := dto.PullRequestDTO{
			PullRequestId:   "stats-pr-" + string(rune(i+'0')),
			PullRequestName: "Test PR " + string(rune(i+'0')),
			AuthorId:        "stats-user-1",
			Status:          "OPEN",
		}

		bodyBytes, _ := json.Marshal(prPayload)
		req := httptest.NewRequest(http.MethodPost, "/v1/api/pullRequest/create", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
	}

	// Get stats
	req = httptest.NewRequest(http.MethodGet, "/v1/api/stats/assignments", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Stats endpoint should return 200")

	// Verify response structure
	var statsResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &statsResp)
	require.NoError(t, err, "Response should be valid JSON")

	stats, ok := statsResp["stats"]
	assert.True(t, ok, "Response should contain 'stats' key")
	assert.NotNil(t, stats, "Stats should not be nil")

	t.Log("Statistics test passed")
}

// TestIntegrationBulkDeactivationWithReassignment tests bulk deactivation and PR reassignment
func TestIntegrationBulkDeactivationWithReassignment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, s := setupTestRouter(t)
	if s == nil {
		return
	}
	defer s.DB.Close()

	// Setup: Create team with users
	teamPayload := dto.TeamDTO{
		TeamName: "deactivation-team",
		Members: []dto.TeamMemberDTO{
			{UserId: "deact-user-1", Username: "user1", IsActive: true},
			{UserId: "deact-user-2", Username: "user2", IsActive: true},
			{UserId: "deact-user-3", Username: "user3", IsActive: true},
			{UserId: "deact-user-4", Username: "user4", IsActive: true},
		},
	}

	bodyBytes, _ := json.Marshal(teamPayload)
	req := httptest.NewRequest(http.MethodPost, "/v1/api/team/add", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Create open PR
	prPayload := dto.PullRequestDTO{
		PullRequestId:   "deact-pr-1",
		PullRequestName: "Open PR for deactivation test",
		AuthorId:        "deact-user-1",
		Status:          "OPEN",
	}

	bodyBytes, _ = json.Marshal(prPayload)
	req = httptest.NewRequest(http.MethodPost, "/v1/api/pullRequest/create", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Assign reviewers
	req = httptest.NewRequest(http.MethodGet, "/v1/api/users/getReview?pull_request_id=deact-pr-1&team_name=deactivation-team", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Bulk deactivate one user
	bulkDeactivatePayload := dto.TeamBulkDeactivateDTO{
		TeamName: "deactivation-team",
		UserIds:  []string{"deact-user-1"},
	}

	bodyBytes, _ = json.Marshal(bulkDeactivatePayload)
	req = httptest.NewRequest(http.MethodPost, "/v1/api/team/bulkDeactivate", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Bulk deactivation should succeed")

	// Verify result
	var result map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.NotNil(t, result["result"], "Should have result")

	t.Log("Bulk deactivation with reassignment test passed")
}

// TestIntegrationUserActivationToggle tests setting user active status
func TestIntegrationUserActivationToggle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router, s := setupTestRouter(t)
	if s == nil {
		return
	}
	defer s.DB.Close()

	// Create team with user
	teamPayload := dto.TeamDTO{
		TeamName: "activation-team",
		Members: []dto.TeamMemberDTO{
			{UserId: "act-user-1", Username: "testuser", IsActive: true},
		},
	}

	bodyBytes, _ := json.Marshal(teamPayload)
	req := httptest.NewRequest(http.MethodPost, "/v1/api/team/add", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Deactivate user
	deactivatePayload := dto.UserIsActiveDTO{
		UserId:   "act-user-1",
		IsActive: false,
	}

	bodyBytes, _ = json.Marshal(deactivatePayload)
	req = httptest.NewRequest(http.MethodPost, "/v1/api/users/setIsActive", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should deactivate user")

	// Reactivate user
	activatePayload := dto.UserIsActiveDTO{
		UserId:   "act-user-1",
		IsActive: true,
	}

	bodyBytes, _ = json.Marshal(activatePayload)
	req = httptest.NewRequest(http.MethodPost, "/v1/api/users/setIsActive", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should reactivate user")

	t.Log("User activation toggle test passed")
}
