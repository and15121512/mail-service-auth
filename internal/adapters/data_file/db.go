package data_file

import "context"

type DataFile struct {
}

func New(ctx context.Context, pgconn string) (*DataFile, error) {
	return &DataFile{}, nil
}
