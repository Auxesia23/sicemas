package database

import (
	"context"
	"log"
	"sicemas/internal/utils"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type seeder struct {
	db *sqlx.DB
	e  *casbin.Enforcer
}

type userSeed struct {
	NIP          string
	NamaLengkap  string
	Jabatan      string
	UnitKerja    string
	Email        string
	NomorTelepon string
}

func NewSeeder(db *sqlx.DB, e *casbin.Enforcer) *seeder {
	return &seeder{
		db,
		e,
	}
}

func (s *seeder) UserSeeder() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	users := []*userSeed{
		{
			NIP:          "200308232022111001",
			NamaLengkap:  "Fajar Gustiana",
			Jabatan:      "Developer",
			UnitKerja:    "Universitas Komputer Indonesia",
			Email:        "fajargustiana2@gmail.com",
			NomorTelepon: "085720928785",
		},
	}

	query := `
		INSERT INTO users (nip_index,password_hash,nama_lengkap,nip,jabatan,unit_kerja,email,nomor_telepon)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id;
	`
	for _, u := range users {
		var id uuid.UUID

		nipIndex := utils.HashIndex(u.NIP)
		passwordHash, err := utils.HashPassword(u.NomorTelepon)
		if err != nil {
			return err
		}
		encryptedNIP, err := utils.Encrypt(u.NIP)
		if err != nil {
			return err
		}
		encryptedEmail, err := utils.Encrypt(u.Email)
		if err != nil {
			return err
		}
		encryptedTelepon, err := utils.Encrypt(u.NomorTelepon)
		if err != nil {
			return err
		}

		err = s.db.GetContext(ctx, &id, query,
			nipIndex,
			passwordHash,
			u.NamaLengkap,
			encryptedNIP,
			u.Jabatan,
			u.UnitKerja,
			encryptedEmail,
			encryptedTelepon,
		)

		if err != nil {
			pqErr := err.(*pq.Error)
			if pqErr.Code == "23505" {
				log.Println("User already exist, seed skipped!")
				return nil
			}
			return err
		}
		s.e.AddGroupingPolicy(id.String(), "admin")

		log.Printf("User %v added succesfuly!", strings.ToUpper(u.NamaLengkap))
	}
	return nil
}

func (s *seeder) PolicySeeder() {
	s.e.AddPolicy("admin", "user", "create")
	s.e.AddPolicy("admin", "user", "read")
	s.e.AddPolicy("admin", "user", "update")
	s.e.AddPolicy("admin", "user", "delete")

	s.e.AddPolicy("admin", "policy", "create")
	s.e.AddPolicy("admin", "policy", "read")
	s.e.AddPolicy("admin", "policy", "update")
	s.e.AddPolicy("admin", "policy", "delete")

	s.e.AddPolicy("admin", "role", "create")
	s.e.AddPolicy("admin", "role", "read")
	s.e.AddPolicy("admin", "role", "update")
	s.e.AddPolicy("admin", "role", "delete")

	s.e.AddPolicy("operator", "situs", "create")
	s.e.AddPolicy("operator", "situs", "read_all")
	s.e.AddPolicy("operator", "situs", "update_all")
	s.e.AddPolicy("operator", "situs", "delete")

	s.e.AddPolicy("penyuluh", "situs", "create")
	s.e.AddPolicy("penyuluh", "situs", "read_own")
	s.e.AddPolicy("penyuluh", "situs", "update_own")
}

func (s *seeder) RoleSeeder() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	roles := []string{
		"admin", "operator", "penyuluh",
	}

	query := `
		INSERT INTO roles (name)
		VALUES($1);
	`
	for _, r := range roles {
		_, err := s.db.ExecContext(ctx, query, r)
		if err != nil {
			pqErr := err.(*pq.Error)
			if pqErr.Code == "23505" {
				log.Println("Role already exist, seed skipped!")
				return nil
			}
			return err
		}
	}
	return nil
}
