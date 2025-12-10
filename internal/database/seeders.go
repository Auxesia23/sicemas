package database

import (
	"context"
	"log"
	"situs-keagamaan/internal/utils"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type seeder struct {
	db *sqlx.DB
}

type userSeed struct {
	NIP          string
	NIPIndex     []byte
	NamaLengkap  string
	Peran        string
	Email        string
	NomorTelepon string
}

func newSeeder(db *sqlx.DB) *seeder {
	return &seeder{
		db,
	}
}

func (s *seeder) UserSeeder() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	users := []*userSeed{
		{
			NIP:          "200308232022111001",
			NIPIndex:     utils.HashIndex("200308232022111001"),
			NamaLengkap:  "Fajar Gustiana",
			Peran:        "admin",
			Email:        "fajargustiana2@gmail.com",
			NomorTelepon: "085720928785",
		},
	}

	query := `
		INSERT INTO users (nip_index,nama_lengkap,peran,nip,email,nomor_telepon)
		VALUES($1,$2,$3,$4,$5,$6);
	`
	for _, u := range users {
		_, err := s.db.ExecContext(ctx, query,
			u.NIPIndex,
			u.NamaLengkap,
			u.Peran,
			utils.Encrypt(u.NIP),
			utils.Encrypt(u.Email),
			utils.Encrypt(u.NomorTelepon),
		)

		if err != nil {
			pqErr := err.(*pq.Error)
			if pqErr.Code == "23505" {
				log.Println("User already exist, seed skipped!")
				return nil
			}
			return err
		}
		log.Printf("User %v added succesfuly!", strings.ToUpper(u.NamaLengkap))
	}
	return nil
}
