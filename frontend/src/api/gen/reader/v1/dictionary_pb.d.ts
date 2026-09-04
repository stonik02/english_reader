import * as jspb from 'google-protobuf'



export class LookupWordRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): LookupWordRequest;

  getBookId(): string;
  setBookId(value: string): LookupWordRequest;

  getChapterId(): string;
  setChapterId(value: string): LookupWordRequest;

  getSelectedText(): string;
  setSelectedText(value: string): LookupWordRequest;

  getSentenceText(): string;
  setSentenceText(value: string): LookupWordRequest;

  getEpubCfi(): string;
  setEpubCfi(value: string): LookupWordRequest;

  getSourceLanguage(): string;
  setSourceLanguage(value: string): LookupWordRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LookupWordRequest.AsObject;
  static toObject(includeInstance: boolean, msg: LookupWordRequest): LookupWordRequest.AsObject;
  static serializeBinaryToWriter(message: LookupWordRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LookupWordRequest;
  static deserializeBinaryFromReader(message: LookupWordRequest, reader: jspb.BinaryReader): LookupWordRequest;
}

export namespace LookupWordRequest {
  export type AsObject = {
    accessToken: string,
    bookId: string,
    chapterId: string,
    selectedText: string,
    sentenceText: string,
    epubCfi: string,
    sourceLanguage: string,
  }
}

export class DictionarySense extends jspb.Message {
  getId(): number;
  setId(value: number): DictionarySense;

  getPartOfSpeech(): string;
  setPartOfSpeech(value: string): DictionarySense;

  getTranslationsList(): Array<string>;
  setTranslationsList(value: Array<string>): DictionarySense;
  clearTranslationsList(): DictionarySense;
  addTranslations(value: string, index?: number): DictionarySense;

  getExampleEn(): string;
  setExampleEn(value: string): DictionarySense;

  getExampleRu(): string;
  setExampleRu(value: string): DictionarySense;

  getSourceUrl(): string;
  setSourceUrl(value: string): DictionarySense;

  getAttribution(): string;
  setAttribution(value: string): DictionarySense;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DictionarySense.AsObject;
  static toObject(includeInstance: boolean, msg: DictionarySense): DictionarySense.AsObject;
  static serializeBinaryToWriter(message: DictionarySense, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DictionarySense;
  static deserializeBinaryFromReader(message: DictionarySense, reader: jspb.BinaryReader): DictionarySense;
}

export namespace DictionarySense {
  export type AsObject = {
    id: number,
    partOfSpeech: string,
    translationsList: Array<string>,
    exampleEn: string,
    exampleRu: string,
    sourceUrl: string,
    attribution: string,
  }
}

export class SentenceTranslation extends jspb.Message {
  getTranslatedText(): string;
  setTranslatedText(value: string): SentenceTranslation;

  getProviderError(): string;
  setProviderError(value: string): SentenceTranslation;

  getResultCase(): SentenceTranslation.ResultCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SentenceTranslation.AsObject;
  static toObject(includeInstance: boolean, msg: SentenceTranslation): SentenceTranslation.AsObject;
  static serializeBinaryToWriter(message: SentenceTranslation, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SentenceTranslation;
  static deserializeBinaryFromReader(message: SentenceTranslation, reader: jspb.BinaryReader): SentenceTranslation;
}

export namespace SentenceTranslation {
  export type AsObject = {
    translatedText: string,
    providerError: string,
  }

  export enum ResultCase { 
    RESULT_NOT_SET = 0,
    TRANSLATED_TEXT = 1,
    PROVIDER_ERROR = 2,
  }
}

export class SourceMetadata extends jspb.Message {
  getSource(): string;
  setSource(value: string): SourceMetadata;

