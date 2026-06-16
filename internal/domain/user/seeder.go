package user

import (
	"go-library/internal/config"
	"log"
)

func SeedAdmin(repo Repository, cfg *config.Config) {
	if cfg.AdminEmail == "" || cfg.AdminPassword == "" {
		log.Println("Super admin credentials not set, skipping seeding")
		return
	}

	existing, err := repo.GetUserByEmail(cfg.AdminEmail)
	if err != nil {
		log.Println("Failed to check super admin existence:", err)
		return
	}

	if existing != nil {
		log.Println("Super admin already exists, skipping seeding")
		return
	}

	name := cfg.AdminName
	if name == "" {
		name = "Super Admin"
	}

	admin := User{
		Name:  name,
		Email: cfg.AdminEmail,
		Role:  RoleSuperAdmin,
	}

	if err := admin.hashPassword(cfg.AdminPassword); err != nil {
		log.Println("Failed to hash super admin password:", err)
		return
	}

	if err := repo.RegisterUser(&admin); err != nil {
		log.Println("Failed to seed super admin:", err)
		return
	}

	log.Println("Super admin seeded successfully")
}
