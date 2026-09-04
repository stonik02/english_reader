import * as jspb from 'google-protobuf'



export class Empty extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Empty.AsObject;
  static toObject(includeInstance: boolean, msg: Empty): Empty.AsObject;
  static serializeBinaryToWriter(message: Empty, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Empty;
  static deserializeBinaryFromReader(message: Empty, reader: jspb.BinaryReader): Empty;
}

export namespace Empty {
  export type AsObject = {
  }
}

export class User extends jspb.Message {
  getId(): string;
  setId(value: string): User;

  getEmail(): string;
  setEmail(value: string): User;

  getCreatedAt(): string;
  setCreatedAt(value: string): User;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): User.AsObject;
  static toObject(includeInstance: boolean, msg: User): User.AsObject;
  static serializeBinaryToWriter(message: User, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): User;
  static deserializeBinaryFromReader(message: User, reader: jspb.BinaryReader): User;
}

export namespace User {
  export type AsObject = {
    id: string,
    email: string,
    createdAt: string,
  }
}

export class AuthResponse extends jspb.Message {
  getUser(): User | undefined;
  setUser(value?: User): AuthResponse;
  hasUser(): boolean;
  clearUser(): AuthResponse;

  getAccessToken(): string;
  setAccessToken(value: string): AuthResponse;

  getRefreshToken(): string;
  setRefreshToken(value: string): AuthResponse;

  getAccessTokenExpiresUnix(): number;
  setAccessTokenExpiresUnix(value: number): AuthResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AuthResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AuthResponse): AuthResponse.AsObject;
  static serializeBinaryToWriter(message: AuthResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AuthResponse;
  static deserializeBinaryFromReader(message: AuthResponse, reader: jspb.BinaryReader): AuthResponse;
}

export namespace AuthResponse {
  export type AsObject = {
    user?: User.AsObject,
    accessToken: string,
    refreshToken: string,
    accessTokenExpiresUnix: number,
  }
}

export class RegisterRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): RegisterRequest;

  getPassword(): string;
  setPassword(value: string): RegisterRequest;

  getDeviceLabel(): string;
  setDeviceLabel(value: string): RegisterRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterRequest): RegisterRequest.AsObject;
  static serializeBinaryToWriter(message: RegisterRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterRequest;
  static deserializeBinaryFromReader(message: RegisterRequest, reader: jspb.BinaryReader): RegisterRequest;
}

export namespace RegisterRequest {
  export type AsObject = {
    email: string,
    password: string,
    deviceLabel: string,
  }
}

export class LoginRequest extends jspb.Message {
  getEmail(): string;
  setEmail(value: string): LoginRequest;

  getPassword(): string;
  setPassword(value: string): LoginRequest;

  getDeviceLabel(): string;
  setDeviceLabel(value: string): LoginRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LoginRequest.AsObject;
  static toObject(includeInstance: boolean, msg: LoginRequest): LoginRequest.AsObject;
  static serializeBinaryToWriter(message: LoginRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LoginRequest;
  static deserializeBinaryFromReader(message: LoginRequest, reader: jspb.BinaryReader): LoginRequest;
}

export namespace LoginRequest {
  export type AsObject = {
    email: string,
    password: string,
    deviceLabel: string,
  }
}

export class RefreshRequest extends jspb.Message {
  getRefreshToken(): string;
  setRefreshToken(value: string): RefreshRequest;

  getDeviceLabel(): string;
  setDeviceLabel(value: string): RefreshRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RefreshRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RefreshRequest): RefreshRequest.AsObject;
  static serializeBinaryToWriter(message: RefreshRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RefreshRequest;
  static deserializeBinaryFromReader(message: RefreshRequest, reader: jspb.BinaryReader): RefreshRequest;
}

export namespace RefreshRequest {
  export type AsObject = {
    refreshToken: string,
    deviceLabel: string,
  }
}

export class LogoutRequest extends jspb.Message {
  getRefreshToken(): string;
  setRefreshToken(value: string): LogoutRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LogoutRequest.AsObject;
  static toObject(includeInstance: boolean, msg: LogoutRequest): LogoutRequest.AsObject;
  static serializeBinaryToWriter(message: LogoutRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LogoutRequest;
  static deserializeBinaryFromReader(message: LogoutRequest, reader: jspb.BinaryReader): LogoutRequest;
}

export namespace LogoutRequest {
  export type AsObject = {
    refreshToken: string,
  }
}

export class GetMeRequest extends jspb.Message {
  getAccessToken(): string;
  setAccessToken(value: string): GetMeRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMeRequest): GetMeRequest.AsObject;
  static serializeBinaryToWriter(message: GetMeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMeRequest;
  static deserializeBinaryFromReader(message: GetMeRequest, reader: jspb.BinaryReader): GetMeRequest;
}

export namespace GetMeRequest {
  export type AsObject = {
    accessToken: string,
  }
}

export class UserResponse extends jspb.Message {
  getUser(): User | undefined;
  setUser(value?: User): UserResponse;
  hasUser(): boolean;
  clearUser(): UserResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UserResponse): UserResponse.AsObject;
  static serializeBinaryToWriter(message: UserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserResponse;
  static deserializeBinaryFromReader(message: UserResponse, reader: jspb.BinaryReader): UserResponse;
}

export namespace UserResponse {
  export type AsObject = {
    user?: User.AsObject,
  }
}

