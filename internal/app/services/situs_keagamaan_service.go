package services

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"mime/multipart"
	apperror "sicemas/internal/app/appError"
	"sicemas/internal/app/repositories"
	"sicemas/internal/dto"
	"sicemas/internal/entity"
	"sicemas/internal/utils"

	"github.com/casbin/casbin/v2"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"golang.org/x/sync/errgroup"
)

type SitusKeagamaanService interface {
	CreateSitusKeagamaan(ctx context.Context, in *dto.SitusKeagamaanRequest, author, actorId uuid.UUID) (uuid.UUID, error)
	UploadFotoSitus(ctx context.Context, foto []multipart.File, situsId uuid.UUID) error
	DeleteFotoSitus(ctx context.Context, situsId uuid.UUID, foto *dto.DeleteFoto) error
	GetAllSitusKeagamaan(ctx context.Context, userId uuid.UUID) ([]dto.SitusKeagamaanResponse, error)
	GetDetailSitusKeagamaan(ctx context.Context, situsId, userId uuid.UUID) (*dto.SitusKeagamaanDetailResponse, error)
	ExportSitusToExcel(ctx context.Context, jenisSitus string, actorId uuid.UUID) (*excelize.File, error)
	UpdateSitus(ctx context.Context, id, userId uuid.UUID, in *dto.SitusKeagamaanUpdate, actorId uuid.UUID) error
	DeleteSitus(ctx context.Context, id uuid.UUID, actorId uuid.UUID) error
	VerifySitus(ctx context.Context, situsId uuid.UUID, in *dto.VerifikasiSitusRequest, actorId uuid.UUID) error
}

type situsKeagamaanServiceImpl struct {
	situsRepo     repositories.SitusKeagamaanRepository
	fotoSitusRepo repositories.FotoSitusRepository
	activityRepo  repositories.ActivityRepository
	cloudinary    *cloudinary.Cloudinary
	enforcer      *casbin.Enforcer
	logger        *slog.Logger
}

func NewSitusKeagamaanService(
	situsRepo repositories.SitusKeagamaanRepository,
	fotoSitusRepo repositories.FotoSitusRepository,
	activityRepo repositories.ActivityRepository,
	cloudinary *cloudinary.Cloudinary,
	enforcer *casbin.Enforcer,
	logger *slog.Logger,
) SitusKeagamaanService {
	return &situsKeagamaanServiceImpl{
		situsRepo:     situsRepo,
		fotoSitusRepo: fotoSitusRepo,
		activityRepo:  activityRepo,
		cloudinary:    cloudinary,
		enforcer:      enforcer,
		logger:        logger,
	}
}

func (s *situsKeagamaanServiceImpl) CreateSitusKeagamaan(ctx context.Context, in *dto.SitusKeagamaanRequest, author, actorId uuid.UUID) (uuid.UUID, error) {
	s.logger.Info("creating new situs keagamaan", "nama", in.Nama, "author_id", author, "actor_id", actorId)

	if len(in.Detail) == 0 {
		in.Detail = []byte("{}")
	}
	var err error
	in.NomorBadanHukum, err = utils.Encrypt(in.NomorBadanHukum)
	if err != nil {
		s.logger.Error("failed to encrypt NomorBadanHukum", "error", err)
		return uuid.Nil, err
	}
	in.NomorAIW, err = utils.Encrypt(in.NomorAIW)
	if err != nil {
		s.logger.Error("failed to encrypt NomorAIW", "error", err)
		return uuid.Nil, err
	}
	in.NomorSertifikatWakaf, err = utils.Encrypt(in.NomorSertifikatWakaf)
	if err != nil {
		s.logger.Error("failed to encrypt NomorSertifikatWakaf", "error", err)
		return uuid.Nil, err
	}
	in.NomorTelepon, err = utils.Encrypt(in.NomorTelepon)
	if err != nil {
		s.logger.Error("failed to encrypt NomorTelepon", "error", err)
		return uuid.Nil, err
	}
	in.Email, err = utils.Encrypt(in.Email)
	if err != nil {
		s.logger.Error("failed to encrypt Email", "error", err)
		return uuid.Nil, err
	}
	for i, nomor := range in.NomorTelponPengurus {
		encryptedNomor, err := utils.Encrypt(nomor)
		if err != nil {
			s.logger.Error("failed to encrypt phone number", "index", i, "error", err)
			return uuid.Nil, err
		}
		in.NomorTelponPengurus[i] = encryptedNomor
	}

	id, err := s.situsRepo.Create(ctx, in, author)
	if err != nil {
		s.logger.Error("failed to create situs keagamaan", "error", err)
		return uuid.Nil, apperror.NewInternal("Terjadi kesalahan.")
	}

	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Menambahkan situs",
		EntityType: "SITUS",
		EntityID:   id.String(),
		TargetName: in.Nama,
	})

	s.logger.Info("situs keagamaan created successfully", "id", id, "nama", in.Nama)
	return id, nil
}

