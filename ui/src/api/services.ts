import config from '../config';
import { fetchRest } from './rest';
import {
  AccountType,
  ProfileResponseDTO,
  AccountResponseDTO,
  QualificationRequestDTO,
  WorkExperienceRequestDTO,
  LessonRequestDTO,
  LessonResponseDTO,
  ProfileRequestDTO,
  AccountRequestDTO,
  LoginRequestDTO,
  LoginResponseDTO,
} from './definitions';

export class Services {
  private headers: { [key: string]: string } = {};

  constructor(bearerToken?: string) {
    if (bearerToken !== undefined) {
      this.setBearerToken(bearerToken);
    }
  }

  private setBearerToken(bearerToken: string) {
    this.headers['Authorization'] = `Bearer ${bearerToken}`;
  }

  async createAccount(acc: AccountRequestDTO): Promise<void> {
    await fetchRest(`${config.apiUrl}/accounts`, {
      method: 'POST',
      body: JSON.stringify(acc),
    });
  }

  async login(req: LoginRequestDTO): Promise<LoginResponseDTO> {
    const res = await fetchRest(`${config.apiUrl}/auth/login`, {
      method: 'POST',
      body: JSON.stringify(req),
    });

    return (await res.json()) as LoginResponseDTO;
  }

  async readAccountByID(id: string): Promise<AccountResponseDTO> {
    const res = await fetchRest(`${config.apiUrl}/accounts/${id}`, {
      headers: this.headers,
    });

    return (await res.json()) as AccountResponseDTO;
  }

  async accountHasProfile(id: string, type: AccountType): Promise<boolean> {
    const res = await fetchRest(`${config.apiUrl}/${type}s/${id}/profile`, this.headers, [200, 404]);

    if (res.status === 200) {
      return true;
    } else {
      return false;
    }
  }

  async readProfileByAccountID(id: string, type: AccountType): Promise<ProfileResponseDTO> {
    const res = await fetchRest(`${config.apiUrl}/${type}s/${id}/profile`, {
      headers: this.headers,
    });

    return (await res.json()) as ProfileResponseDTO;
  }

  async readProfileByAccount(acc: AccountResponseDTO): Promise<ProfileResponseDTO> {
    return this.readProfileByAccountID(acc.id, acc.type);
  }

  async createProfileByAccount(acc: AccountResponseDTO, profile: ProfileRequestDTO): Promise<void> {
    await fetchRest(`${config.apiUrl}/${acc.type}s/${acc.id}/profile`, {
      method: 'POST',
      body: JSON.stringify(profile),
      headers: this.headers,
    });
  }

  async createQualificationOnProfileID(
    profileId: string,
    accountType: AccountType,
    qual: QualificationRequestDTO,
  ): Promise<void> {
    await fetchRest(`${config.apiUrl}/${accountType}s/${profileId}/profile/qualifications`, {
      method: 'POST',
      headers: this.headers,
      body: JSON.stringify(qual),
    });
  }

  async deleteQualificationOnProfileID(profileId: string, accountType: AccountType, qualId: string): Promise<void> {
    await fetchRest(`${config.apiUrl}/${accountType}s/${profileId}/profile/qualifications/${qualId}`, {
      headers: this.headers,
      method: 'DELETE',
    });
  }

  async createWorkExperienceOnProfileID(
    profileId: string,
    accountType: AccountType,
    exp: WorkExperienceRequestDTO,
  ): Promise<void> {
    await fetchRest(`${config.apiUrl}/${accountType}s/${profileId}/profile/work-experience`, {
      method: 'POST',
      headers: this.headers,
      body: JSON.stringify(exp),
    });
  }

  async deleteWorkExperienceOnProfileID(profileId: string, accountType: AccountType, expId: string): Promise<void> {
    await fetchRest(`${config.apiUrl}/${accountType}s/${profileId}/profile/work-experience/${expId}`, {
      headers: this.headers,
      method: 'DELETE',
    });
  }

  async updateDescriptionOnProfileID(profileId: string, accountType: AccountType, description: string): Promise<void> {
    await fetchRest(`${config.apiUrl}/${accountType}s/${profileId}/profile/description`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify({
        value: description,
      }),
    });
  }

  async updateAvailabilityOnProfileID(
    profileId: string,
    accountType: AccountType,
    availability: boolean[],
  ): Promise<void> {
    await fetchRest(`${config.apiUrl}/${accountType}s/${profileId}/profile/availability`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify({
        value: availability,
      }),
    });
  }

  async updateAvatarOnProfileID(profileId: string, accountType: AccountType, base64Avatar: string): Promise<void> {
    await fetchRest(`${config.apiUrl}/${accountType}s/${profileId}/profile/avatar`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify({
        value: base64Avatar,
      }),
    });
  }

  async createLesson(lesson: LessonRequestDTO): Promise<void> {
    await fetchRest(`${config.apiUrl}/lessons`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify(lesson),
    });
  }

  async readLessonsByAccountId(accountId: string): Promise<LessonResponseDTO[]> {
    const res = await fetchRest(`${config.apiUrl}/accounts/${accountId}/lessons`, {
      headers: this.headers,
      method: 'GET',
    });

    return (await res.json()) as LessonResponseDTO[];
  }

  async readLesson(lessonId: string): Promise<LessonResponseDTO> {
    const res = await fetchRest(`${config.apiUrl}/lessons/${lessonId}`, {
      headers: this.headers,
      method: 'GET',
    });

    return (await res.json()) as LessonResponseDTO;
  }

  async readLessonByAccountId(lessonId: string): Promise<LessonResponseDTO> {
    const res = await fetchRest(`${config.apiUrl}/lessons/${lessonId}`, {
      headers: this.headers,
      method: 'GET',
    });

    return (await res.json()) as LessonResponseDTO;
  }

  async updateLessonStageAccept(lesson_id: string): Promise<void> {
    await fetchRest(`${config.apiUrl}/lessons/${lesson_id}/accept`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify({
        stage_detail: 'Lesson accepted',
      }),
    });
  }

  async updateLessonStageDeny(lesson_id: string, stage_detail: string): Promise<void> {
    await fetchRest(`${config.apiUrl}/lessons/${lesson_id}/deny`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify({
        stage_detail,
      }),
    });
  }

  async updateLessonStageCancel(lesson_id: string, stage_detail: string): Promise<void> {
    await fetchRest(`${config.apiUrl}/lessons/${lesson_id}/cancel`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify({
        stage_detail,
      }),
    });
  }
}
