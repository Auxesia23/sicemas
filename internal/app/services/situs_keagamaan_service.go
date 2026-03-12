package services

import (
	"context"
	"database/sql"
	"errors"
	"mime/multipart"
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/repositories"
	"situs-keagamaan/internal/dto"
	"situs-keagamaan/internal/utils"

	"github.com/casbin/casbin/v2"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type SitusKeagamaanService interface {
	CreateSitusKeagamaan(ctx context.Context, in *dto.SitusKeagamaanRequest, author uuid.UUID) (uuid.UUID, error)
	UploadFotoSitus(ctx context.Context, foto []multipart.File, situsId uuid.UUID) error
	GetAllSitusKeagamaan(ctx context.Context) ([]dto.SitusKeagamaanResponse, error)
	GetDetailSitusKeagamaan(ctx context.Context, id uuid.UUID) (*dto.SitusKeagamaanDetailResponse, error)
}

type situsKeagamaanServiceImpl struct {
	situsRepo     repositories.SitusKeagamaanRepository
	fotoSitusRepo repositories.FotoSitusRepository
	cloudinary    *cloudinary.Cloudinary
	enforcer      *casbin.Enforcer
}

func NewSitusKeagamaanService(
	situsRepo repositories.SitusKeagamaanRepository,
	fotoSitusRepo repositories.FotoSitusRepository,
	cloudinary *cloudinary.Cloudinary,
	enforcer *casbin.Enforcer,
) SitusKeagamaanService {
	return &situsKeagamaanServiceImpl{
		situsRepo:     situsRepo,
		fotoSitusRepo: fotoSitusRepo,
		cloudinary:    cloudinary,
		enforcer:      enforcer,
	}
}

func (s *situsKeagamaanServiceImpl) CreateSitusKeagamaan(ctx context.Context, in *dto.SitusKeagamaanRequest, author uuid.UUID) (uuid.UUID, error) {
	if len(in.Detail) == 0 {
		in.Detail = []byte("{}")
	}
	var err error
	in.NomorBadanHukum, err = utils.Encrypt(in.NomorBadanHukum)
	if err != nil {
		return uuid.Nil, err
	}
	in.NomorAIW, err = utils.Encrypt(in.NomorAIW)
	if err != nil {
		return uuid.Nil, err
	}
	in.NomorSertifikatWakaf, err = utils.Encrypt(in.NomorSertifikatWakaf)
	if err != nil {
		return uuid.Nil, err
	}
	in.NomorTelepon, err = utils.Encrypt(in.NomorTelepon)
	if err != nil {
		return uuid.Nil, err
	}
	in.Email, err = utils.Encrypt(in.Email)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := s.situsRepo.Create(ctx, in, author)
	if err != nil {
		return uuid.Nil, apperror.NewInternal("Terjadi kesalahan.")
	}
	return id, nil
}

func (s *situsKeagamaanServiceImpl) UploadFotoSitus(ctx context.Context, foto []multipart.File, situsId uuid.UUID) error {
	if len(foto) == 0 {
		return nil
	}

	results := make([]dto.FotoSitusPayload, len(foto))

	g, gCtx := errgroup.WithContext(ctx)

	for i, f := range foto {
		index, currentFile := i, f

		g.Go(func() error {
			publicId, _ := uuid.NewV7()
			resp, err := s.cloudinary.Upload.Upload(gCtx, currentFile, uploader.UploadParams{
				PublicID:       publicId.String(),
				Folder:         "situs_keagamaan_assets",
				UniqueFilename: api.Bool(false),
				Overwrite:      api.Bool(true),
			})
			if err != nil {
				return err
			}

			results[index] = dto.FotoSitusPayload{
				SitusID:  situsId,
				ImageURL: resp.SecureURL,
				PublicID: resp.PublicID,
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return apperror.NewInternal("Terjadi kesalahan")
	}

	err := s.fotoSitusRepo.BulkCreate(ctx, results)
	if err != nil {
		return apperror.NewInternal("Terjadi kesalahan")
	}

	return nil
}

func (s *situsKeagamaanServiceImpl) GetAllSitusKeagamaan(ctx context.Context) ([]dto.SitusKeagamaanResponse, error) {
	situs, err := s.situsRepo.ReadAll(ctx)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}
	if len(situs) == 0 {
		return []dto.SitusKeagamaanResponse{}, nil
	}
	return situs, nil
}

func (s *situsKeagamaanServiceImpl) GetDetailSitusKeagamaan(ctx context.Context, id uuid.UUID) (*dto.SitusKeagamaanDetailResponse, error) {
	situs, err := s.situsRepo.ReadDetail(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("Situs tidak ditemukan")
		}
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}

	situs.NomorTelepon, err = utils.Decrypt(situs.NomorTelepon)
	if err != nil {
		return nil, err
	}
	situs.Email, err = utils.Decrypt(situs.Email)
	if err != nil {
		return nil, err
	}
	situs.NomorBadanHukum, err = utils.Decrypt(situs.NomorBadanHukum)
	if err != nil {
		return nil, err
	}
	situs.NomorAIW, err = utils.Decrypt(situs.NomorAIW)
	if err != nil {
		return nil, err
	}
	situs.NomorSertifikatWakaf, err = utils.Decrypt(situs.NomorSertifikatWakaf)
	if err != nil {
		return nil, err
	}

	return situs, nil
}
