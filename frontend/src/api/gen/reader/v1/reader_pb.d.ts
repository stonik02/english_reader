import * as jspb from 'google-protobuf'



export class Chapter extends jspb.Message {
  getId(): string;
  setId(value: string): Chapter;

  getHref(): string;
  setHref(value: string): Chapter;

  getSequence(): number;
  setSequence(value: number): Chapter;

  getStartCfi(): string;
  setStartCfi(value: string): Chapter;

  getEndCfi(): string;
  setEndCfi(value: string): Chapter;

  getSanitizedHtml(): string;
  setSanitizedHtml(value: string): Chapter;

  getTotalChapters(): number;
  setTotalChapters(value: number): Chapter;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Chapter.AsObject;
  static toObject(includeInstance: boolean, msg: Chapter): Chapter.AsObject;
  static serializeBinaryToWriter(message: Chapter, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Chapter;
  static deserializeBinaryFromReader(message: Chapter, reader: jspb.BinaryReader): Chapter;
}

export namespace Chapter {
  export type AsObject = {
    id: string,
    href: string,
    sequence: number,
    startCfi: string,
    endCfi: string,
    sanitizedHtml: string,
    totalChapters: number,
  }
}

export class ReadingProgress extends jspb.Message {
  getChapterId(): string;
  setChapterId(value: string): ReadingProgress;

  getEpubCfi(): string;
  setEpubCfi(value: string): ReadingProgress;

  getProgressPercent(): number;
  setProgressPercent(value: number): ReadingProgress;

  getRevision(): number;
  setRevision(value: number): ReadingProgress;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReadingProgress.AsObject;
  static toObject(includeInstance: boolean, msg: ReadingProgress): ReadingProgress.AsObject;
  static serializeBinaryToWriter(message: ReadingProgress, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReadingProgress;
  static deserializeBinaryFromReader(message: ReadingProgress, reader: jspb.BinaryReader): ReadingProgress;
}

export namespace ReadingProgress {
  export type AsObject = {
    chapterId: string,
    epubCfi: string,
    progressPercent: number,
    revision: number,
  }
}

export class ReaderSettings extends jspb.Message {
  getFontScale(): number;
  setFontScale(value: number): ReaderSettings;

  getTheme(): string;
  setTheme(value: string): ReaderSettings;

  getLineHeight(): number;
  setLineHeight(value: number): ReaderSettings;

  getHighlightColor(): string;
  setHighlightColor(value: string): ReaderSettings;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReaderSettings.AsObject;
  static toObject(includeInstance: boolean, msg: ReaderSettings): ReaderSettings.AsObject;
  static serializeBinaryToWriter(message: ReaderSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReaderSettings;
  static deserializeBinaryFromReader(message: ReaderSettings, reader: jspb.BinaryReader): ReaderSettings;
}

export namespace ReaderSettings {
  export type AsObject = {
    fontScale: number,
    theme: string,
    lineHeight: number,
    highlightColor: string,
  }
}

export class ReadingState extends jspb.Message {
  getChapter(): Chapter | undefined;
  setChapter(value?: Chapter): ReadingState;
  hasChapter(): boolean;
  clearChapter(): ReadingState;

  getProgress(): ReadingProgress | undefined;
  setProgress(value?: ReadingProgress): ReadingState;
  hasProgress(): boolean;
  clearProgress(): ReadingState;

  getSettings(): ReaderSettings | undefined;
  setSettings(value?: ReaderSettings): ReadingState;
  hasSettings(): boolean;
  clearSettings(): ReadingState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ReadingState.AsObject;
  static toObject(includeInstance: boolean, msg: ReadingState): ReadingState.AsObject;
  static serializeBinaryToWriter(message: ReadingState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ReadingState;
  static deserializeBinaryFromReader(message: ReadingState, reader: jspb.BinaryReader): ReadingState;
}

export namespace ReadingState {
  export type AsObject = {
    chapter?: Chapter.AsObject,
    progress?: ReadingProgress.AsObject,
    settings?: ReaderSettings.AsObject,
  }
}

export class GetReadingStateRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): GetReadingStateRequest;