func (s *situsKeagamaanServiceImpl) UploadFotoSitus(ctx context.Context, foto []multipart.File, situsId uuid.UUID) error {
	s.logger.Info("uploading photos for situs", "situs_id", situsId, "photo_count", len(foto))

	if len(foto) == 0 {
		s.logger.Warn("no photos provided for upload")
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
				s.logger.Error("failed to upload photo to Cloudinary", "situs_id", situsId, "error", err)
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
		s.logger.Error("failed during photo upload process", "situs_id", situsId, "error", err)
		return apperror.NewInternal("Terjadi kesalahan")
	}

	err := s.fotoSitusRepo.BulkCreate(ctx, results)
	if err != nil {
		s.logger.Error("failed to save photo records to database", "situs_id", situsId, "error", err)
		return apperror.NewInternal("Terjadi kesalahan")
	}

	s.logger.Info("photos uploaded successfully", "situs_id", situsId, "uploaded_count", len(results))
	return nil
}

func (s *situsKeagamaanServiceImpl) DeleteFotoSitus(ctx context.Context, situsId uuid.UUID, foto *dto.DeleteFoto) error {
	s.logger.Info("deleting situs photos", "situs_id", situsId, "photo_count", len(foto.IDs))

	if len(foto.IDs) == 0 {
		s.logger.Warn("no photo IDs provided for deletion")
		return nil
	}

	publicIDs, err := s.fotoSitusRepo.GetPublicIDs(ctx, foto.IDs)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("photo not found for deletion", "situs_id", situsId)
			return apperror.NewNotFound("Foto situs tidak ditemukan")
		}
		s.logger.Error("failed to get photo public IDs", "situs_id", situsId, "error", err)
		return apperror.NewInternal("Terjadi kesalahan")
	}

	if err := s.fotoSitusRepo.Delete(ctx, foto.IDs, situsId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("photo not found during deletion", "situs_id", situsId)
			return apperror.NewNotFound("Foto situs tidak ditemukan")
		}
		s.logger.Error("failed to delete photo records", "situs_id", situsId, "error", err)
		return apperror.NewInternal("Terjadi kesalahan")
	}

	_, err = s.cloudinary.Admin.DeleteAssets(ctx, admin.DeleteAssetsParams{PublicIDs: publicIDs})
	if err != nil {
		s.logger.Error("failed to delete photos from Cloudinary", "situs_id", situsId, "error", err)
		return apperror.NewInternal("Terjadi kesalahan")
	}

	s.logger.Info("situs photos deleted successfully", "situs_id", situsId)
	return nil
}

