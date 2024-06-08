package mock

import (
	"context"
	"math/rand"

	"fiber-admin/internal/pkg/dao/mods"
	"fiber-admin/internal/pkg/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocumentationDaoMock struct {
	DocumentationMap map[primitive.ObjectID]*entity.DocumentationModel
	DocumentationIDs []primitive.ObjectID
	DocumentationDao mods.DocumentationDao
}

func NewDocumentationDaoMock(documentationDao mods.DocumentationDao) *DocumentationDaoMock {
	return &DocumentationDaoMock{
		DocumentationMap: make(map[primitive.ObjectID]*entity.DocumentationModel),
		DocumentationDao: documentationDao,
	}
}

func NewDocumentationDaoMockWithRandomData(n int, documentationDao mods.DocumentationDao) *DocumentationDaoMock {
	documentationDaoMock := NewDocumentationDaoMock(documentationDao)
	for i := 0; i < n; i++ {
		documentation := documentationDaoMock.GenerateDocumentationModel()
		documentationDaoMock.DocumentationMap[documentation.DocumentID] = documentation
		documentationDaoMock.DocumentationIDs = append(documentationDaoMock.DocumentationIDs, documentation.DocumentID)
	}
	return documentationDaoMock
}

func (m *DocumentationDaoMock) Create(documentation *entity.DocumentationModel) {
	m.DocumentationMap[documentation.DocumentID] = documentation
}

func (m *DocumentationDaoMock) Get(documentationID primitive.ObjectID) (*entity.DocumentationModel, error) {
	documentation, ok := m.DocumentationMap[documentationID]
	if !ok {
		return nil, nil
	}
	return documentation, nil
}

func (m *DocumentationDaoMock) RandomDocumentationID() primitive.ObjectID {
	return m.DocumentationIDs[rand.Intn(len(m.DocumentationIDs))]
}

func (m *DocumentationDaoMock) GenerateDocumentationModel() *entity.DocumentationModel {
	title, content := GenerateDocumentation()
	documentationID, err := m.DocumentationDao.InsertDocumentation(context.Background(), title, content)
	if err != nil {
		panic(err)
	}

	documentation, err := m.DocumentationDao.GetDocumentationByID(context.Background(), documentationID)
	if err != nil {
		panic(err)
	}
	return documentation
}

func GenerateDocumentation() (string, string) {
	return RandomString(10), RandomString(10)
}

func (m *DocumentationDaoMock) Delete() {
	for _, documentationID := range m.DocumentationIDs {
		_ = m.DocumentationDao.DeleteDocumentation(context.Background(), documentationID)
	}
}