  getSourceVersion(): string;
  setSourceVersion(value: string): SourceMetadata;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SourceMetadata.AsObject;
  static toObject(includeInstance: boolean, msg: SourceMetadata): SourceMetadata.AsObject;
  static serializeBinaryToWriter(message: SourceMetadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SourceMetadata;
  static deserializeBinaryFromReader(message: SourceMetadata, reader: jspb.BinaryReader): SourceMetadata;
}

export namespace SourceMetadata {
  export type AsObject = {
    source: string,
    sourceVersion: string,
  }
}

export class LookupWordResponse extends jspb.Message {
  getNormalizedLemma(): string;
  setNormalizedLemma(value: string): LookupWordResponse;

  getSensesList(): Array<DictionarySense>;
  setSensesList(value: Array<DictionarySense>): LookupWordResponse;
  clearSensesList(): LookupWordResponse;
  addSenses(value?: DictionarySense, index?: number): DictionarySense;

  getSentenceTranslation(): SentenceTranslation | undefined;
  setSentenceTranslation(value?: SentenceTranslation): LookupWordResponse;
  hasSentenceTranslation(): boolean;
  clearSentenceTranslation(): LookupWordResponse;

  getContextVerified(): boolean;
  setContextVerified(value: boolean): LookupWordResponse;

  getSource(): SourceMetadata | undefined;
  setSource(value?: SourceMetadata): LookupWordResponse;
  hasSource(): boolean;
  clearSource(): LookupWordResponse;

  getAlreadySaved(): boolean;
  setAlreadySaved(value: boolean): LookupWordResponse;

  getLemmaId(): number;
  setLemmaId(value: number): LookupWordResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LookupWordResponse.AsObject;
  static toObject(includeInstance: boolean, msg: LookupWordResponse): LookupWordResponse.AsObject;
  static serializeBinaryToWriter(message: LookupWordResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LookupWordResponse;
  static deserializeBinaryFromReader(message: LookupWordResponse, reader: jspb.BinaryReader): LookupWordResponse;
}

export namespace LookupWordResponse {
  export type AsObject = {
    normalizedLemma: string,
    sensesList: Array<DictionarySense.AsObject>,
    sentenceTranslation?: SentenceTranslation.AsObject,
    contextVerified: boolean,
    source?: SourceMetadata.AsObject,
    alreadySaved: boolean,
    lemmaId: number,
  }
}

export class TranslateTextRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): TranslateTextRequest;

  getBookId(): string;
  setBookId(value: string): TranslateTextRequest;

  getChapterId(): string;
  setChapterId(value: string): TranslateTextRequest;

  getText(): string;
  setText(value: string): TranslateTextRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TranslateTextRequest.AsObject;
  static toObject(includeInstance: boolean, msg: TranslateTextRequest): TranslateTextRequest.AsObject;
  static serializeBinaryToWriter(message: TranslateTextRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TranslateTextRequest;
  static deserializeBinaryFromReader(message: TranslateTextRequest, reader: jspb.BinaryReader): TranslateTextRequest;
}

export namespace TranslateTextRequest {
  export type AsObject = {
    accessToken: string,
    bookId: string,
    chapterId: string,
    text: string,
  }
}

export class TranslateTextResponse extends jspb.Message {
  getSentenceTranslation(): SentenceTranslation | undefined;
  setSentenceTranslation(value?: SentenceTranslation): TranslateTextResponse;
  hasSentenceTranslation(): boolean;
  clearSentenceTranslation(): TranslateTextResponse;

  getContextVerified(): boolean;
  setContextVerified(value: boolean): TranslateTextResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TranslateTextResponse.AsObject;
  static toObject(includeInstance: boolean, msg: TranslateTextResponse): TranslateTextResponse.AsObject;
  static serializeBinaryToWriter(message: TranslateTextResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TranslateTextResponse;
  static deserializeBinaryFromReader(message: TranslateTextResponse, reader: jspb.BinaryReader): TranslateTextResponse;
}

export namespace TranslateTextResponse {
  export type AsObject = {
    sentenceTranslation?: SentenceTranslation.AsObject,
    contextVerified: boolean,
  }
}