func (s *situsKeagamaanServiceImpl) GetAllSitusKeagamaan(ctx context.Context, userId uuid.UUID) ([]dto.SitusKeagamaanResponse, error) {
	s.logger.Info("fetching all situs keagamaan", "user_id", userId)

	canReadAll, err := s.enforcer.Enforce(userId.String(), "situs", "read_all")
	if err != nil {
		s.logger.Error("failed to validate access permission", "user_id", userId, "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan validasi akses.")
	}

	var situs []dto.SitusKeagamaanResponse

	if canReadAll {
		situs, err = s.situsRepo.ReadAll(ctx)
		if err != nil {
			s.logger.Error("failed to read all situs", "error", err)
			return nil, apperror.NewInternal("Terjadi kesalahan validasi akses.")
		}
	} else {
		situs, err = s.situsRepo.ReadOwn(ctx, userId)
		if err != nil {
			s.logger.Error("failed to read user's own situs", "user_id", userId, "error", err)
			return nil, apperror.NewInternal("Terjadi kesalahan validasi akses.")
		}
	}

	if len(situs) == 0 {
		s.logger.Info("no situs found")
		return []dto.SitusKeagamaanResponse{}, nil
	}

	s.logger.Info("fetched situs successfully", "count", len(situs))
	return situs, nil
}

func (s *situsKeagamaanServiceImpl) GetDetailSitusKeagamaan(ctx context.Context, situsId, userId uuid.UUID) (*dto.SitusKeagamaanDetailResponse, error) {
	s.logger.Info("fetching detail situs keagamaan", "situs_id", situsId, "user_id", userId)

	canReadAll, err := s.enforcer.Enforce(userId.String(), "situs", "read_all")
	if err != nil {
		s.logger.Error("failed to validate access permission", "user_id", userId, "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan validasi akses.")
	}

	if !canReadAll {
		err := s.situsRepo.CheckOwnership(ctx, situsId, userId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				s.logger.Warn("user does not have access to view this situs", "situs_id", situsId, "user_id", userId)
				return nil, apperror.NewForbidden("Anda tidak memiliki hak akses untuk melihat situs ini")
			}
			s.logger.Error("failed to check ownership", "situs_id", situsId, "error", err)
			return nil, apperror.NewInternal("Terjadi kesalahan.")
		}
	}

	situs, err := s.situsRepo.ReadDetail(ctx, situsId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("situs not found for detail", "situs_id", situsId)
			return nil, apperror.NewNotFound("Situs tidak ditemukan")
		}
		s.logger.Error("failed to read situs detail", "situs_id", situsId, "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan.")
	}

	situs.NomorTelepon, err = utils.Decrypt(situs.NomorTelepon)
	if err != nil {
		s.logger.Error("failed to decrypt phone number", "situs_id", situsId, "error", err)
		return nil, err
	}
	situs.Email, err = utils.Decrypt(situs.Email)
	if err != nil {
		s.logger.Error("failed to decrypt email", "situs_id", situsId, "error", err)
		return nil, err
	}
	situs.NomorBadanHukum, err = utils.Decrypt(situs.NomorBadanHukum)
	if err != nil {
		s.logger.Error("failed to decrypt NomorBadanHukum", "situs_id", situsId, "error", err)
		return nil, err
	}
	situs.NomorAIW, err = utils.Decrypt(situs.NomorAIW)
	if err != nil {
		s.logger.Error("failed to decrypt NomorAIW", "situs_id", situsId, "error", err)
		return nil, err
	}
	situs.NomorSertifikatWakaf, err = utils.Decrypt(situs.NomorSertifikatWakaf)
	if err != nil {
		s.logger.Error("failed to decrypt NomorSertifikatWakaf", "situs_id", situsId, "error", err)
		return nil, err
	}
	for i, nomor := range situs.NomorTelponPengurus {
		decryptedNomor, err := utils.Decrypt(nomor)
		if err != nil {
			s.logger.Error("failed to decrypt manager phone number", "situs_id", situsId, "index", i, "error", err)
			return nil, err
		}
		situs.NomorTelponPengurus[i] = decryptedNomor
	}

	fotoSitus, err := s.fotoSitusRepo.GetBySitusID(ctx, situsId)
	if err != nil {
		s.logger.Error("failed to get photos for situs", "situs_id", situsId, "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan")
	}
	if len(*fotoSitus) == 0 {
		situs.Foto = &[]dto.FotoResponse{}
	} else {
		situs.Foto = fotoSitus
	}

	s.logger.Info("situs detail fetched successfully", "situs_id", situsId)
	return situs, nil
}

