import { AuthContextValues } from '../api/auth';
import { ReadProfileDTO } from '../api/definitions';
import { fetchRest } from '../api/rest';
import config from '../config';

export async function GetProfile(auth: AuthContextValues): Promise<ReadProfileDTO> {
  const res = await fetchRest(`${config.apiUrl}/${auth.account?.type}s/${auth.claims?.sub}/profile`, {
    headers: {
      Authoriziation: `Bearer ${auth.bearerToken}`,
    },
  });
  return (await res.json()) as ReadProfileDTO;
}
