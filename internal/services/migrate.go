package services

func (s *Conn) AutoMigrate() error {
	//if s.db.Migrator().HasTable(&User{}) {
	//	return nil
	//}

	err := s.Db.Migrator().DropTable()
	if err != nil {
		return err
	}

	// AutoMigrate function will ONLY create tables, missing columns and missing indexes, and WON'T change existing column's type or delete unused columns
	err = s.Db.Migrator().AutoMigrate()
	if err != nil {
		// If there is an error while migrating, log the error message and stop the program
		return err
	}
	return nil
}
