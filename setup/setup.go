package setup

func Setup() {
	Db()
	Migrate()
	Seed()
	Resources()
}
