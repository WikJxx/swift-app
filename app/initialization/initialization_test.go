package initialization

import (
	"context"
	"path/filepath"
	"runtime"
	"swift-app/database"
	utils "swift-app/internal/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeDatabase_WithTestContainer(t *testing.T) {
	ctx := context.Background()

	uri, err := utils.MongoContainer.ConnectionString(ctx)
	assert.NoError(t, err)

	err = InitializeDatabase(uri, "swiftTestDB", "swiftCodes")
	assert.NoError(t, err)

	assert.NotNil(t, database.GetCollection())
}

func TestImportData_WithTestContainer(t *testing.T) {
	ctx := context.Background()

	uri, err := utils.MongoContainer.ConnectionString(ctx)
	assert.NoError(t, err)

	err = InitializeDatabase(uri, "swiftTestDB", "swiftCodes")
	assert.NoError(t, err)

	_, currentFilePath, _, _ := runtime.Caller(0)
	testCSV := filepath.Join(filepath.Dir(currentFilePath), "test_data", "swift_test.csv")

	hqAdded, hqSkipped, branchesAdded, branchesDuplicate, branchesMissingHQ, branchesSkipped, err := ImportData(testCSV)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, hqAdded, 1)
	assert.Equal(t, 0, hqSkipped)
	assert.GreaterOrEqual(t, branchesAdded, 0)
	assert.Equal(t, 0, branchesDuplicate)
	assert.Equal(t, 1, branchesMissingHQ)
	assert.Equal(t, 1, branchesSkipped)

	t.Cleanup(func() {
		collection := database.GetCollection()
		_, _ = collection.DeleteMany(ctx, struct{}{})
		_ = database.CloseMongoDB()
	})
}