func (s *situsKeagamaanServiceImpl) ExportSitusToExcel(ctx context.Context, jenisSitus string, actorId uuid.UUID) (*excelize.File, error) {
	s.logger.Info("exporting situs to Excel", "jenis_situs", jenisSitus, "user_id", actorId)

	situsList, err := s.situsRepo.ReadAllDetail(ctx, jenisSitus)
	if err != nil {
		s.logger.Error("failed to fetch situs data for export", "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan saat mengambil data.")
	}

	for i := range situsList {
		situsList[i].NomorTelepon, err = utils.Decrypt(situsList[i].NomorTelepon)
		if err != nil {
			s.logger.Error("failed to decrypt phone number", "situs_id", situsList[i].ID, "error", err)
			return nil, err
		}

		situsList[i].Email, err = utils.Decrypt(situsList[i].Email)
		if err != nil {
			s.logger.Error("failed to decrypt email", "situs_id", situsList[i].ID, "error", err)
			return nil, err
		}

		situsList[i].NomorBadanHukum, err = utils.Decrypt(situsList[i].NomorBadanHukum)
		if err != nil {
			s.logger.Error("failed to decrypt NomorBadanHukum", "situs_id", situsList[i].ID, "error", err)
			return nil, err
		}

		situsList[i].NomorAIW, err = utils.Decrypt(situsList[i].NomorAIW)
		if err != nil {
			s.logger.Error("failed to decrypt NomorAIW", "situs_id", situsList[i].ID, "error", err)
			return nil, err
		}

		situsList[i].NomorSertifikatWakaf, err = utils.Decrypt(situsList[i].NomorSertifikatWakaf)
		if err != nil {
			s.logger.Error("failed to decrypt NomorSertifikatWakaf", "situs_id", situsList[i].ID, "error", err)
			return nil, err
		}

		for j, nomor := range situsList[i].NomorTelponPengurus {
			decryptedNomor, err := utils.Decrypt(nomor)
			if err != nil {
				s.logger.Error("failed to decrypt manager phone number", "situs_id", situsList[i].ID, "index", j, "error", err)
				return nil, err
			}
			situsList[i].NomorTelponPengurus[j] = decryptedNomor
		}
	}

	excelFile, err := utils.ExportSitusToExcel(situsList)
	if err != nil {
		s.logger.Error("failed to generate Excel file", "error", err)
		return nil, apperror.NewInternal("Terjadi kesalahan saat men-generate Excel.")
	}

	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Mengekspor data situs",
		EntityType: "SITUS",
		EntityID:   uuid.Nil.String(),
		TargetName: jenisSitus,
	})

	s.logger.Info("situs Excel file generated successfully", "row_count", len(situsList))
	return excelFile, nil
}

