package auth

import "github.com/casbin/casbin/v2"

type seeder struct {
	e *casbin.Enforcer
}

func newSeeder(e *casbin.Enforcer) *seeder {
	return &seeder{
		e,
	}
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

	s.e.AddPolicy("operator", "situs", "create")
	s.e.AddPolicy("operator", "situs", "read")
	s.e.AddPolicy("operator", "situs", "update")
	s.e.AddPolicy("operator", "situs", "delete")

	s.e.AddPolicy("penyuluh", "situs", "create")
	s.e.AddPolicy("penyuluh", "situs", "read")
	s.e.AddPolicy("penyuluh", "situs", "update")
}
