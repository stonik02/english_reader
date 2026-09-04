package grpcserver

import (
	"context"
	"errors"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	domain "github.com/deniskrylov/english-reader/backend/internal/domain/reader"
	getadjacent "github.com/deniskrylov/english-reader/backend/internal/usecase/reader/getadjacentchapter"
	getchapter "github.com/deniskrylov/english-reader/backend/internal/usecase/reader/getchapter"
	getstate "github.com/deniskrylov/english-reader/backend/internal/usecase/reader/getreadingstate"
	getsettings "github.com/deniskrylov/english-reader/backend/internal/usecase/reader/getsettings"
	saveprogress "github.com/deniskrylov/english-reader/backend/internal/usecase/reader/saveprogress"
	updatesettings "github.com/deniskrylov/english-reader/backend/internal/usecase/reader/updatesettings"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReaderService struct {
	readerv1.UnimplementedReaderServiceServer
	parse          func(string) (string, error)
	state          *getstate.UseCase
	chapter        *getchapter.UseCase
	adjacent       *getadjacent.UseCase
	progress       *saveprogress.UseCase
	settings       *getsettings.UseCase
	updateSettings *updatesettings.UseCase
}

func NewReaderService(parse func(string) (string, error), state *getstate.UseCase, chapter *getchapter.UseCase, adjacent *getadjacent.UseCase, progress *saveprogress.UseCase, settings *getsettings.UseCase, updateSettings *updatesettings.UseCase) *ReaderService {
	return &ReaderService{parse: parse, state: state, chapter: chapter, adjacent: adjacent, progress: progress, settings: settings, updateSettings: updateSettings}
}
func (s *ReaderService) subject(v string) (string, error) {
	id, e := s.parse(v)
	if e != nil {
		return "", status.Error(codes.Unauthenticated, "invalid access token")
	}
	return id, nil
}
func (s *ReaderService) GetReadingState(c context.Context, q *readerv1.GetReadingStateRequest) (*readerv1.ReadingState, error) {
	id, e := s.subject(q.GetAccessToken())
	if e != nil {
		return nil, e
	}
	v, e := s.state.Execute(c, id, q.GetBookId())
	if e != nil {
		return nil, readerError(e)
	}
	return &readerv1.ReadingState{Chapter: chapter(v.Chapter), Progress: progress(v.Progress), Settings: settings(v.Settings)}, nil
}
func (s *ReaderService) GetChapter(c context.Context, q *readerv1.GetChapterRequest) (*readerv1.Chapter, error) {
	id, e := s.subject(q.GetAccessToken())
	if e != nil {
		return nil, e
	}
	v, e := s.chapter.Execute(c, id, q.GetBookId(), q.GetChapterId(), q.GetHref())
	if e != nil {
		return nil, readerError(e)
	}
	return chapter(v), nil
}
func (s *ReaderService) GetAdjacentChapter(c context.Context, q *readerv1.GetAdjacentChapterRequest) (*readerv1.Chapter, error) {
	id, e := s.subject(q.GetAccessToken())
	if e != nil {
		return nil, e
	}
	v, e := s.adjacent.Execute(c, id, q.GetBookId(), q.GetChapterId(), int(q.GetDirection()))
	if e != nil {
		return nil, readerError(e)
	}
	return chapter(v), nil
}
func (s *ReaderService) SaveReadingProgress(c context.Context, q *readerv1.SaveReadingProgressRequest) (*readerv1.ReadingProgress, error) {
	id, e := s.subject(q.GetAccessToken())
	if e != nil {
		return nil, e
	}
	v, e := s.progress.Execute(c, id, q.GetBookId(), saveprogress.Request{ChapterID: q.GetChapterId(), EPUBCFI: q.GetEpubCfi(), ProgressPercent: q.GetProgressPercent(), Revision: q.GetRevision()})
	if e != nil {
		return nil, readerError(e)
	}
	return progress(v), nil
}
func (s *ReaderService) GetReaderSettings(c context.Context, q *readerv1.GetReaderSettingsRequest) (*readerv1.ReaderSettings, error) {
	id, e := s.subject(q.GetAccessToken())
	if e != nil {
		return nil, e
	}
	v, e := s.settings.Execute(c, id)
	if e != nil {
		return nil, readerError(e)
	}
	return settings(v), nil
}
func (s *ReaderService) UpdateReaderSettings(c context.Context, q *readerv1.UpdateReaderSettingsRequest) (*readerv1.ReaderSettings, error) {
	id, e := s.subject(q.GetAccessToken())
	if e != nil {
		return nil, e
	}
	v, e := s.updateSettings.Execute(c, id, domain.Settings{FontScale: int(q.GetFontScale()), Theme: q.GetTheme(), LineHeight: q.GetLineHeight(), HighlightColor: q.GetHighlightColor()})
	if e != nil {
		return nil, readerError(e)
	}
	return settings(v), nil
}
func chapter(v domain.Chapter) *readerv1.Chapter {
	return &readerv1.Chapter{Id: v.ID, Href: v.Href, Sequence: int32(v.Sequence), StartCfi: v.StartCFI, EndCfi: v.EndCFI, SanitizedHtml: v.SanitizedHTML, TotalChapters: int32(v.TotalChapters)}
}
func progress(v domain.Progress) *readerv1.ReadingProgress {
	return &readerv1.ReadingProgress{ChapterId: v.ChapterID, EpubCfi: v.EPUBCFI, ProgressPercent: v.ProgressPercent, Revision: v.Revision}
}
func settings(v domain.Settings) *readerv1.ReaderSettings {
	return &readerv1.ReaderSettings{FontScale: int32(v.FontScale), Theme: v.Theme, LineHeight: v.LineHeight, HighlightColor: v.HighlightColor}
}
func readerError(e error) error {
	if errors.Is(e, domain.ErrInvalidInput) {
		return status.Error(codes.InvalidArgument, e.Error())
	}
	if errors.Is(e, domain.ErrNotFound) {
		return status.Error(codes.NotFound, e.Error())
	}
	if errors.Is(e, domain.ErrNotReady) {
		return status.Error(codes.FailedPrecondition, e.Error())
	}
	return status.Error(codes.Internal, "internal server error")
}