func (s *situsKeagamaanServiceImpl) UpdateSitus(ctx context.Context, id, userId uuid.UUID, in *dto.SitusKeagamaanUpdate, actorId uuid.UUID) error {
	s.logger.Info("updating situs", "situs_id", id, "nama", in.Nama, "user_id", userId, "actor_id", actorId)

	if err := s.situsRepo.CheckOwnership(ctx, id, userId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("user does not have permission to update this situs", "situs_id", id, "user_id", userId)
			return apperror.NewForbidden("Anda tidak memiliki izin untuk mengubah situs ini")
		}
		s.logger.Error("failed to check ownership", "situs_id", id, "error", err)
		return apperror.NewInternal("Terjadi kesahalahan.")
	}

	if len(in.Detail) == 0 {
		in.Detail = []byte("{}")
	}

	var err error
	in.NomorBadanHukum, err = utils.Encrypt(in.NomorBadanHukum)
	if err != nil {
		s.logger.Error("failed to encrypt NomorBadanHukum", "error", err)
		return err
	}

	in.NomorAIW, err = utils.Encrypt(in.NomorAIW)
	if err != nil {
		s.logger.Error("failed to encrypt NomorAIW", "error", err)
		return err
	}

	in.NomorSertifikatWakaf, err = utils.Encrypt(in.NomorSertifikatWakaf)
	if err != nil {
		s.logger.Error("failed to encrypt NomorSertifikatWakaf", "error", err)
		return err
	}

	in.NomorTelepon, err = utils.Encrypt(in.NomorTelepon)
	if err != nil {
		s.logger.Error("failed to encrypt phone number", "error", err)
		return err
	}

	in.Email, err = utils.Encrypt(in.Email)
	if err != nil {
		s.logger.Error("failed to encrypt email", "error", err)
		return err
	}
	for i, nomor := range in.NomorTelponPengurus {
		encryptedNomor, err := utils.Encrypt(nomor)
		if err != nil {
			s.logger.Error("failed to encrypt manager phone number", "index", i, "error", err)
			return err
		}
		in.NomorTelponPengurus[i] = encryptedNomor
	}

	err = s.situsRepo.Update(ctx, id, in)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("situs not found for update", "situs_id", id)
			return apperror.NewNotFound("Situs tidak ditemukan")
		}
		s.logger.Error("failed to update situs", "error", err)
		return apperror.NewInternal("Terjadi kesalahan saat memperbarui data situs.")
	}

	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Memperbarui data situs",
		EntityType: "SITUS",
		EntityID:   id.String(),
		TargetName: in.Nama,
	})

	s.logger.Info("situs updated successfully", "situs_id", id)
	return nil
}

func (s *situsKeagamaanServiceImpl) DeleteSitus(ctx context.Context, id, actorId uuid.UUID) error {
	s.logger.Info("deleting situs", "situs_id", id, "actor_id", actorId)

	situsTarget, err := s.situsRepo.ReadDetail(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("situs not found for deletion", "situs_id", id)
			return apperror.NewNotFound("Situs tidak ditemukan")
		}
		s.logger.Error("failed to read situs for deletion", "error", err)
		return apperror.NewInternal("Terjadi kesalahan.")
	}

	if err := s.situsRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("situs not found during deletion", "situs_id", id)
			return apperror.NewNotFound("Situs tidak ditemukan")
		}
		s.logger.Error("failed to delete situs", "error", err)
		return apperror.NewInternal("Terjadi kesalahan.")
	}

	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Menghapus situs",
		EntityType: "SITUS",
		EntityID:   id.String(),
		TargetName: situsTarget.Nama,
	})

	s.logger.Info("situs deleted successfully", "situs_id", id)
	return nil
}

func (s *situsKeagamaanServiceImpl) VerifySitus(ctx context.Context, id uuid.UUID, in *dto.VerifikasiSitusRequest, actorId uuid.UUID) error {
	s.logger.Info("verifying situs", "situs_id", id, "actor_id", actorId)

	situsTarget, err := s.situsRepo.ReadDetail(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("situs not found for verification", "situs_id", id)
			return apperror.NewNotFound("Situs tidak ditemukan")
		}
		s.logger.Error("failed to read situs for verification", "error", err)
		return apperror.NewInternal("Terjadi kesalahan.")
	}

	if err := s.situsRepo.Verify(ctx, id, in); err != nil {
		s.logger.Error("failed to verify situs", "situs_id", id, "error", err)
		return apperror.NewInternal("Terjadi kesalahan.")
	}

	_ = s.activityRepo.InsertActivity(ctx, &entity.Activity{
		UserID:     actorId,
		ActionType: "Memverifikasi situs",
		EntityType: "SITUS",
		EntityID:   id.String(),
		TargetName: situsTarget.Nama,
	})

	s.logger.Info("situs verified successfully", "situs_id", id)
	return nil
}
