package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	readerv1 "github.com/deniskrylov/english-reader/backend/gen/reader/v1"
	"github.com/deniskrylov/english-reader/backend/internal/config"
	"github.com/deniskrylov/english-reader/backend/internal/database"
	ggetme "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/auth/getme"
	glogin "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/auth/login"
	glogout "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/auth/logout"
	grefresh "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/auth/refresh"
	gregistration "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/auth/register"
	dictionaryhandler "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/dictionary/lookupword"
	grpcLibraryAdd "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/library/addtomylibrary"
	grpcLibraryDelete "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/library/deletebook"
	grpcLibraryGet "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/library/getbook"
	grpcLibraryCatalog "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/library/listcatalog"
	grpcLibraryMine "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/library/listmybooks"
	grpcLibraryRemove "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/library/removefrommylibrary"
	grpcLibraryUpload "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/library/uploadbook"
	vocabularydelete "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/vocabulary/deleteentry"
	vocabularyhighlights "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/vocabulary/gethighlights"
	vocabularylist "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/vocabulary/listentries"
	vocabularysave "github.com/deniskrylov/english-reader/backend/internal/handler/grpc/vocabulary/saveentry"
	hgetme "github.com/deniskrylov/english-reader/backend/internal/handler/http/auth/getme"
	hlogin "github.com/deniskrylov/english-reader/backend/internal/handler/http/auth/login"
	hlogout "github.com/deniskrylov/english-reader/backend/internal/handler/http/auth/logout"
	hrefresh "github.com/deniskrylov/english-reader/backend/internal/handler/http/auth/refresh"
	hregister "github.com/deniskrylov/english-reader/backend/internal/handler/http/auth/register"
	ladd "github.com/deniskrylov/english-reader/backend/internal/handler/http/library/addtomylibrary"
	ldelete "github.com/deniskrylov/english-reader/backend/internal/handler/http/library/deletebook"
	lget "github.com/deniskrylov/english-reader/backend/internal/handler/http/library/getbook"
	lcover "github.com/deniskrylov/english-reader/backend/internal/handler/http/library/getcover"
	lcatalog "github.com/deniskrylov/english-reader/backend/internal/handler/http/library/listcatalog"
	lmine "github.com/deniskrylov/english-reader/backend/internal/handler/http/library/listmybooks"
	lremove "github.com/deniskrylov/english-reader/backend/internal/handler/http/library/removefrommylibrary"
	lupload "github.com/deniskrylov/english-reader/backend/internal/handler/http/library/upload"
	radjacent "github.com/deniskrylov/english-reader/backend/internal/handler/http/reader/getadjacentchapter"
	rchapter "github.com/deniskrylov/english-reader/backend/internal/handler/http/reader/getchapter"
	rstate "github.com/deniskrylov/english-reader/backend/internal/handler/http/reader/getreadingstate"
	rsettings "github.com/deniskrylov/english-reader/backend/internal/handler/http/reader/getsettings"
	rprogress "github.com/deniskrylov/english-reader/backend/internal/handler/http/reader/saveprogress"
	rupdatesettings "github.com/deniskrylov/english-reader/backend/internal/handler/http/reader/updatesettings"
	repoauth "github.com/deniskrylov/english-reader/backend/internal/repository/postgres/auth"
	repodictionary "github.com/deniskrylov/english-reader/backend/internal/repository/postgres/dictionary"
	repolibrary "github.com/deniskrylov/english-reader/backend/internal/repository/postgres/library"
	reporeader "github.com/deniskrylov/english-reader/backend/internal/repository/postgres/reader"
	repovocabulary "github.com/deniskrylov/english-reader/backend/internal/repository/postgres/vocabulary"
	"github.com/deniskrylov/english-reader/backend/internal/service/epubparser"
	"github.com/deniskrylov/english-reader/backend/internal/service/epubstorage"
	"github.com/deniskrylov/english-reader/backend/internal/service/libretranslate"
	"github.com/deniskrylov/english-reader/backend/internal/service/morphology"
	"github.com/deniskrylov/english-reader/backend/internal/service/password"
	"github.com/deniskrylov/english-reader/backend/internal/service/refreshtoken"
	"github.com/deniskrylov/english-reader/backend/internal/service/token"
	"github.com/deniskrylov/english-reader/backend/internal/service/tokenizer"
	"github.com/deniskrylov/english-reader/backend/internal/service/wordnormalizer"
	"github.com/deniskrylov/english-reader/backend/internal/transport/grpcserver"
	"github.com/deniskrylov/english-reader/backend/internal/transport/httpserver"
	getme "github.com/deniskrylov/english-reader/backend/internal/usecase/auth/getme"
	login "github.com/deniskrylov/english-reader/backend/internal/usecase/auth/login"
	logout "github.com/deniskrylov/english-reader/backend/internal/usecase/auth/logout"
	refresh "github.com/deniskrylov/english-reader/backend/internal/usecase/auth/refresh"
	register "github.com/deniskrylov/english-reader/backend/internal/usecase/auth/register"
	dictionarylookup "github.com/deniskrylov/english-reader/backend/internal/usecase/dictionary/lookupword"
	libraryadd "github.com/deniskrylov/english-reader/backend/internal/usecase/library/addtomylibrary"
	librarydelete "github.com/deniskrylov/english-reader/backend/internal/usecase/library/deletebook"
	libraryget "github.com/deniskrylov/english-reader/backend/internal/usecase/library/getbook"
	librarycover "github.com/deniskrylov/english-reader/backend/internal/usecase/library/getcover"
	librarycatalog "github.com/deniskrylov/english-reader/backend/internal/usecase/library/listcatalog"
	librarymine "github.com/deniskrylov/english-reader/backend/internal/usecase/library/listmybooks"
	libraryremove "github.com/deniskrylov/english-reader/backend/internal/usecase/library/removefrommylibrary"
	libraryupload "github.com/deniskrylov/english-reader/backend/internal/usecase/library/upload"
	readeradjacent "github.com/deniskrylov/english-reader/backend/internal/usecase/reader/getadjacentchapter"
	readerchapter "github.com/deniskrylov/english-reader/backend/internal/usecase/reader/getchapter"
	readerstate "github.com/deniskrylov/english-reader/backend/internal/usecase/reader/getreadingstate"
	readersettings "github.com/deniskrylov/english-reader/backend/internal/usecase/reader/getsettings"
	readerprogress "github.com/deniskrylov/english-reader/backend/internal/usecase/reader/saveprogress"
	readerupdatesettings "github.com/deniskrylov/english-reader/backend/internal/usecase/reader/updatesettings"
	vocabularydeleteusecase "github.com/deniskrylov/english-reader/backend/internal/usecase/vocabulary/deleteentry"
	vocabularyhighlightsusecase "github.com/deniskrylov/english-reader/backend/internal/usecase/vocabulary/gethighlights"
	vocabularylistusecase "github.com/deniskrylov/english-reader/backend/internal/usecase/vocabulary/listentries"
	vocabularysaveusecase "github.com/deniskrylov/english-reader/backend/internal/usecase/vocabulary/saveentry"
	"github.com/deniskrylov/english-reader/backend/internal/worker/epubingestion"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	if err := run(logger); err != nil {
		logger.Error("server stopped with error", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	go epubingestion.New(repolibrary.New(pool), epubparser.New(), epubstorage.New(cfg.BookStoragePath, epubstorage.DefaultMaxSize), logger, time.Second).Run(ctx)

	grpcTransport := grpcserver.New(logger)
	readerv1.RegisterAuthServiceServer(grpcTransport.GRPC, newGRPCAuthService(pool, cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL))
	readerv1.RegisterReaderServiceServer(grpcTransport.GRPC, newGRPCReaderService(pool, cfg.JWTSecret, cfg.AccessTokenTTL))
	readerv1.RegisterLibraryServiceServer(grpcTransport.GRPC, newGRPCLibraryService(pool, cfg.JWTSecret, cfg.AccessTokenTTL, cfg.BookStoragePath))
	readerv1.RegisterDictionaryServiceServer(grpcTransport.GRPC, newGRPCDictionaryService(pool, cfg))
	readerv1.RegisterVocabularyServiceServer(grpcTransport.GRPC, newGRPCVocabularyService(pool, cfg.JWTSecret, cfg.AccessTokenTTL, cfg.LookupWordMaxLength))
	grpcTransport.SetServing(true)
	grpcListener, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		return err
	}
	defer grpcListener.Close()

	httpTransport := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           newHTTPHandler(pool, cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL, cfg.BookStoragePath, cfg.FrontendOrigin, cfg.Environment != "development"),
		ReadHeaderTimeout: cfg.ShutdownTimeout,
	}

	errChannel := make(chan error, 2)
	go func() {
		logger.Info("grpc server started", "address", cfg.GRPCAddr, "environment", cfg.Environment)
		if serveErr := grpcTransport.GRPC.Serve(grpcListener); serveErr != nil {
			errChannel <- serveErr
		}
	}()
	go func() {
		logger.Info("http server started", "address", cfg.HTTPAddr)
		if serveErr := httpTransport.ListenAndServe(); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			errChannel <- serveErr
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	case serveErr := <-errChannel:
		return serveErr
	}

	grpcTransport.SetServing(false)
	grpcTransport.GRPC.GracefulStop()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	return httpTransport.Shutdown(shutdownCtx)
}

