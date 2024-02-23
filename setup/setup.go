package setup

func Setup() {
	Db()
	Migrate()
	Seed()
	Resources()
	Scanner()
	Encoder()
	Copier()
}
