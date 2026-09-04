import * as jspb from 'google-protobuf'

import * as reader_v1_auth_pb from '../../reader/v1/auth_pb'; // proto import: "reader/v1/auth.proto"


export class UploadBookMetadata extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): UploadBookMetadata;

  getFilename(): string;
  setFilename(value: string): UploadBookMetadata;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadBookMetadata.AsObject;
  static toObject(includeInstance: boolean, msg: UploadBookMetadata): UploadBookMetadata.AsObject;
  static serializeBinaryToWriter(message: UploadBookMetadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadBookMetadata;
  static deserializeBinaryFromReader(message: UploadBookMetadata, reader: jspb.BinaryReader): UploadBookMetadata;
}

export namespace UploadBookMetadata {
  export type AsObject = {
    accessToken: string,
    filename: string,
  }
}

export class UploadBookRequest extends jspb.Message {
  getMetadata(): UploadBookMetadata | undefined;
  setMetadata(value?: UploadBookMetadata): UploadBookRequest;
  hasMetadata(): boolean;
  clearMetadata(): UploadBookRequest;

  getChunk(): Uint8Array | string;
  getChunk_asU8(): Uint8Array;
  getChunk_asB64(): string;
  setChunk(value: Uint8Array | string): UploadBookRequest;

  getPayloadCase(): UploadBookRequest.PayloadCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadBookRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UploadBookRequest): UploadBookRequest.AsObject;
  static serializeBinaryToWriter(message: UploadBookRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadBookRequest;
  static deserializeBinaryFromReader(message: UploadBookRequest, reader: jspb.BinaryReader): UploadBookRequest;
}

export namespace UploadBookRequest {
  export type AsObject = {
    metadata?: UploadBookMetadata.AsObject,
    chunk: Uint8Array | string,
  }

  export enum PayloadCase { 
    PAYLOAD_NOT_SET = 0,
    METADATA = 1,
    CHUNK = 2,
  }
}

export class Book extends jspb.Message {
  getId(): string;
  setId(value: string): Book;

  getTitle(): string;
  setTitle(value: string): Book;

  getAuthor(): string;
  setAuthor(value: string): Book;

  getStatus(): string;
  setStatus(value: string): Book;

  getUploadedByUserId(): string;
  setUploadedByUserId(value: string): Book;

  getCreatedAt(): string;
  setCreatedAt(value: string): Book;

  getCoverUrl(): string;
  setCoverUrl(value: string): Book;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Book.AsObject;
  static toObject(includeInstance: boolean, msg: Book): Book.AsObject;
  static serializeBinaryToWriter(message: Book, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Book;
  static deserializeBinaryFromReader(message: Book, reader: jspb.BinaryReader): Book;
}

export namespace Book {
  export type AsObject = {
    id: string,
    title: string,
    author: string,
    status: string,
    uploadedByUserId: string,
    createdAt: string,
    coverUrl: string,
  }
}

export class UserBook extends jspb.Message {
  getBook(): Book | undefined;
  setBook(value?: Book): UserBook;
  hasBook(): boolean;
  clearBook(): UserBook;

  getAddedAt(): string;
  setAddedAt(value: string): UserBook;

  getAddedVia(): string;
  setAddedVia(value: string): UserBook;

  getProgressPercent(): number;
  setProgressPercent(value: number): UserBook;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserBook.AsObject;
  static toObject(includeInstance: boolean, msg: UserBook): UserBook.AsObject;
  static serializeBinaryToWriter(message: UserBook, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserBook;
  static deserializeBinaryFromReader(message: UserBook, reader: jspb.BinaryReader): UserBook;
}

export namespace UserBook {
  export type AsObject = {
    book?: Book.AsObject,
    addedAt: string,
    addedVia: string,
    progressPercent: number,
  }
}

export class BookPage extends jspb.Message {
  getBooksList(): Array<Book>;
  setBooksList(value: Array<Book>): BookPage;
  clearBooksList(): BookPage;
  addBooks(value?: Book, index?: number): Book;

  getNextCursor(): string;
  setNextCursor(value: string): BookPage;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BookPage.AsObject;
  static toObject(includeInstance: boolean, msg: BookPage): BookPage.AsObject;
  static serializeBinaryToWriter(message: BookPage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BookPage;
  static deserializeBinaryFromReader(message: BookPage, reader: jspb.BinaryReader): BookPage;
}

export namespace BookPage {
  export type AsObject = {
    booksList: Array<Book.AsObject>,
    nextCursor: string,
  }
}

export class UserBookPage extends jspb.Message {
  getBooksList(): Array<UserBook>;
  setBooksList(value: Array<UserBook>): UserBookPage;
  clearBooksList(): UserBookPage;
  addBooks(value?: UserBook, index?: number): UserBook;