func newHTTPHandler(pool *pgxpool.Pool, secret string, accessTTL, refreshTTL time.Duration, bookStoragePath, frontendOrigin string, secureCookie bool) http.Handler {
	repository := repoauth.New(pool)
	passwordService := password.New()
	tokenService := token.New(secret, accessTTL)
	refreshService := refreshtoken.New()
	registerUseCase := register.New(repository, passwordService, tokenService, refreshService, refreshTTL)
	loginUseCase := login.New(repository, repository, passwordService, tokenService, refreshService, refreshTTL)
	refreshUseCase := refresh.New(repository, repository, repository, tokenService, refreshService, refreshTTL)
	logoutUseCase := logout.New(repository, refreshService)
	getMeUseCase := getme.New(tokenService, repository)
	libraryRepository := repolibrary.New(pool)
	libraryStorage := epubstorage.New(bookStoragePath, epubstorage.DefaultMaxSize)
	readerRepository := reporeader.New(pool)
	return httpserver.New(pool, frontendOrigin,
		hregister.New(registerUseCase, secureCookie, refreshTTL),
		hlogin.New(loginUseCase, secureCookie, refreshTTL),
		hrefresh.New(refreshUseCase, secureCookie, refreshTTL),
		hlogout.New(logoutUseCase, secureCookie),
		hgetme.New(getMeUseCase),
		lupload.New(libraryupload.New(libraryStorage, libraryRepository), tokenService),
		lcatalog.New(librarycatalog.New(libraryRepository)),
		lget.New(libraryget.New(libraryRepository)),
		lcover.New(librarycover.New(libraryRepository, libraryStorage), tokenService),
		ldelete.New(librarydelete.New(libraryRepository, libraryStorage), tokenService),
		ladd.New(libraryadd.New(libraryRepository), tokenService),
		lmine.New(librarymine.New(libraryRepository), tokenService),
		lremove.New(libraryremove.New(libraryRepository), tokenService),
		rstate.New(readerstate.New(readerRepository), tokenService),
		rchapter.New(readerchapter.New(readerRepository), tokenService),
		radjacent.New(readeradjacent.New(readerRepository), tokenService),
		rprogress.New(readerprogress.New(readerRepository), tokenService),
		rsettings.New(readersettings.New(readerRepository), tokenService),
		rupdatesettings.New(readerupdatesettings.New(readerRepository), tokenService),
	)
}

