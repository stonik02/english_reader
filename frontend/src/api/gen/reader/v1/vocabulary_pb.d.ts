import * as jspb from 'google-protobuf'



export class VocabularySense extends jspb.Message {
  getId(): number;
  setId(value: number): VocabularySense;

  getPartOfSpeech(): string;
  setPartOfSpeech(value: string): VocabularySense;

  getTranslationsList(): Array<string>;
  setTranslationsList(value: Array<string>): VocabularySense;
  clearTranslationsList(): VocabularySense;
  addTranslations(value: string, index?: number): VocabularySense;

  getExampleEn(): string;
  setExampleEn(value: string): VocabularySense;

  getExampleRu(): string;
  setExampleRu(value: string): VocabularySense;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VocabularySense.AsObject;
  static toObject(includeInstance: boolean, msg: VocabularySense): VocabularySense.AsObject;
  static serializeBinaryToWriter(message: VocabularySense, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VocabularySense;
  static deserializeBinaryFromReader(message: VocabularySense, reader: jspb.BinaryReader): VocabularySense;
}

export namespace VocabularySense {
  export type AsObject = {
    id: number,
    partOfSpeech: string,
    translationsList: Array<string>,
    exampleEn: string,
    exampleRu: string,
  }
}

export class VocabularyEntry extends jspb.Message {
  getId(): string;
  setId(value: string): VocabularyEntry;

  getLemmaId(): number;
  setLemmaId(value: number): VocabularyEntry;

  getLemma(): string;
  setLemma(value: string): VocabularyEntry;

  getSourceForm(): string;
  setSourceForm(value: string): VocabularyEntry;

  getChosenSense(): VocabularySense | undefined;
  setChosenSense(value?: VocabularySense): VocabularyEntry;
  hasChosenSense(): boolean;
  clearChosenSense(): VocabularyEntry;

  getCreatedAt(): string;
  setCreatedAt(value: string): VocabularyEntry;

  getUpdatedAt(): string;
  setUpdatedAt(value: string): VocabularyEntry;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): VocabularyEntry.AsObject;
  static toObject(includeInstance: boolean, msg: VocabularyEntry): VocabularyEntry.AsObject;
  static serializeBinaryToWriter(message: VocabularyEntry, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): VocabularyEntry;
  static deserializeBinaryFromReader(message: VocabularyEntry, reader: jspb.BinaryReader): VocabularyEntry;
}

export namespace VocabularyEntry {
  export type AsObject = {
    id: string,
    lemmaId: number,
    lemma: string,
    sourceForm: string,
    chosenSense?: VocabularySense.AsObject,
    createdAt: string,
    updatedAt: string,
  }
}

export class SaveEntryRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): SaveEntryRequest;

  getLemmaId(): number;
  setLemmaId(value: number): SaveEntryRequest;

  getChosenSenseId(): number;
  setChosenSenseId(value: number): SaveEntryRequest;
  hasChosenSenseId(): boolean;
  clearChosenSenseId(): SaveEntryRequest;

  getSourceForm(): string;
  setSourceForm(value: string): SaveEntryRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveEntryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SaveEntryRequest): SaveEntryRequest.AsObject;
  static serializeBinaryToWriter(message: SaveEntryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveEntryRequest;
  static deserializeBinaryFromReader(message: SaveEntryRequest, reader: jspb.BinaryReader): SaveEntryRequest;
}

export namespace SaveEntryRequest {
  export type AsObject = {
    accessToken: string,
    lemmaId: number,
    chosenSenseId?: number,
    sourceForm: string,
  }

  export enum ChosenSenseIdCase { 
    _CHOSEN_SENSE_ID_NOT_SET = 0,
    CHOSEN_SENSE_ID = 3,
  }
}

export class SaveEntryResponse extends jspb.Message {
  getEntry(): VocabularyEntry | undefined;
  setEntry(value?: VocabularyEntry): SaveEntryResponse;
  hasEntry(): boolean;
  clearEntry(): SaveEntryResponse;

  getAlreadySaved(): boolean;
  setAlreadySaved(value: boolean): SaveEntryResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SaveEntryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SaveEntryResponse): SaveEntryResponse.AsObject;
  static serializeBinaryToWriter(message: SaveEntryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SaveEntryResponse;
  static deserializeBinaryFromReader(message: SaveEntryResponse, reader: jspb.BinaryReader): SaveEntryResponse;
}

export namespace SaveEntryResponse {
  export type AsObject = {
    entry?: VocabularyEntry.AsObject,
    alreadySaved: boolean,
  }
}

export class ListEntriesRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): ListEntriesRequest;

  getCursor(): string;
  setCursor(value: string): ListEntriesRequest;

  getLimit(): number;
  setLimit(value: number): ListEntriesRequest;

  getQuery(): string;
  setQuery(value: string): ListEntriesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEntriesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListEntriesRequest): ListEntriesRequest.AsObject;
  static serializeBinaryToWriter(message: ListEntriesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEntriesRequest;
  static deserializeBinaryFromReader(message: ListEntriesRequest, reader: jspb.BinaryReader): ListEntriesRequest;
}

