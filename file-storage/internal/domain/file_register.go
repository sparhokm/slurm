package domain

type FileRegister struct {
	ID          string
	Filepath    string
	ContentType string
	Size        int64
	Version     int64
}
