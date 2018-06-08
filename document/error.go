package document

//ErrSchemaInfoNotFound error to indicate specified document schema is not found in database
type ErrSchemaInfoNotFound struct {
	msg string
}

func (err ErrSchemaInfoNotFound) Error() string { return err.msg }

//ErrSchemaInfoAlreadyExists error in indicate schema info already registered in database
type ErrSchemaInfoAlreadyExists struct {
	msg string
}

func (err ErrSchemaInfoAlreadyExists) Error() string { return err.msg }

//ErrDraftNotFound error in indicate schema has no draft in database
type ErrDraftNotFound struct {
	msg string
}

func (err ErrDraftNotFound) Error() string { return err.msg }