  getBookId(): string;
  setBookId(value: string): GetReadingStateRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetReadingStateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetReadingStateRequest): GetReadingStateRequest.AsObject;
  static serializeBinaryToWriter(message: GetReadingStateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetReadingStateRequest;
  static deserializeBinaryFromReader(message: GetReadingStateRequest, reader: jspb.BinaryReader): GetReadingStateRequest;
}

export namespace GetReadingStateRequest {
  export type AsObject = {
    accessToken: string,
    bookId: string,
  }
}

export class GetChapterRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): GetChapterRequest;

  getBookId(): string;
  setBookId(value: string): GetChapterRequest;

  getChapterId(): string;
  setChapterId(value: string): GetChapterRequest;

  getHref(): string;
  setHref(value: string): GetChapterRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetChapterRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetChapterRequest): GetChapterRequest.AsObject;
  static serializeBinaryToWriter(message: GetChapterRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetChapterRequest;
  static deserializeBinaryFromReader(message: GetChapterRequest, reader: jspb.BinaryReader): GetChapterRequest;
}

export namespace GetChapterRequest {
  export type AsObject = {
    accessToken: string,
    bookId: string,
    chapterId: string,
    href: string,
  }
}

export class GetAdjacentChapterRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): GetAdjacentChapterRequest;

  getBookId(): string;
  setBookId(value: string): GetAdjacentChapterRequest;

  getChapterId(): string;
  setChapterId(value: string): GetAdjacentChapterRequest;

  getDirection(): number;
  setDirection(value: number): GetAdjacentChapterRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetAdjacentChapterRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetAdjacentChapterRequest): GetAdjacentChapterRequest.AsObject;
  static serializeBinaryToWriter(message: GetAdjacentChapterRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetAdjacentChapterRequest;
  static deserializeBinaryFromReader(message: GetAdjacentChapterRequest, reader: jspb.BinaryReader): GetAdjacentChapterRequest;
}

export namespace GetAdjacentChapterRequest {
  export type AsObject = {
    accessToken: string,
    bookId: string,
    chapterId: string,
    direction: number,
  }
}

export class SaveReadingProgressRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): SaveReadingProgressRequest;

  getBookId(): string;
  setBookId(value: string): SaveReadingProgressRequest;

  getChapterId(): string;
  setChapterId(value: string): SaveReadingProgressRequest;

  getEpubCfi(): string;
  setEpubCfi(value: string): SaveReadingProgressRequest;

  getProgressPercent(): number;
  setProgressPercent(value: number): SaveReadingProgressRequest;

  getRevision(): number;
  setRevision(value: number): SaveReadingProgressRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveReadingProgressRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SaveReadingProgressRequest): SaveReadingProgressRequest.AsObject;
  static serializeBinaryToWriter(message: SaveReadingProgressRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveReadingProgressRequest;
  static deserializeBinaryFromReader(message: SaveReadingProgressRequest, reader: jspb.BinaryReader): SaveReadingProgressRequest;
}

export namespace SaveReadingProgressRequest {
  export type AsObject = {
    accessToken: string,
    bookId: string,
    chapterId: string,
    epubCfi: string,
    progressPercent: number,
    revision: number,
  }
}

export class GetReaderSettingsRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): GetReaderSettingsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetReaderSettingsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetReaderSettingsRequest): GetReaderSettingsRequest.AsObject;
  static serializeBinaryToWriter(message: GetReaderSettingsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetReaderSettingsRequest;
  static deserializeBinaryFromReader(message: GetReaderSettingsRequest, reader: jspb.BinaryReader): GetReaderSettingsRequest;
}

export namespace GetReaderSettingsRequest {
  export type AsObject = {
    accessToken: string,
  }
}

export class UpdateReaderSettingsRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): UpdateReaderSettingsRequest;

  getFontScale(): number;
  setFontScale(value: number): UpdateReaderSettingsRequest;

  getTheme(): string;
  setTheme(value: string): UpdateReaderSettingsRequest;

  getLineHeight(): number;
  setLineHeight(value: number): UpdateReaderSettingsRequest;

  getHighlightColor(): string;
  setHighlightColor(value: string): UpdateReaderSettingsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateReaderSettingsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateReaderSettingsRequest): UpdateReaderSettingsRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateReaderSettingsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateReaderSettingsRequest;
  static deserializeBinaryFromReader(message: UpdateReaderSettingsRequest, reader: jspb.BinaryReader): UpdateReaderSettingsRequest;
}

export namespace UpdateReaderSettingsRequest {
  export type AsObject = {
    accessToken: string,
    fontScale: number,
    theme: string,
    lineHeight: number,
    highlightColor: string,
  }
}

