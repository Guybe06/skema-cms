package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"skema-api/core/cache"
	"skema-api/core/mailer"
	"skema-api/features/accounts"
	accountsrepo "skema-api/features/accounts/repository"
	accountssvc "skema-api/features/accounts/service"
	"skema-api/features/apikeys"
	apikeysrepo "skema-api/features/apikeys/repository"
	apikeyssvc "skema-api/features/apikeys/service"
	"skema-api/features/collections"
	collectionsrepo "skema-api/features/collections/repository"
	collectionssvc "skema-api/features/collections/service"
	"skema-api/features/connections"
	connectionsrepo "skema-api/features/connections/repository"
	connectionssvc "skema-api/features/connections/service"
	"skema-api/features/content"
	contentsvc "skema-api/features/content/service"
	"skema-api/features/memberships"
	membershipsrepo "skema-api/features/memberships/repository"
	membershipssvc "skema-api/features/memberships/service"
	"skema-api/features/organizations"
	orgsrepo "skema-api/features/organizations/repository"
	orgssvc "skema-api/features/organizations/service"
	"skema-api/features/publicapi"
	"skema-api/features/users"
	usersrepo "skema-api/features/users/repository"
	userssvc "skema-api/features/users/service"
)

const testDBName = "skemacms_test"
const testEncKey = "test-encryption-key-32-chars-pad!"

var (
	testPool      *pgxpool.Pool
	testRouter    *gin.Engine
	testJWTSecret = "jwt-secret-pour-les-tests-uniquement"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../.env")
	gin.SetMode(gin.TestMode)

	ctx := context.Background()
	ensureTestDB(ctx)

	pool, err := pgxpool.New(ctx, buildDSN(testDBName))
	if err != nil {
		panic(fmt.Sprintf("connexion à la base de test échouée : %v", err))
	}
	testPool = pool

	runMigrations(ctx)
	testRouter = buildRouter()

	code := m.Run()
	pool.Close()
	os.Exit(code)
}

func ensureTestDB(ctx context.Context) {
	conn, err := pgx.Connect(ctx, buildDSN("postgres"))
	if err != nil {
		panic(fmt.Sprintf("connexion à postgres échouée : %v", err))
	}
	defer conn.Close(ctx)

	var exists bool
	conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)", testDBName).Scan(&exists)
	if !exists {
		conn.Exec(ctx, "CREATE DATABASE "+testDBName)
	}
}

func runMigrations(ctx context.Context) {
	files, err := os.ReadDir("../../migrations")
	if err != nil {
		panic(fmt.Sprintf("lecture des migrations échouée : %v", err))
	}

	conn, err := testPool.Acquire(ctx)
	if err != nil {
		panic(err)
	}
	defer conn.Release()

	for _, f := range files {
		if f.IsDir() || len(f.Name()) < 4 || f.Name()[len(f.Name())-6:] != "up.sql" {
			continue
		}
		sql, err := os.ReadFile("../../migrations/" + f.Name())
		if err != nil {
			panic(fmt.Sprintf("lecture de %s échouée : %v", f.Name(), err))
		}
		conn.Exec(ctx, string(sql))
	}
}

func buildRouter() *gin.Engine {
	r := gin.New()
	c := cache.New("")
	m := mailer.New("test@skemacms.com", "")

	v1 := r.Group("/v1")

	accountsRepo := accountsrepo.New(testPool)
	accountsSvc := accountssvc.NewForTest(accountsRepo, c, m, testJWTSecret, "http://localhost:3001")
	accounts.RegisterRoutes(v1, accountsSvc, testJWTSecret)

	users.RegisterRoutes(v1, userssvc.New(usersrepo.New(testPool)), testJWTSecret)

	orgsRepository := orgsrepo.New(testPool)
	organizations.RegisterRoutes(v1, orgssvc.New(orgsRepository), testJWTSecret)
	memberships.RegisterRoutes(v1, membershipssvc.New(membershipsrepo.New(testPool), orgsRepository, m, "http://localhost:3001"), testJWTSecret)

	connSvc := connectionssvc.New(connectionsrepo.New(testPool), orgsRepository, testEncKey)
	connections.RegisterRoutes(v1, connSvc, testJWTSecret)

	collRepo := collectionsrepo.New(testPool)
	collections.RegisterRoutes(v1, collectionssvc.New(collRepo, orgsRepository, connSvc), testJWTSecret)
	content.RegisterRoutes(v1, contentsvc.New(collRepo, orgsRepository, connSvc), testJWTSecret)

	apikeySvc := apikeyssvc.New(apikeysrepo.New(testPool), orgsRepository)
	apikeys.RegisterRoutes(v1, apikeySvc, testJWTSecret)

	pub := v1.Group("/pub")
	publicapi.RegisterRoutes(pub, apikeySvc, orgsRepository, collRepo, connSvc)

	return r
}

func buildDSN(dbName string) string {
	host := getenv("CMS_DB_HOST", "localhost")
	port := getenv("CMS_DB_PORT", "5432")
	user := getenv("CMS_DB_USER", "postgres")
	pass := getenv("CMS_DB_PASSWORD", "")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbName)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func dbPort() int {
	v := getenv("CMS_DB_PORT", "5432")
	n := 5432
	fmt.Sscanf(v, "%d", &n)
	return n
}

func truncateTables(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	// Supprimer les tables dynamiques créées par les tests (préfixe test_)
	rows, err := testPool.Query(ctx,
		`SELECT tablename FROM pg_tables WHERE schemaname='public' AND tablename LIKE 'test_%'`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name string
			rows.Scan(&name)
			testPool.Exec(ctx, fmt.Sprintf(`DROP TABLE IF EXISTS "%s" CASCADE`, name))
		}
	}

	_, err = testPool.Exec(ctx,
		`TRUNCATE users, sessions, verification_tokens, organizations, memberships, connections,
		         collections, collection_fields, api_keys RESTART IDENTITY CASCADE`,
	)
	if err != nil {
		t.Fatalf("nettoyage des tables échoué : %v", err)
	}
}

func do(method, path string, body any, token string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func doWithKey(method, path string, body any, apiKey string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("X-Api-Key", apiKey)
	}
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func decode(t *testing.T, w *httptest.ResponseRecorder, target any) {
	t.Helper()
	if err := json.NewDecoder(w.Body).Decode(target); err != nil {
		t.Fatalf("décodage de la réponse échoué : %v", err)
	}
}

func assertStatus(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if w.Code != expected {
		t.Fatalf("statut attendu %d, obtenu %d - corps : %s", expected, w.Code, w.Body.String())
	}
}

func logResponse(t *testing.T, label string, w *httptest.ResponseRecorder) {
	t.Helper()
	t.Logf(">>> %s [%d]\n%s", label, w.Code, w.Body.String())
}

func signupAndGetToken(t *testing.T) string {
	t.Helper()
	w := do(http.MethodPost, "/v1/accounts/signup", map[string]string{
		"email":      "test@skemacms.com",
		"password":   "TestP@ss123!",
		"first_name": "Test",
		"last_name":  "User",
	}, "")
	assertStatus(t, w, http.StatusCreated)

	var resp struct {
		Data struct {
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	decode(t, w, &resp)
	return resp.Data.AccessToken
}
