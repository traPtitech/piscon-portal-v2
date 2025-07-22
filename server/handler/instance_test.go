package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/handler/openapi"
	repomock "github.com/traPtitech/piscon-portal-v2/server/repository/mock"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	usecasemock "github.com/traPtitech/piscon-portal-v2/server/usecase/mock"
	"github.com/traPtitech/piscon-portal-v2/server/utils/ptr"
	"go.uber.org/mock/gomock"
)

func TestGetTeamInstances(t *testing.T) {
	teamID := uuid.New()
	instances := []domain.Instance{
		{
			ID:     uuid.New(),
			TeamID: teamID,
			Index:  1,
			Infra: domain.InfraInstance{
				PublicIP:  ptr.Of("1.2.3.4"),
				PrivateIP: ptr.Of("10.0.0.1"),
				Status:    domain.InstanceStatusRunning,
			},
			CreatedAt: time.Now(),
		},
	}
	tests := []struct {
		name         string
		mockSetup    func(*usecasemock.MockUseCase)
		expectedCode int
		expectedLen  int
		expectErr    bool
	}{
		{
			name: "success",
			mockSetup: func(m *usecasemock.MockUseCase) {
				m.EXPECT().GetTeamInstances(gomock.Any(), teamID).Return(instances, nil)
			},
			expectedCode: http.StatusOK,
			expectedLen:  1,
			expectErr:    false,
		},
		{
			name: "error",
			mockSetup: func(m *usecasemock.MockUseCase) {
				m.EXPECT().GetTeamInstances(gomock.Any(), teamID).Return(nil, errors.New("unexpected error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedLen:  0,
			expectErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repoMock := repomock.NewMockRepository(ctrl)
			useCaseMock := usecasemock.NewMockUseCase(ctrl)
			e := echo.New()
			tt.mockSetup(useCaseMock)

			req := httptest.NewRequest(http.MethodGet, "/teams/"+teamID.String()+"/instances", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("teamID")
			c.SetParamValues(teamID.String())
			h := NewHandler(useCaseMock, repoMock, nil)

			err := h.GetTeamInstances(c)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			if tt.expectedCode == http.StatusOK {
				var res []openapi.Instance
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
				assert.Len(t, res, tt.expectedLen)
				if tt.expectedLen > 0 {
					assert.Equal(t, instances[0].ID, uuid.UUID(res[0].ID))
				}
			}
		})
	}
}

func TestCreateTeamInstance(t *testing.T) {
	teamID := uuid.New()
	instance := domain.Instance{
		ID:     uuid.New(),
		TeamID: teamID,
		Index:  1,
		Infra: domain.InfraInstance{
			PublicIP:  ptr.Of("1.2.3.4"),
			PrivateIP: ptr.Of("10.0.0.1"),
			Status:    domain.InstanceStatusRunning,
		},
		CreatedAt: time.Now(),
	}
	tests := []struct {
		name         string
		mockSetup    func(*usecasemock.MockUseCase)
		expectedCode int
		expectID     uuid.UUID
	}{
		{
			name: "success",
			mockSetup: func(m *usecasemock.MockUseCase) {
				m.EXPECT().CreateInstance(gomock.Any(), teamID).Return(instance, nil)
			},
			expectedCode: http.StatusCreated,
			expectID:     instance.ID,
		},
		{
			name: "bad request",
			mockSetup: func(m *usecasemock.MockUseCase) {
				m.EXPECT().CreateInstance(gomock.Any(), teamID).
					Return(domain.Instance{}, usecase.NewUseCaseError(domain.ErrInstanceLimitExceeded))
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "internal error",
			mockSetup: func(m *usecasemock.MockUseCase) {
				m.EXPECT().CreateInstance(gomock.Any(), teamID).Return(domain.Instance{}, errors.New("unexpected error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repoMock := repomock.NewMockRepository(ctrl)
			useCaseMock := usecasemock.NewMockUseCase(ctrl)
			e := echo.New()
			tt.mockSetup(useCaseMock)

			req := httptest.NewRequest(http.MethodPost, "/teams/"+teamID.String()+"/instances", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("teamID")
			c.SetParamValues(teamID.String())
			h := NewHandler(useCaseMock, repoMock, nil)

			err := h.CreateTeamInstance(c)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			if tt.expectedCode == http.StatusCreated {
				var res openapi.Instance
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
				assert.Equal(t, tt.expectID, uuid.UUID(res.ID))
			}
		})
	}
}

func TestDeleteTeamInstance(t *testing.T) {
	teamID := uuid.New()
	instanceID := uuid.New()
	tests := []struct {
		name         string
		mockSetup    func(*usecasemock.MockUseCase)
		expectedCode int
	}{
		{
			name: "success",
			mockSetup: func(m *usecasemock.MockUseCase) {
				m.EXPECT().GetInstance(gomock.Any(), instanceID).Return(domain.Instance{
					ID:     instanceID,
					TeamID: teamID,
				}, nil)
				m.EXPECT().DeleteInstance(gomock.Any(), instanceID).Return(nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "instance not found",
			mockSetup: func(m *usecasemock.MockUseCase) {
				m.EXPECT().GetInstance(gomock.Any(), instanceID).Return(domain.Instance{}, usecase.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "not found: instance does not belong to  the team",
			mockSetup: func(m *usecasemock.MockUseCase) {
				m.EXPECT().GetInstance(gomock.Any(), instanceID).Return(domain.Instance{
					ID:     instanceID,
					TeamID: uuid.New(), // Different team ID
				}, nil)
				// No call to DeleteInstance expected, as the instance does not belong to the team
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "internal error",
			mockSetup: func(m *usecasemock.MockUseCase) {
				m.EXPECT().GetInstance(gomock.Any(), instanceID).Return(domain.Instance{
					ID:     instanceID,
					TeamID: teamID,
				}, nil)
				m.EXPECT().DeleteInstance(gomock.Any(), instanceID).Return(errors.New("unexpected error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repoMock := repomock.NewMockRepository(ctrl)
			useCaseMock := usecasemock.NewMockUseCase(ctrl)
			e := echo.New()
			tt.mockSetup(useCaseMock)

			req := httptest.NewRequest(http.MethodDelete, "/teams/"+teamID.String()+"/instances/"+instanceID.String(), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("teamID", "instanceID")
			c.SetParamValues(teamID.String(), instanceID.String())
			h := NewHandler(useCaseMock, repoMock, nil)

			err := h.DeleteTeamInstance(c)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestPatchTeamInstance(t *testing.T) {
	teamID := uuid.New()
	instanceID := uuid.New()
	tests := []struct {
		name         string
		body         openapi.PatchTeamInstanceReq
		mockSetup    func(*usecasemock.MockUseCase)
		expectedCode int
	}{
		{
			name: "success: start instance",
			body: openapi.PatchTeamInstanceReq{Operation: openapi.InstanceOperationStart},
			mockSetup: func(m *usecasemock.MockUseCase) {
				m.EXPECT().GetInstance(gomock.Any(), instanceID).Return(domain.Instance{
					ID:     instanceID,
					TeamID: teamID,
				}, nil)
				m.EXPECT().UpdateInstance(gomock.Any(), instanceID, gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "success: stop instance",
			body: openapi.PatchTeamInstanceReq{Operation: openapi.InstanceOperationStop},
			mockSetup: func(m *usecasemock.MockUseCase) {
				m.EXPECT().GetInstance(gomock.Any(), instanceID).Return(domain.Instance{
					ID:     instanceID,
					TeamID: teamID,
				}, nil)
				m.EXPECT().UpdateInstance(gomock.Any(), instanceID, gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "bad request",
			body: openapi.PatchTeamInstanceReq{Operation: openapi.InstanceOperation("invalid")},
			mockSetup: func(_ *usecasemock.MockUseCase) {
				// no call expected
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "instance not found",
			body: openapi.PatchTeamInstanceReq{Operation: openapi.InstanceOperationStart},
			mockSetup: func(m *usecasemock.MockUseCase) {
				m.EXPECT().GetInstance(gomock.Any(), instanceID).Return(domain.Instance{}, usecase.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "instance does not belong to the team",
			body: openapi.PatchTeamInstanceReq{Operation: openapi.InstanceOperationStart},
			mockSetup: func(m *usecasemock.MockUseCase) {
				m.EXPECT().GetInstance(gomock.Any(), instanceID).Return(domain.Instance{
					ID:     instanceID,
					TeamID: uuid.New(), // Different team ID
				}, nil)
				// No call to UpdateInstance expected, as the instance does not belong to the team
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "internal error",
			body: openapi.PatchTeamInstanceReq{Operation: openapi.InstanceOperationStart},
			mockSetup: func(m *usecasemock.MockUseCase) {
				m.EXPECT().GetInstance(gomock.Any(), instanceID).Return(domain.Instance{
					ID:     instanceID,
					TeamID: teamID,
				}, nil)
				m.EXPECT().UpdateInstance(gomock.Any(), instanceID, gomock.Any()).Return(errors.New("unexpected error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repoMock := repomock.NewMockRepository(ctrl)
			useCaseMock := usecasemock.NewMockUseCase(ctrl)
			e := echo.New()
			if tt.mockSetup != nil {
				tt.mockSetup(useCaseMock)
			}

			bodyBytes := lo.Must(tt.body.MarshalJSON())
			req := httptest.NewRequest(http.MethodPatch, "/teams/"+teamID.String()+"/instances/"+instanceID.String(),
				bytes.NewReader(bodyBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("teamID", "instanceID")
			c.SetParamValues(teamID.String(), instanceID.String())
			h := NewHandler(useCaseMock, repoMock, nil)

			err := h.PatchTeamInstance(c)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestGetInstances(t *testing.T) {
	instances := []domain.Instance{
		{
			ID:     uuid.New(),
			TeamID: uuid.New(),
			Index:  1,
			Infra: domain.InfraInstance{
				PublicIP:  ptr.Of("1.2.3.4"),
				PrivateIP: ptr.Of("10.0.0.1"),
				Status:    domain.InstanceStatusRunning,
			},
			CreatedAt: time.Now(),
		},
		{
			ID:     uuid.New(),
			TeamID: uuid.New(),
			Index:  2,
			Infra: domain.InfraInstance{
				PublicIP:  ptr.Of("5.6.7.8"),
				PrivateIP: ptr.Of("10.0.0.2"),
				Status:    domain.InstanceStatusStopped,
			},
			CreatedAt: time.Now(),
		},
	}
	tests := []struct {
		name         string
		mockSetup    func(*usecasemock.MockUseCase)
		expectedCode int
		expectedLen  int
	}{
		{
			name: "success",
			mockSetup: func(m *usecasemock.MockUseCase) {
				m.EXPECT().GetAllInstances(gomock.Any()).Return(instances, nil)
			},
			expectedCode: http.StatusOK,
			expectedLen:  2,
		},
		{
			name: "internal error",
			mockSetup: func(m *usecasemock.MockUseCase) {
				m.EXPECT().GetAllInstances(gomock.Any()).Return(nil, errors.New("unexpected error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedLen:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repoMock := repomock.NewMockRepository(ctrl)
			useCaseMock := usecasemock.NewMockUseCase(ctrl)
			e := echo.New()
			tt.mockSetup(useCaseMock)

			req := httptest.NewRequest(http.MethodGet, "/instances", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := NewHandler(useCaseMock, repoMock, nil)

			err := h.GetInstances(c)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			if tt.expectedCode == http.StatusOK {
				var res []openapi.Instance
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
				assert.Len(t, res, tt.expectedLen)
			}
		})
	}
}
