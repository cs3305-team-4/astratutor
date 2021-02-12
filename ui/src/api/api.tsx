import React, { useContext } from 'react';

import { API, Services } from './services';
import { AccountResponseDTO } from './definitions';
import { Route, Redirect, RouteProps } from 'react-router-dom';

import jwt_decode from 'jwt-decode';

import Config from '../config';
import { fetchRest } from './rest';

export interface AuthClaims {
  iss: string;
  sub: string;
  aud: string;
  exp: number;
  iat: number;
  nonce: string;

  email: string;
}
export interface APIContextValues {
  claims: AuthClaims | undefined;
  account: AccountResponseDTO | undefined;
  bearerToken: string | undefined;
  services: Services | undefined;

  isLoggedIn(): boolean;
  loginFromJwt(jwt: string): Promise<void>;

  // Tries to login from browser cache
  loginSilent(): Promise<void>;

  // Returns true if an autologin attempt has finished
  // This does not mean successful, use isLoggedIn to check that
  loginSilentFinished(): boolean;

  logout(): void;
}

// Only to be consumed by the root component, use APIContext.Consumer to get these values in a child component
export function useApiValues(): APIContextValues {
  const [authValues, setAuthValues] = React.useState<APIContextValues>({
    claims: undefined,
    account: undefined,
    bearerToken: undefined,
    services: undefined,
    isLoggedIn() {
      return false;
    },

    async loginSilent(): Promise<void> {
      const jwt = window.localStorage.getItem('auth-jwt');
      if (jwt !== null) {
        try {
          const claims = jwt_decode(jwt) as AuthClaims;

          const services = new Services(jwt);
          const account = await services.readAccountByID(claims.sub);

          window.localStorage.setItem('auth-jwt', jwt);

          setAuthValues({
            claims,
            account,
            bearerToken: jwt,
            services,
            loginSilent: authValues.loginSilent,
            loginSilentFinished: () => true,
            isLoggedIn: () => true,
            loginFromJwt: authValues.loginFromJwt,
            logout: authValues.logout,
          });
        } catch (e) {
          setAuthValues({
            claims: undefined,
            account: undefined,
            bearerToken: undefined,
            services: undefined,
            loginSilent: authValues.loginSilent,
            loginSilentFinished: () => true,
            isLoggedIn: () => true,
            loginFromJwt: authValues.loginFromJwt,
            logout: authValues.logout,
          });
        }
      } else {
        setAuthValues({
          claims: undefined,
          account: undefined,
          bearerToken: undefined,
          services: undefined,
          loginSilent: authValues.loginSilent,
          loginSilentFinished: () => true,
          isLoggedIn: () => false,
          loginFromJwt: authValues.loginFromJwt,
          logout: authValues.logout,
        });
      }
    },

    loginSilentFinished() {
      return false;
    },

    async loginFromJwt(jwt: string): Promise<void> {
      try {
        const claims = jwt_decode(jwt) as AuthClaims;

        const services = new Services(jwt);
        const account = await services.readAccountByID(claims.sub);

        window.localStorage.setItem('auth-jwt', jwt);

        const newAuthValues = {
          claims,
          account,
          bearerToken: jwt,
          services,
          loginSilent: authValues.loginSilent,
          loginSilentFinished: () => true,
          isLoggedIn: () => true,
          loginFromJwt: authValues.loginFromJwt,
          logout: authValues.logout,
        };

        setAuthValues(newAuthValues);
      } catch (e) {
        setAuthValues({
          claims: undefined,
          account: undefined,
          bearerToken: undefined,
          services: undefined,
          loginSilent: authValues.loginSilent,
          loginSilentFinished: () => true,
          isLoggedIn: () => true,
          loginFromJwt: authValues.loginFromJwt,
          logout: authValues.logout,
        });

        throw new Error(`error logging in with jwt: ${e}`);
      }
    },

    logout() {
      window.localStorage.removeItem('auth-jwt');
      window.location.href = '/';
    },
  });

  return authValues;
}

// default value passed here is irrelevant, will be overriden later
export const APIContext = React.createContext<APIContextValues>({} as APIContextValues);

type PrivateRouteProps = RouteProps;

export function PrivateRoute(props: PrivateRouteProps) {
  const { ...rest } = props;
  const auth = useContext(APIContext);

  if (auth.isLoggedIn()) {
    return <Route {...rest} />;
  } else {
    return <Redirect to="/login" />;
  }
}
