package routing

//ErrNoMatch error to indicate no matching URL and method signature
type ErrNoMatch struct {
	Msg string
}

func (err ErrNoMatch) Error() string { return err.Msg }

//ErrInvalidInputData error to indicate invalid input data
type ErrInvalidInputData struct {
	Msg string
}

func (err ErrInvalidInputData) Error() string { return err.Msg }

//ErrNotAuthorize error to indicate requestor is not authorize to access the resource
type ErrNotAuthorize struct {
	Msg string
}

func (err ErrNotAuthorize) Error() string { return err.Msg }

//ErrNotAuthenticate error to indicate requestor is not login yet
type ErrNotAuthenticate struct {
	Msg string
}

func (err ErrNotAuthenticate) Error() string { return err.Msg }
