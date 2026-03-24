package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"mime/multipart"
	apperror "situs-keagamaan/internal/app/appError"
	"situs-keagamaan/internal/app/repositories"
	"situs-keagamaan/internal/dto"
	"situs-keagamaan/internal/entity"
	"situs-keagamaan/internal/utils"

	"github.com/casbin/casbin/v2"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type SitusKeagamaanService interface {
	CreateSitusKeagamaan(ctx context.Context, in *dto.SitusKeagamaanRequest, author, actorId uuid.UUID) (uuid.UUID, error)
	UploadFotoSitus(ctx context.Context, foto []multipart.File, situsId uuid.UUID) error
	DeleteFotoSitus(ctx context.Context, situsId uuid.UUID, foto *dto.DeleteFoto) error
	GetAllSitusKeagamaan(ctx context.Context, userId uuid.UUID) ([]dto.SitusKeagamaanResponse, error)
	GetDetailSitusKeagamaan(ctx context.Context, situsId, userId uuid.UUID) (*dto.SitusKeagamaanDetailResponse, error)
	UpdateSitus(ctx context.Context, id, userId uuid.UUID, in *dto.SitusKeagamaanUpdate, actorId uuid.UUID) error
	DeleteSitus(ctx context.Context, id uuid.UUID, actorId uuid.UUID) error
	VerifySitus(ctx context.Context, situsId uuid.UUID, in *dto.VerifikasiSitusRequest, actorId uuid.UUID) error

	GetAllSitusForPublic(ctx context.Context, filter dto.PublicListFilter) ([]dto.SitusPublicListResponse, error)
	GetSitusDetailForPublic(ctx context.Context, id uuid.UUID) (*dto.SitusPublicDetailResponse, error)
	GetLandingStats(ctx context.Context) (*dto.LandingStatsResponse, error)
}

type situsKeagamaanServiceImpl struct {
	situsRepo     repositories.SitusKeagamaanRepository
	fotoSitusRepo repositories.FotoSitusRepository
	activityRepo  repositories.ActivityRepository
	cloudinary    *cloudinary.Cloudinary
	enforcer      *casbin.Enforcer
}

func NewSitusKeagamaanService(
	situsRepo repositories.SitusKeagamaanRepository,
	fotoSitusRepo repositories.FotoSitusRepository,
	activityRepo repositories.ActivityRepository,
	cloudinary *cloudinary.Cloudinary,
	enforcer *casbin.Enforcer,
) SitusKeagamaanService {
	return &situsKeagamaanServiceImpl{
		situsRepo:     situsRepo,
		fotoSitusRepo: fotoSitusRepo,
		activityRepo:  activityRepo,
		cloudinary:    cloudinary,
		enforcer:      enforcer,
	}
}

func (s *situsKeagamaanServiceImpl) CreateSitusKeagamaan(ctx context.Context, in *dto.SitusKeagamaanRequest, author, actorId uuid.UUID) (uuid.UUID, error) {
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
	for i, nomor := range in.NomorTelponPengurus {
		encryptedNomor, err := utils.Encrypt(nomor)
		if err != nil {
			return uuid.Nil, err
		}
		in.NomorTelponPengurus[i] = encryptedNomor
	}

	id, err := s.situsRepo.Create(ctx, in, author)
	if err != nil {
		return uuid.Nil, apperror.NewInternal("Terjadi kesalahan.")
	}
	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Menambahkan situs",
		EntityType: "SITUS",
		EntityID:   id.String(),
		TargetName: in.Nama,
	})
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

func (s *situsKeagamaanServiceImpl) DeleteFotoSitus(ctx context.Context, situsId uuid.UUID, foto *dto.DeleteFoto) error {
	if len(foto.IDs) == 0 {
		return nil
	}
	publicIDs, err := s.fotoSitusRepo.GetPublicIDs(ctx, foto.IDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFound("Foto situs tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi kesalahan")
	}

	if err := s.fotoSitusRepo.Delete(ctx, foto.IDs, situsId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFound("Foto situs tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi kesalahan")
	}

	_, err = s.cloudinary.Admin.DeleteAssets(ctx, admin.DeleteAssetsParams{PublicIDs: publicIDs})
	if err != nil {
		return apperror.NewInternal("Terjadi kesalahan")
	}

	return nil
}

func (s *situsKeagamaanServiceImpl) GetAllSitusKeagamaan(ctx context.Context, userId uuid.UUID) ([]dto.SitusKeagamaanResponse, error) {
	canReadAll, err := s.enforcer.Enforce(userId.String(), "situs", "read_all")
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan validasi akses.")
	}

	var situs []dto.SitusKeagamaanResponse

	if canReadAll {
		situs, err = s.situsRepo.ReadAll(ctx)
		if err != nil {
			return nil, apperror.NewInternal("Terjadi kesalahan validasi akses.")
		}
	} else {
		situs, err = s.situsRepo.ReadOwn(ctx, userId)
		if err != nil {
			return nil, apperror.NewInternal("Terjadi kesalahan validasi akses.")
		}
	}

	if len(situs) == 0 {
		return []dto.SitusKeagamaanResponse{}, nil
	}

	return situs, nil
}