func newGRPCAuthService(pool *pgxpool.Pool, secret string, accessTTL, refreshTTL time.Duration) *grpcserver.AuthService {
	repository := repoauth.New(pool)
	passwordService := password.New()
	tokenService := token.New(secret, accessTTL)
	refreshService := refreshtoken.New()
	return grpcserver.NewAuthService(
		gregistration.New(register.New(repository, passwordService, tokenService, refreshService, refreshTTL)),
		glogin.New(login.New(repository, repository, passwordService, tokenService, refreshService, refreshTTL)),
		grefresh.New(refresh.New(repository, repository, repository, tokenService, refreshService, refreshTTL)),
		glogout.New(logout.New(repository, refreshService)),
		ggetme.New(getme.New(tokenService, repository)),
	)
}

func newGRPCReaderService(pool *pgxpool.Pool, secret string, accessTTL time.Duration) *grpcserver.ReaderService {
	repository := reporeader.New(pool)
	tokens := token.New(secret, accessTTL)

	return grpcserver.NewReaderService(
		tokens.Parse,
		readerstate.New(repository),
		readerchapter.New(repository),
		readeradjacent.New(repository),
		readerprogress.New(repository),
		readersettings.New(repository),
		readerupdatesettings.New(repository),
	)
}

func newGRPCDictionaryService(pool *pgxpool.Pool, cfg config.Config) *grpcserver.DictionaryService {
	tokens := token.New(cfg.JWTSecret, cfg.AccessTokenTTL)
	lookup := dictionarylookup.New(
		wordnormalizer.New(cfg.LookupWordMaxLength),
		morphology.New(),
		repodictionary.New(pool),
		repodictionary.New(pool),
		repovocabulary.New(pool),
		libretranslate.New(cfg.TranslateURL, cfg.TranslateTimeout, cfg.LookupSentenceMaxLength),
		"libretranslate-en-ru",
		cfg.TranslateCacheTTL,
	)

	return grpcserver.NewDictionaryService(dictionaryhandler.New(lookup, tokens))
}

