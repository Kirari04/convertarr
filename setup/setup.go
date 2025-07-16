package setup

// this function runs all setup functions
// APP => initializes app => generates random string for jwt secret
// DB => initializes database => creates database file sqllite
// MIGRATE => migrates database
// SEED => seeds database
// RESOURCES => initializes go routines for server resources monitoring
// SCANNER => initializes scanner go routine => scans folders for new files
// ENCODER => initializes encoder go routine => encodes files
// COPIER => initializes copier go routine => copies files to temp folder
func Setup() {
	App()
	Db()
	Migrate()
	Seed()
	Resources()
	Scanner()
	Encoder()
	Copier()
}
