
import React from 'react';
import { AuthClaims } from "../api/auth"

export interface AuthContextValues {
  claims: AuthClaims | undefined;
  isLoggedIn(): boolean;
  loginFromJwt(jwt: string): void ;
}

export default React.createContext<AuthContextValues>({
  claims: undefined,
  isLoggedIn: () => false,
  loginFromJwt: (jwt: string) => {}
});