export namespace ListEntriesRequest {
  export type AsObject = {
    accessToken: string,
    cursor: string,
    limit: number,
    query: string,
  }
}

export class ListEntriesResponse extends jspb.Message {
  getEntriesList(): Array<VocabularyEntry>;
  setEntriesList(value: Array<VocabularyEntry>): ListEntriesResponse;
  clearEntriesList(): ListEntriesResponse;
  addEntries(value?: VocabularyEntry, index?: number): VocabularyEntry;

  getNextCursor(): string;
  setNextCursor(value: string): ListEntriesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEntriesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListEntriesResponse): ListEntriesResponse.AsObject;
  static serializeBinaryToWriter(message: ListEntriesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEntriesResponse;
  static deserializeBinaryFromReader(message: ListEntriesResponse, reader: jspb.BinaryReader): ListEntriesResponse;
}

export namespace ListEntriesResponse {
  export type AsObject = {
    entriesList: Array<VocabularyEntry.AsObject>,
    nextCursor: string,
  }
}

export class DeleteEntryRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): DeleteEntryRequest;

  getEntryId(): string;
  setEntryId(value: string): DeleteEntryRequest;

  getLemmaId(): number;
  setLemmaId(value: number): DeleteEntryRequest;

  getTargetCase(): DeleteEntryRequest.TargetCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteEntryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteEntryRequest): DeleteEntryRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteEntryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteEntryRequest;
  static deserializeBinaryFromReader(message: DeleteEntryRequest, reader: jspb.BinaryReader): DeleteEntryRequest;
}

export namespace DeleteEntryRequest {
  export type AsObject = {
    accessToken: string,
    entryId: string,
    lemmaId: number,
  }

  export enum TargetCase { 
    TARGET_NOT_SET = 0,
    ENTRY_ID = 2,
    LEMMA_ID = 3,
  }
}

export class DeleteEntryResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteEntryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteEntryResponse): DeleteEntryResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteEntryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteEntryResponse;
  static deserializeBinaryFromReader(message: DeleteEntryResponse, reader: jspb.BinaryReader): DeleteEntryResponse;
}

export namespace DeleteEntryResponse {
  export type AsObject = {
  }
}

export class GetHighlightsRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): GetHighlightsRequest;

  getBookId(): string;
  setBookId(value: string): GetHighlightsRequest;

  getChapterId(): string;
  setChapterId(value: string): GetHighlightsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetHighlightsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetHighlightsRequest): GetHighlightsRequest.AsObject;
  static serializeBinaryToWriter(message: GetHighlightsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetHighlightsRequest;
  static deserializeBinaryFromReader(message: GetHighlightsRequest, reader: jspb.BinaryReader): GetHighlightsRequest;
}

export namespace GetHighlightsRequest {
  export type AsObject = {
    accessToken: string,
    bookId: string,
    chapterId: string,
  }
}

export class HighlightToken extends jspb.Message {
  getLemmaId(): number;
  setLemmaId(value: number): HighlightToken;

  getLemma(): string;
  setLemma(value: string): HighlightToken;

  getPositionsList(): Array<number>;
  setPositionsList(value: Array<number>): HighlightToken;
  clearPositionsList(): HighlightToken;
  addPositions(value: number, index?: number): HighlightToken;

  getMatchKind(): HighlightToken.MatchKind;
  setMatchKind(value: HighlightToken.MatchKind): HighlightToken;

  getTextsList(): Array<string>;
  setTextsList(value: Array<string>): HighlightToken;
  clearTextsList(): HighlightToken;
  addTexts(value: string, index?: number): HighlightToken;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): HighlightToken.AsObject;
  static toObject(includeInstance: boolean, msg: HighlightToken): HighlightToken.AsObject;
  static serializeBinaryToWriter(message: HighlightToken, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): HighlightToken;
  static deserializeBinaryFromReader(message: HighlightToken, reader: jspb.BinaryReader): HighlightToken;
}

export namespace HighlightToken {
  export type AsObject = {
    lemmaId: number,
    lemma: string,
    positionsList: Array<number>,
    matchKind: HighlightToken.MatchKind,
    textsList: Array<string>,
  }

  export enum MatchKind { 
    MATCH_KIND_UNSPECIFIED = 0,
    MATCH_KIND_LEMMA = 1,
    MATCH_KIND_EXACT_FALLBACK = 2,
  }
}

export class GetHighlightsResponse extends jspb.Message {
  getTokensList(): Array<HighlightToken>;
  setTokensList(value: Array<HighlightToken>): GetHighlightsResponse;
  clearTokensList(): GetHighlightsResponse;
  addTokens(value?: HighlightToken, index?: number): HighlightToken;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetHighlightsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetHighlightsResponse): GetHighlightsResponse.AsObject;
  static serializeBinaryToWriter(message: GetHighlightsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetHighlightsResponse;
  static deserializeBinaryFromReader(message: GetHighlightsResponse, reader: jspb.BinaryReader): GetHighlightsResponse;
}

export namespace GetHighlightsResponse {
  export type AsObject = {
    tokensList: Array<HighlightToken.AsObject>,
  }
}

