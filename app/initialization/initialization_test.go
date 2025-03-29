package initialization

import (
	"context"
	"path/filepath"
	"runtime"
	"swift-app/database"
	testutils "swift-app/internal/testutils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeDatabase_WithTestContainer(t *testing.T) {
	ctx := context.Background()

	uri, err := testutils.MongoContainer.ConnectionString(ctx)
	assert.NoError(t, err)

	err = InitializeDatabase(uri, "swiftTestDB", "swiftCodes")
	assert.NoError(t, err)

	assert.NotNil(t, database.GetCollection())
}

func TestImportData_WithTestContainer(t *testing.T) {
	ctx := context.Background()

	uri, err := testutils.MongoContainer.ConnectionString(ctx)
	assert.NoError(t, err)

	err = InitializeDatabase(uri, "swiftTestDB", "swiftCodes")
	assert.NoError(t, err)

	_, currentFilePath, _, _ := runtime.Caller(0)
	testCSV := filepath.Join(filepath.Dir(currentFilePath), "test_data", "swift_test.csv")

	summary, err := ImportData(testCSV)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, summary.HQAdded, 1)
	assert.Equal(t, 0, summary.HQSkipped)
	assert.GreaterOrEqual(t, summary.BranchesAdded, 0)
	assert.Equal(t, 0, summary.BranchesDuplicate)
	assert.Equal(t, 1, summary.BranchesMissingHQ)
	assert.Equal(t, 1, summary.BranchesSkipped)

	t.Cleanup(func() {
		collection := database.GetCollection()
		_, _ = collection.DeleteMany(ctx, struct{}{})
		_ = database.CloseMongoDB()
	})
}
