package grpcserver

import (
	"context"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	"google.golang.org/grpc"
)

type LibraryService struct {
	readerv1.UnimplementedLibraryServiceServer
	uploadBook          UploadBookHandler
	listCatalog         ListCatalogHandler
	getBook             GetBookHandler
	addToMyLibrary      AddToMyLibraryHandler
	listMyBooks         ListMyBooksHandler
	removeFromMyLibrary RemoveFromMyLibraryHandler
	deleteBook          DeleteBookHandler
}

type UploadBookHandler interface {
	UploadBook(grpc.ClientStreamingServer[readerv1.UploadBookRequest, readerv1.Book]) error
}

type ListCatalogHandler interface {
	ListCatalog(context.Context, *readerv1.ListCatalogRequest) (*readerv1.BookPage, error)
}

type GetBookHandler interface {
	GetBook(context.Context, *readerv1.GetBookRequest) (*readerv1.Book, error)
}

type AddToMyLibraryHandler interface {
	AddToMyLibrary(context.Context, *readerv1.AddToMyLibraryRequest) (*readerv1.UserBook, error)
}

type ListMyBooksHandler interface {
	ListMyBooks(context.Context, *readerv1.ListMyBooksRequest) (*readerv1.UserBookPage, error)
}

type RemoveFromMyLibraryHandler interface {
	RemoveFromMyLibrary(context.Context, *readerv1.RemoveFromMyLibraryRequest) (*readerv1.Empty, error)
}
type DeleteBookHandler interface {
	DeleteBook(context.Context, *readerv1.DeleteBookRequest) (*readerv1.Empty, error)
}

func NewLibraryService(
	uploadBook UploadBookHandler,
	listCatalog ListCatalogHandler,
	getBook GetBookHandler,
	addToMyLibrary AddToMyLibraryHandler,
	listMyBooks ListMyBooksHandler,
	removeFromMyLibrary RemoveFromMyLibraryHandler,
	deleteBook DeleteBookHandler,
) *LibraryService {
	return &LibraryService{
		uploadBook:          uploadBook,
		listCatalog:         listCatalog,
		getBook:             getBook,
		addToMyLibrary:      addToMyLibrary,
		listMyBooks:         listMyBooks,
		removeFromMyLibrary: removeFromMyLibrary,
		deleteBook:          deleteBook,
	}
}

func (s *LibraryService) DeleteBook(ctx context.Context, request *readerv1.DeleteBookRequest) (*readerv1.Empty, error) {
	return s.deleteBook.DeleteBook(ctx, request)
}

func (s *LibraryService) UploadBook(stream grpc.ClientStreamingServer[readerv1.UploadBookRequest, readerv1.Book]) error {
	return s.uploadBook.UploadBook(stream)
}

func (s *LibraryService) ListCatalog(ctx context.Context, request *readerv1.ListCatalogRequest) (*readerv1.BookPage, error) {
	return s.listCatalog.ListCatalog(ctx, request)
}

func (s *LibraryService) GetBook(ctx context.Context, request *readerv1.GetBookRequest) (*readerv1.Book, error) {
	return s.getBook.GetBook(ctx, request)
}

func (s *LibraryService) AddToMyLibrary(ctx context.Context, request *readerv1.AddToMyLibraryRequest) (*readerv1.UserBook, error) {
	return s.addToMyLibrary.AddToMyLibrary(ctx, request)
}

func (s *LibraryService) ListMyBooks(ctx context.Context, request *readerv1.ListMyBooksRequest) (*readerv1.UserBookPage, error) {
	return s.listMyBooks.ListMyBooks(ctx, request)
}

func (s *LibraryService) RemoveFromMyLibrary(ctx context.Context, request *readerv1.RemoveFromMyLibraryRequest) (*readerv1.Empty, error) {
	return s.removeFromMyLibrary.RemoveFromMyLibrary(ctx, request)
}