func (s *situsKeagamaanServiceImpl) GetDetailSitusKeagamaan(ctx context.Context, situsId, userId uuid.UUID) (*dto.SitusKeagamaanDetailResponse, error) {
	canReadAll, err := s.enforcer.Enforce(userId.String(), "situs", "read_all")
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan validasi akses.")
	}

	if !canReadAll {
		err := s.situsRepo.CheckOwnership(ctx, situsId, userId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, apperror.NewForbidden("Anda tidak memiliki hak akses untuk melihat situs ini")
			}
			return nil, apperror.NewInternal("Terjadi kesalahan.")
		}
	}

	situs, err := s.situsRepo.ReadDetail(ctx, situsId)
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
	for i, nomor := range situs.NomorTelponPengurus {
		decryptedNomor, err := utils.Decrypt(nomor)
		if err != nil {
			return nil, err
		}
		situs.NomorTelponPengurus[i] = decryptedNomor
	}

	fotoSitus, err := s.fotoSitusRepo.GetBySitusID(ctx, situsId)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan")
	}
	if len(*fotoSitus) == 0 {
		situs.Foto = &[]dto.FotoResponse{}
	} else {
		situs.Foto = fotoSitus
	}

	return situs, nil
}

func (s *situsKeagamaanServiceImpl) UpdateSitus(ctx context.Context, id, userId uuid.UUID, in *dto.SitusKeagamaanUpdate, actorId uuid.UUID) error {
	if err := s.situsRepo.CheckOwnership(ctx, id, userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewForbidden("Anda tidak memiliki izin untuk mengubah situs ini")
		}
		return apperror.NewInternal("Terjadi kesahalahan.")
	}
	if len(in.Detail) == 0 {
		in.Detail = []byte("{}")
	}

	var err error
	in.NomorBadanHukum, err = utils.Encrypt(in.NomorBadanHukum)
	if err != nil {
		return err
	}

	in.NomorAIW, err = utils.Encrypt(in.NomorAIW)
	if err != nil {
		return err
	}

	in.NomorSertifikatWakaf, err = utils.Encrypt(in.NomorSertifikatWakaf)
	if err != nil {
		return err
	}

	in.NomorTelepon, err = utils.Encrypt(in.NomorTelepon)
	if err != nil {
		return err
	}

	in.Email, err = utils.Encrypt(in.Email)
	if err != nil {
		return err
	}
	for i, nomor := range in.NomorTelponPengurus {
		encryptedNomor, err := utils.Encrypt(nomor)
		if err != nil {
			return err
		}
		in.NomorTelponPengurus[i] = encryptedNomor
	}

	err = s.situsRepo.Update(ctx, id, in)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFound("Situs tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi kesalahan saat memperbarui data situs.")
	}
	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Memperbarui data situs",
		EntityType: "SITUS",
		EntityID:   id.String(),
		TargetName: in.Nama,
	})

	return nil
}

func (s *situsKeagamaanServiceImpl) DeleteSitus(ctx context.Context, id, actorId uuid.UUID) error {
	situsTarget, err := s.situsRepo.ReadDetail(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFound("Situs tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi kesalahan.")
	}
	if err := s.situsRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFound("Situs tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi kesalahan.")
	}
	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Menghapus situs",
		EntityType: "SITUS",
		EntityID:   id.String(),
		TargetName: situsTarget.Nama,
	})
	return nil
}

func (s *situsKeagamaanServiceImpl) VerifySitus(ctx context.Context, id uuid.UUID, in *dto.VerifikasiSitusRequest, actorId uuid.UUID) error {
	situsTarget, err := s.situsRepo.ReadDetail(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NewNotFound("Situs tidak ditemukan")
		}
		return apperror.NewInternal("Terjadi kesalahan.")
	}
	if err := s.situsRepo.Verify(ctx, id, in); err != nil {
		return apperror.NewInternal("Terjadi kesalahan.")
	}
	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Memverifikasi situs",
		EntityType: "SITUS",
		EntityID:   id.String(),
		TargetName: situsTarget.Nama,
	})
	return nil
}

func (s *situsKeagamaanServiceImpl) GetAllSitusForPublic(ctx context.Context, filter dto.PublicListFilter) ([]dto.SitusPublicListResponse, error) {
	situs, err := s.situsRepo.ReadForPublic(ctx, filter)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}
	if len(situs) == 0 {
		return []dto.SitusPublicListResponse{}, nil
	}
	return situs, nil
}

func (s *situsKeagamaanServiceImpl) GetSitusDetailForPublic(ctx context.Context, id uuid.UUID) (*dto.SitusPublicDetailResponse, error) {
	situs, err := s.situsRepo.ReadDetailForPublic(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NewNotFound("Situs tidak ditemukan")
		}
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}

	situs.Fasilitas = make(map[string]any)
	if len(situs.FasilitasJSON) > 0 {
		_ = json.Unmarshal(situs.FasilitasJSON, &situs.Fasilitas)
	}

	situs.Galeri = []string{}
	if len(situs.GaleriJSON) > 0 {
		_ = json.Unmarshal(situs.GaleriJSON, &situs.Galeri)
	}

	return situs, nil
}

func (s *situsKeagamaanServiceImpl) GetLandingStats(ctx context.Context) (*dto.LandingStatsResponse, error) {
	stats, err := s.situsRepo.GetLandingStats(ctx)
	if err != nil {
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}
	return stats, nil
}
