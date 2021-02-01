
export interface AuthClaims {
    iss: string;
    sub: string;
    aud: string;
    exp: number;
    iat: number;
    nonce: string;

    email: string;
}