  getNextCursor(): string;
  setNextCursor(value: string): UserBookPage;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserBookPage.AsObject;
  static toObject(includeInstance: boolean, msg: UserBookPage): UserBookPage.AsObject;
  static serializeBinaryToWriter(message: UserBookPage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserBookPage;
  static deserializeBinaryFromReader(message: UserBookPage, reader: jspb.BinaryReader): UserBookPage;
}

export namespace UserBookPage {
  export type AsObject = {
    booksList: Array<UserBook.AsObject>,
    nextCursor: string,
  }
}

export class ListCatalogRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): ListCatalogRequest;

  getCursor(): string;
  setCursor(value: string): ListCatalogRequest;

  getLimit(): number;
  setLimit(value: number): ListCatalogRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListCatalogRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListCatalogRequest): ListCatalogRequest.AsObject;
  static serializeBinaryToWriter(message: ListCatalogRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListCatalogRequest;
  static deserializeBinaryFromReader(message: ListCatalogRequest, reader: jspb.BinaryReader): ListCatalogRequest;
}

export namespace ListCatalogRequest {
  export type AsObject = {
    accessToken: string,
    cursor: string,
    limit: number,
  }
}

export class GetBookRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): GetBookRequest;

  getBookId(): string;
  setBookId(value: string): GetBookRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetBookRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetBookRequest): GetBookRequest.AsObject;
  static serializeBinaryToWriter(message: GetBookRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetBookRequest;
  static deserializeBinaryFromReader(message: GetBookRequest, reader: jspb.BinaryReader): GetBookRequest;
}

export namespace GetBookRequest {
  export type AsObject = {
    accessToken: string,
    bookId: string,
  }
}

export class AddToMyLibraryRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): AddToMyLibraryRequest;

  getBookId(): string;
  setBookId(value: string): AddToMyLibraryRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddToMyLibraryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddToMyLibraryRequest): AddToMyLibraryRequest.AsObject;
  static serializeBinaryToWriter(message: AddToMyLibraryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddToMyLibraryRequest;
  static deserializeBinaryFromReader(message: AddToMyLibraryRequest, reader: jspb.BinaryReader): AddToMyLibraryRequest;
}

export namespace AddToMyLibraryRequest {
  export type AsObject = {
    accessToken: string,
    bookId: string,
  }
}

export class ListMyBooksRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): ListMyBooksRequest;

  getCursor(): string;
  setCursor(value: string): ListMyBooksRequest;

  getLimit(): number;
  setLimit(value: number): ListMyBooksRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListMyBooksRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListMyBooksRequest): ListMyBooksRequest.AsObject;
  static serializeBinaryToWriter(message: ListMyBooksRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListMyBooksRequest;
  static deserializeBinaryFromReader(message: ListMyBooksRequest, reader: jspb.BinaryReader): ListMyBooksRequest;
}

export namespace ListMyBooksRequest {
  export type AsObject = {
    accessToken: string,
    cursor: string,
    limit: number,
  }
}

export class RemoveFromMyLibraryRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): RemoveFromMyLibraryRequest;

  getBookId(): string;
  setBookId(value: string): RemoveFromMyLibraryRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RemoveFromMyLibraryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RemoveFromMyLibraryRequest): RemoveFromMyLibraryRequest.AsObject;
  static serializeBinaryToWriter(message: RemoveFromMyLibraryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RemoveFromMyLibraryRequest;
  static deserializeBinaryFromReader(message: RemoveFromMyLibraryRequest, reader: jspb.BinaryReader): RemoveFromMyLibraryRequest;
}

export namespace RemoveFromMyLibraryRequest {
  export type AsObject = {
    accessToken: string,
    bookId: string,
  }
}

export class DeleteBookRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): DeleteBookRequest;

  getBookId(): string;
  setBookId(value: string): DeleteBookRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteBookRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteBookRequest): DeleteBookRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteBookRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteBookRequest;
  static deserializeBinaryFromReader(message: DeleteBookRequest, reader: jspb.BinaryReader): DeleteBookRequest;
}

export namespace DeleteBookRequest {
  export type AsObject = {
    accessToken: string,
    bookId: string,
  }
}

