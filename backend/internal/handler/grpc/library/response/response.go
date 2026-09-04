package response

import (
	"errors"
	"time"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/library"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Book(book domain.Book) *readerv1.Book {
	return &readerv1.Book{
		Id:               book.ID,
		Title:            book.Title,
		Author:           book.Author,
		Status:           book.Status,
		CoverUrl:         book.CoverURL,
		UploadedByUserId: book.UploadedByUserID,
		CreatedAt:        book.CreatedAt.Format(time.RFC3339Nano),
	}
}

func UserBook(userBook domain.UserBook) *readerv1.UserBook {
	return &readerv1.UserBook{
		Book:            Book(userBook.Book),
		AddedAt:         userBook.AddedAt.Format(time.RFC3339Nano),
		AddedVia:        userBook.AddedVia,
		ProgressPercent: userBook.ProgressPercent,
	}
}

func Error(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidUpload), errors.Is(err, domain.ErrTooLarge):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrNotReady):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