func newGRPCLibraryService(pool *pgxpool.Pool, secret string, accessTTL time.Duration, storagePath string) *grpcserver.LibraryService {
	repository := repolibrary.New(pool)
	tokens := token.New(secret, accessTTL)
	upload := libraryupload.New(epubstorage.New(storagePath, epubstorage.DefaultMaxSize), repository)

	return grpcserver.NewLibraryService(
		grpcLibraryUpload.New(upload, tokens),
		grpcLibraryCatalog.New(librarycatalog.New(repository)),
		grpcLibraryGet.New(libraryget.New(repository)),
		grpcLibraryAdd.New(libraryadd.New(repository), tokens),
		grpcLibraryMine.New(librarymine.New(repository), tokens),
		grpcLibraryRemove.New(libraryremove.New(repository), tokens),
		grpcLibraryDelete.New(librarydelete.New(repository, epubstorage.New(storagePath, epubstorage.DefaultMaxSize)), tokens),
	)
}

func newGRPCVocabularyService(pool *pgxpool.Pool, secret string, accessTTL time.Duration, maxWordLength int) *grpcserver.VocabularyService {
	repository := repovocabulary.New(pool)
	tokens := token.New(secret, accessTTL)
	return grpcserver.NewVocabularyService(
		vocabularysave.New(vocabularysaveusecase.New(repository), tokens),
		vocabularylist.New(vocabularylistusecase.New(repository, wordnormalizer.New(maxWordLength)), tokens),
		vocabularydelete.New(vocabularydeleteusecase.New(repository), tokens),
		vocabularyhighlights.New(
			vocabularyhighlightsusecase.New(
				repository,
				repository,
				tokenizer.New(wordnormalizer.New(maxWordLength), morphology.New(), 1_000_000, 50_000),
				2_000,
				10_000,
			),
			tokens,
		),
	)
}
