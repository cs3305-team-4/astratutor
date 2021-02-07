import React, { useContext } from 'react';

import { Route, Redirect, RouteProps } from 'react-router-dom';

import jwt_decode from 'jwt-decode';

import Config from '../config';
import { fetchRest } from '../api/rest';
import { AccountDTO } from '../api/definitions';

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
    account: AccountDTO | undefined;
    bearerToken: string | undefined;

    isLoggedIn(): boolean;
    loginFromJwt(jwt: string): void;

    // Tries to login from browser cache
    loginSilent(): void;

    // Returns true if an autologin attempt has finished
    // This does not mean successful, use isLoggedIn to check that
    loginSilentFinished(): boolean;

    logout(): void;
}

// Only to be consumed by the root component, use AuthContext.Consumer to get these values in a child component
export function useAuthValues(): AuthContextValues {
    const [authValues, setAuthValues] = React.useState<AuthContextValues>({
        claims: undefined,
        account: undefined,
        bearerToken: undefined,
        isLoggedIn() {
            return false;
        },

        loginSilent() {
            const jwt = window.localStorage.getItem('auth-jwt');
            if (jwt !== null) {
                try {
                    const claims = jwt_decode(jwt) as AuthClaims;

                    // Test by getting the users account info
                    fetchRest(`${Config.apiUrl}/accounts/${claims.sub}`, {
                        headers: {
                            Authorization: `Bearer ${jwt}`,
                        },
                    })
                        .then((res) => res.json())
                        .then((account: AccountDTO) => {
                            window.localStorage.setItem('auth-jwt', jwt);

                            setAuthValues({
                                claims,
                                account,
                                bearerToken: jwt,
                                loginSilent: authValues.loginSilent,
                                loginSilentFinished: () => true,
                                isLoggedIn: () => true,
                                loginFromJwt: authValues.loginFromJwt,
                                logout: authValues.logout,
                            });
                        });
                } catch (e) {
                    setAuthValues({
                        claims: undefined,
                        account: undefined,
                        bearerToken: undefined,
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

        loginFromJwt(jwt: string) {
            try {
                const claims = jwt_decode(jwt) as AuthClaims;

                fetchRest(`${Config.apiUrl}/accounts/${claims.sub}`, {
                    headers: {
                        Authorization: `Bearer ${jwt}`,
                    },
                })
                    .then((res) => res.json())
                    .then((account: AccountDTO) => {
                        window.localStorage.setItem('auth-jwt', jwt);

                        const newAuthValues = {
                            claims,
                            account,
                            loginSilent: authValues.loginSilent,
                            loginSilentFinished: () => true,
                            isLoggedIn: () => true,
                            loginFromJwt: authValues.loginFromJwt,
                            logout: authValues.logout,
                        };

                        setAuthValues(newAuthValues);
                    });
            } catch (e) {
                setAuthValues({
                    claims: undefined,
                    account: undefined,
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

    React.useEffect(() => {}, [authValues]);

    return authValues;
}

// default value passed here is irrelevant, will be overriden later
export const AuthContext = React.createContext<AuthContextValues>({} as AuthContextValues);

type PrivateRouteProps = RouteProps;

export function PrivateRoute(props: PrivateRouteProps) {
    const { ...rest } = props;
    const auth = useContext(AuthContext);

    if (auth.isLoggedIn()) {
        return <Route {...rest} />;
    } else {
        return <Redirect to="/login" />;
    }
}
