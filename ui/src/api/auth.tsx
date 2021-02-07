
import React, { useContext } from 'react';

import {
  Route,
  Redirect,
  RouteProps,
}
from "react-router-dom";

import jwt_decode from 'jwt-decode'

export interface AuthClaims {
    iss: string;
    sub: string;
    aud: string;
    exp: number;
    iat: number;
    nonce: string;

    email: string;
}

export interface AuthContextValues {
  claims: AuthClaims | undefined;
  isLoggedIn(): boolean;
  loginFromLocalStorage(): void;
  loginFromJwt(jwt: string): void ;
  logout(): void;
}

// Only to be consumed by the root component, use AuthContext.Consumer to get these values in a child component
export function useAuthValues() : AuthContextValues {
  const [authValues, setAuthValues] = React.useState<AuthContextValues>({
    claims: undefined,
    isLoggedIn: () => false,
  
    loginFromLocalStorage() {
      if (window.localStorage.getItem("auth-jwt") !== null) {
        let claims = jwt_decode(window.localStorage.getItem("auth-jwt") as string) as AuthClaims;

        setAuthValues((prev: AuthContextValues) => {
          return {
            claims,
            isLoggedIn: () => true,
            loginFromLocalStorage: prev.loginFromLocalStorage,
            loginFromJwt: prev.loginFromJwt,
            logout: prev.logout
          }
        })
      }
    },
    
    loginFromJwt(jwt: string) {
      try {
        let claims = jwt_decode(jwt) as AuthClaims;

        setAuthValues((prev: AuthContextValues) => {
          window.localStorage.setItem("auth-jwt", jwt)

          return {
            claims,
            isLoggedIn: () => true,
            loginFromLocalStorage: prev.loginFromLocalStorage,
            loginFromJwt: prev.loginFromJwt,
            logout: prev.logout
          }
        })
      } catch (e) {
        throw new Error(`issue with token: ${e}`)
      }
    },
  
    logout() {
      window.localStorage.removeItem("auth-jwt")
      window.location.href = "/"
    }
  })

  return authValues
}


// default value passed here is irrelevant, will be overriden later
export let AuthContext = React.createContext<AuthContextValues>({} as AuthContextValues);

interface PrivateRouteProps extends RouteProps { }

export function PrivateRoute(props: PrivateRouteProps) {
  const { ...rest } = props;
  const auth = useContext(AuthContext)

  // TODO(ocanty) need to do login redirect
  return <Route {...rest} />
  // return (<AuthContext.Consumer>
  //   <Route children="" {...rest}/>
  // </AuthContext.Consumer>)
}