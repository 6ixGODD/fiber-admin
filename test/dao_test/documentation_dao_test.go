package dao_test

import (
	"testing"
	"time"

	"fiber-admin/test/wire"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var documentID primitive.ObjectID

func TestInsertDocumentation(t *testing.T) {
	// t.Skip("Skip TestInsertDocumentation")
	var (
		injector         = wire.GetInjector()
		documentationDao = injector.DocumentationDao
		ctx              = injector.Ctx
		title            = "Title"
		content          = "Content"
		err              error
	)

	documentID, err = documentationDao.InsertDocumentation(ctx, title, content)
	assert.NoError(t, err)
	assert.NotEmpty(t, documentID)

	documentation, err := documentationDao.GetDocumentationByID(ctx, documentID)
	assert.NoError(t, err)
	assert.NotNil(t, documentation)
	assert.Equal(t, title, documentation.Title)
	assert.Equal(t, content, documentation.Content)
}

func TestGetDocumentation(t *testing.T) {
	// t.Skip("Skip TestGetDocumentation")
	var (
		injector         = wire.GetInjector()
		documentationDao = injector.DocumentationDao
		ctx              = injector.Ctx
		err              error
	)
	documentation, err := documentationDao.GetDocumentationByID(ctx, documentID)
	assert.NoError(t, err)
	assert.NotNil(t, documentation)
	assert.NotEmpty(t, documentation.DocumentID)
	assert.NotEmpty(t, documentation.Title)
	assert.NotEmpty(t, documentation.Content)
	assert.NotEmpty(t, documentation.CreatedAt)
	assert.NotEmpty(t, documentation.UpdatedAt)
}

func TestGetDocumentationList(t *testing.T) {
	// t.Skip("Skip TestGetDocumentationList")
	var (
		injector         = wire.GetInjector()
		documentationDao = injector.DocumentationDao
		ctx              = injector.Ctx
		createStartTime  = time.Now().Add(-time.Hour)
		createEndTime    = time.Now().Add(time.Hour)
		updateStartTime  = time.Now().Add(-time.Hour)
		updateEndTime    = time.Now().Add(time.Hour)
		err              error
	)
	documentationList, count, err := documentationDao.GetDocumentationList(
		ctx, 0, 10, false, nil, nil, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, count)
	assert.NotEmpty(t, documentationList)
	t.Logf("Documentation Count: %d", *count)
	t.Logf("Documentation List: %v", documentationList)
	t.Logf("=====================================")

	documentationList, count, err = documentationDao.GetDocumentationList(
		ctx, 0, 10, false, &createStartTime, &createEndTime, nil, nil,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, count)
	t.Logf("Create Start Time: %v", createStartTime)
	t.Logf("Create End Time: %v", createEndTime)
	t.Logf("Documentation Count: %d", *count)
	t.Logf("Documentation List: %v", documentationList)
	t.Logf("=====================================")

	documentationList, count, err = documentationDao.GetDocumentationList(
		ctx, 0, 10, false, nil, nil, &updateStartTime, &updateEndTime,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, count)
	t.Logf("Update Start Time: %v", updateStartTime)
	t.Logf("Update End Time: %v", updateEndTime)
	t.Logf("Documentation Count: %d", *count)
	t.Logf("Documentation List: %v", documentationList)
	t.Logf("=====================================")

	documentationList, count, err = documentationDao.GetDocumentationList(
		ctx, 0, 10, false, &createStartTime, &createEndTime, &updateStartTime, &updateEndTime,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, count)
	t.Logf("Create Start Time: %v", createStartTime)
	t.Logf("Create End Time: %v", createEndTime)
	t.Logf("Update Start Time: %v", updateStartTime)
	t.Logf("Update End Time: %v", updateEndTime)
	t.Logf("Documentation Count: %d", *count)
	t.Logf("Documentation List: %v", documentationList)
	t.Logf("=====================================")
}

func TestUpdateDocumentation(t *testing.T) {
	// t.Skip("Skip TestUpdateDocumentation")
	var (
		injector         = wire.GetInjector()
		documentationDao = injector.DocumentationDao
		ctx              = injector.Ctx
		title            = "New Title"
		content          = "New Content"
		err              error
	)

	err = documentationDao.UpdateDocumentation(ctx, documentID, &title, &content)
	assert.NoError(t, err)

	documentation, err := documentationDao.GetDocumentationByID(ctx, documentID)
	assert.NoError(t, err)
	assert.NotNil(t, documentation)
	assert.Equal(t, title, documentation.Title)
	assert.Equal(t, content, documentation.Content)
}

func TestDeleteDocumentation(t *testing.T) {
	// t.Skip("Skip TestDeleteDocumentation")
	var (
		injector         = wire.GetInjector()
		documentationDao = injector.DocumentationDao
		ctx              = injector.Ctx
		err              error
	)
	err = documentationDao.DeleteDocumentation(ctx, documentID)
	assert.NoError(t, err)

	documentation, err := documentationDao.GetDocumentationByID(ctx, documentID)
	assert.Error(t, err)
	assert.Nil(t, documentation)
}

func TestDeleteDocumentationList(t *testing.T) {
	// t.Skip("Skip TestDeleteDocumentationList")
	var (
		injector         = wire.GetInjector()
		documentationDao = injector.DocumentationDao
		ctx              = injector.Ctx
		createStartTime  = time.Now().Add(-time.Hour)
		createEndTime    = time.Now().Add(time.Hour)
		updateStartTime  = time.Now().Add(-time.Hour)
		updateEndTime    = time.Now().Add(time.Hour)
		err              error
	)

	documentationList, count, err := documentationDao.GetDocumentationList(
		ctx, 0, 10, false, &createStartTime, &createEndTime, &updateStartTime, &updateEndTime,
	)
	assert.NoError(t, err)
	assert.NotEmpty(t, count)
	assert.NotEmpty(t, documentationList)
	t.Logf("Documentation Count: %d", *count)
	t.Logf("Documentation List: %v", documentationList)
	t.Logf("=====================================")

	count, err = documentationDao.DeleteDocumentationList(
		ctx, &createStartTime, &createEndTime, &updateStartTime, &updateEndTime,
	)
	assert.NoError(t, err)
	assert.NotNil(t, count)
	t.Logf("Documentation Count: %d", *count)
	t.Logf("=====================================")

	documentationList, count, err = documentationDao.GetDocumentationList(
		ctx, 0, 10, false, &createStartTime, &createEndTime, &updateStartTime, &updateEndTime,
	)
	assert.NoError(t, err)
	assert.Empty(t, documentationList)
	t.Logf("Documentation Count: %d", *count)
	t.Logf("Documentation List: %v", documentationList)
	t.Logf("=====================================")
}
