//go:build wireinject
// +build wireinject

package admin

import (
	handler "e-course-management/internal/admin/delivery/http"
	repository "e-course-management/internal/admin/repository"
	usecase "e-course-management/internal/admin/usecase"

	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitializedService(db *gorm.DB) *handler.AdminHandler {
	wire.Build(
		repository.NewAdminRepository,
		usecase.NewAdminUseCase,
		handler.NewAdminHandler,
	)

	return &handler.AdminHandler{}
}
