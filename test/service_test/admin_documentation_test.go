package service_test

import (
	"testing"

	"fiber-admin/test/mock"
	"fiber-admin/test/wire"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInsertDocumentation(t *testing.T) {
	var (
		injector             = wire.GetInjector()
		ctx                  = injector.Ctx
		documentationService = injector.AdminDocumentationService
		title                = mock.RandomString(10)
		content              = mock.RandomString(10)
	)
	documentationIDHex, err := documentationService.InsertDocumentation(ctx, &title, &content)
	assert.NoError(t, err)
	assert.NotNil(t, documentationIDHex)

	t.Logf("Documentation ID: %s", documentationIDHex)

	documentationID, err := primitive.ObjectIDFromHex(documentationIDHex)
	assert.NoError(t, err)
	documentation, err := injector.DocumentationDao.GetDocumentationByID(ctx, documentationID)
	assert.NoError(t, err)
	assert.NotNil(t, documentation)
	assert.Equal(t, title, documentation.Title)
	assert.Equal(t, content, documentation.Content)

	t.Logf("Documentation Data: %+v", *documentation)
}

func TestUpdateDocumentation(t *testing.T) {
	var (
		injector             = wire.GetInjector()
		ctx                  = injector.Ctx
		documentationService = injector.AdminDocumentationService
		documentationID      = injector.DocumentationDaoMock.RandomDocumentationID()
		title                = mock.RandomString(10)
		content              = mock.RandomString(10)
	)
	err := documentationService.UpdateDocumentation(ctx, &documentationID, &title, &content)
	assert.NoError(t, err)

	documentation, err := injector.DocumentationDao.GetDocumentationByID(ctx, documentationID)
	assert.NoError(t, err)
	assert.NotNil(t, documentation)
	assert.Equal(t, title, documentation.Title)
	assert.Equal(t, content, documentation.Content)

	t.Logf("Documentation Data: %+v", documentation)
}

func TestDeleteDocumentation(t *testing.T) {
	var (
		injector             = wire.GetInjector()
		ctx                  = injector.Ctx
		documentationService = injector.AdminDocumentationService
		documentationID      = injector.DocumentationDaoMock.RandomDocumentationID()
	)
	err := documentationService.DeleteDocumentation(ctx, &documentationID)
	assert.NoError(t, err)

	documentation, err := injector.DocumentationDao.GetDocumentationByID(ctx, documentationID)
	assert.Error(t, err)
	assert.Nil(t, documentation)

	t.Logf("Documentation ID: %s", documentationID.Hex())
}
