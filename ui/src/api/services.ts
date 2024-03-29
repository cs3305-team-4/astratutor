import { PaymentMethod } from '@stripe/stripe-js';

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
  LessonDenyRequestDTO,
  ProfileRequestDTO,
  AccountRequestDTO,
  LoginRequestDTO,
  LoginResponseDTO,
  TurnCredentials,
  SubjectDTO,
  SubjectTaughtDTO,
  TutorSubjectsDTO,
  BillingTutorOnboardURLResponseDTO,
  BillingTutorPanelURLResponseDTO,
  BillingLessonPaymentIntentSecretResponseDTO,
  BillingCardSetupSessionRequestDTO,
  BillingCardSetupSessionResponseDTO,
  BillingCardsResponseDTO,
  BillingPayeePayment,
  BillingPayerPayment,
  BillingPayeesPaymentsResponseDTO,
  BillingPayersPaymentsResponseDTO,
  BillingPayoutInfoResponseDTO,
  SubjectTaughtPriceUpdateRequestDTO,
  SubjectTaughtRequestDTO,
  SubjectTaughtDescriptionUpdateRequestDTO,
  PaginatedResponseDTO,
  LessonCancelRequestDTO,
  LessonRescheduleRequestDTO,
  ReviewDTO,
  ReviewAverageDTO,
  ReviewUpdateDTO,
  ReviewCreateDTO,
  SubjectRequestDTO,
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

  async createSubjectTaughtOnProfileID(
    profileId: string,
    accountType: AccountType,
    SubjectTaught: SubjectTaughtRequestDTO,
  ): Promise<void> {
    await fetchRest(`${config.apiUrl}/${accountType}s/${profileId}/subjects/${SubjectTaught.subject_id}`, {
      method: 'POST',
      headers: this.headers,
      body: JSON.stringify(SubjectTaught),
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

  async updateSubtitleOnProfileID(profileId: string, accountType: AccountType, subtitle: string): Promise<void> {
    await fetchRest(`${config.apiUrl}/${accountType}s/${profileId}/profile/subtitle`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify({
        value: subtitle,
      }),
    });
  }

  async updateSubjectDescriptionOnProfileID(
    profileId: string,
    subjectTaughtID: string,
    accountType: AccountType,
    UpdateDescription: SubjectTaughtDescriptionUpdateRequestDTO,
  ): Promise<void> {
    await fetchRest(`${config.apiUrl}/${accountType}s/${profileId}/subjects/${subjectTaughtID}/description`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify({
        description: UpdateDescription.description,
      }),
    });
  }

  async updateSubjectPriceOnProfileID(
    profileId: string,
    subjectTaughtID: string,
    accountType: AccountType,
    UpdatePrice: SubjectTaughtPriceUpdateRequestDTO,
  ): Promise<void> {
    await fetchRest(`${config.apiUrl}/${accountType}s/${profileId}/subjects/${subjectTaughtID}/cost`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify({
        price: UpdatePrice.price,
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

  async updateLessonStagePaymentRequired(lesson_id: string): Promise<void> {
    await fetchRest(`${config.apiUrl}/lessons/${lesson_id}/payment-required`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify({
        stage_detail: 'Lesson payment requested',
      }),
    });
  }

  async updateLessonStageDeny(lesson_id: string, denyRequest: LessonDenyRequestDTO): Promise<void> {
    await fetchRest(`${config.apiUrl}/lessons/${lesson_id}/deny`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify(denyRequest),
    });
  }

  async updateLessonStageCancel(lesson_id: string, cancelRequest: LessonCancelRequestDTO): Promise<void> {
    await fetchRest(`${config.apiUrl}/lessons/${lesson_id}/cancel`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify(cancelRequest),
    });
  }

  async updateLessonStageReschedule(lesson_id: string, rescheduleRequest: LessonRescheduleRequestDTO): Promise<void> {
    await fetchRest(`${config.apiUrl}/lessons/${lesson_id}/reschedule`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify(rescheduleRequest),
    });
  }

  async updateLessonStageCompleted(lesson_id: string): Promise<void> {
    await fetchRest(`${config.apiUrl}/lessons/${lesson_id}/completed`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify({}),
    });
  }

  async getTurnCredentials(): Promise<TurnCredentials> {
    const res = await fetchRest(`${config.apiUrl}/signalling/credentials`, {
      headers: this.headers,
      method: 'GET',
    });

    return (await res.json()) as TurnCredentials;
  }

  async readSubjects(query: string): Promise<SubjectDTO[]> {
    const res = await fetchRest(`${config.apiUrl}/subjects?query=${query}`);

    return await res.json();
  }

  async readTutors(
    page: number,
    pageSize: number,
    filters?: string[],
    query?: string,
    sort?: string,
  ): Promise<PaginatedResponseDTO<TutorSubjectsDTO[]>> {
    // const url = filters
    //   ? `${config.apiUrl}/subjects/tutors?filter=${filters.join(',')}`
    //   : `${config.apiUrl}/subjects/tutors`;
    const url = `${config.apiUrl}/subjects/tutors?page_size=${pageSize}&page=${page}&filter=${filters?.join(
      ',',
    )}&query=${query}&sort=${sort}`;
    const res = await fetchRest(url);

    return (await res.json()) as PaginatedResponseDTO<TutorSubjectsDTO[]>;
  }

  async readSubjectsTaughtByAccountId(account_id: string): Promise<SubjectTaughtDTO[]> {
    const res = await fetchRest(`${config.apiUrl}/subjects/tutors/${account_id}`);
    return (await res.json()) as SubjectTaughtDTO[];
  }

  async readTutorBillingOnboardStatus(account_id: string): Promise<boolean> {
    const res = await fetchRest(
      `${config.apiUrl}/accounts/${account_id}/billing/tutor-onboard`,
      {
        headers: this.headers,
        method: 'GET',
      },
      [200, 404],
    );

    if (res.status === 200) {
      return true;
    } else {
      return false;
    }
  }

  async readTutorBillingRequirementsMetStatus(account_id: string): Promise<boolean> {
    const res = await fetchRest(
      `${config.apiUrl}/accounts/${account_id}/billing/tutor-requirements-met`,
      {
        headers: this.headers,
        method: 'GET',
      },
      [200, 404],
    );

    if (res.status === 200) {
      return true;
    } else {
      return false;
    }
  }

  async readTutorBillingOnboardUrl(account_id: string): Promise<string> {
    const res = await fetchRest(`${config.apiUrl}/accounts/${account_id}/billing/tutor-onboard-url`, {
      headers: this.headers,
      method: 'GET',
    });

    return ((await res.json()) as BillingTutorOnboardURLResponseDTO).url;
  }

  async readTutorBillingPanelUrl(account_id: string): Promise<string> {
    const res = await fetchRest(`${config.apiUrl}/accounts/${account_id}/billing/tutor-panel-url`, {
      headers: this.headers,
      method: 'GET',
    });

    return ((await res.json()) as BillingTutorPanelURLResponseDTO).url;
  }

  // async createLessonPaymenCheckoutSession(lesson_id: string): Promise<string> {
  //   const res = await fetchRest(`${config.apiUrl}/lessons/${lesson_id}/checkout-session`, {
  //     headers: this.headers,
  //     method: 'POST',
  //   });

  //   return ((await res.json()) as LessonPaymentCheckoutSessionResponseDTO).id;
  // }

  async readLessonBillingPaymentIntentSecret(lesson_id: string): Promise<string> {
    const res = await fetchRest(`${config.apiUrl}/lessons/${lesson_id}/payment-intent-secret`, {
      headers: this.headers,
      method: 'GET',
    });

    return ((await res.json()) as BillingLessonPaymentIntentSecretResponseDTO).id;
  }

  async updateLessonStageScheduled(lesson_id: string): Promise<void> {
    await fetchRest(`${config.apiUrl}/lessons/${lesson_id}/schedule`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify({
        stage_detail: 'Lesson scheduled',
      }),
    });
  }

  async createCardSetupSession(account_id: string, paths: BillingCardSetupSessionRequestDTO): Promise<string> {
    const res = await fetchRest(`${config.apiUrl}/accounts/${account_id}/billing/card-setup-session`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify(paths),
    });

    return ((await res.json()) as BillingCardSetupSessionResponseDTO).id;
  }

  async readCardsByAccount(account_id: string): Promise<PaymentMethod[]> {
    const res = await fetchRest(`${config.apiUrl}/accounts/${account_id}/billing/cards`, {
      headers: this.headers,
      method: 'GET',
    });

    return ((await res.json()) as BillingCardsResponseDTO).cards;
  }

  async deleteCardByAccount(account_id: string, payment_method_id: string): Promise<void> {
    await fetchRest(`${config.apiUrl}/accounts/${account_id}/billing/cards/${payment_method_id}`, {
      headers: this.headers,
      method: 'DELETE',
    });
  }

  async createPayout(account_id: string): Promise<void> {
    await fetchRest(`${config.apiUrl}/accounts/${account_id}/billing/payout`, {
      headers: this.headers,
      method: 'POST',
    });
  }

  async readPayoutInfo(account_id: string): Promise<BillingPayoutInfoResponseDTO> {
    const res = await fetchRest(`${config.apiUrl}/accounts/${account_id}/billing/payout-info`, {
      headers: this.headers,
      method: 'GET',
    });

    return (await res.json()) as BillingPayoutInfoResponseDTO;
  }

  async readPayeesPayments(account_id: string): Promise<BillingPayeePayment[]> {
    const res = await fetchRest(`${config.apiUrl}/accounts/${account_id}/billing/payees-payments`, {
      headers: this.headers,
      method: 'GET',
    });

    return ((await res.json()) as BillingPayeesPaymentsResponseDTO).payments;
  }

  async readPayersPayments(account_id: string): Promise<BillingPayerPayment[]> {
    const res = await fetchRest(`${config.apiUrl}/accounts/${account_id}/billing/payers-payments`, {
      headers: this.headers,
      method: 'GET',
    });

    return ((await res.json()) as BillingPayersPaymentsResponseDTO).payments;
  }

  async tutorGetAllReviews(tid: string): Promise<ReviewDTO[]> {
    const res = await fetchRest(`${config.apiUrl}/reviews/${tid}`);
    return (await res.json()) as ReviewDTO[];
  }

  async tutorRatingAverage(tid: string): Promise<ReviewAverageDTO> {
    const res = await fetchRest(`${config.apiUrl}/reviews/${tid}/average`);
    return (await res.json()) as ReviewAverageDTO;
  }

  async tutorGetSingleReview(tid: string, rid: string): Promise<ReviewDTO> {
    const res = await fetchRest(`${config.apiUrl}/reviews/${tid}/${rid}`);
    return (await res.json()) as ReviewDTO;
  }

  async tutorGetReviewByStudent(tid: string, sid: string): Promise<ReviewDTO> {
    const res = await fetchRest(`${config.apiUrl}/reviews/${tid}/author/${sid}`);
    return (await res.json()) as ReviewDTO;
  }

  async tutorCreateReview(tid: string, review: ReviewCreateDTO): Promise<number> {
    const res = await fetchRest(`${config.apiUrl}/reviews/${tid}`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify(review),
    });
    return res.status;
  }

  async tutorReviewUpdateRating(tid: string, rid: string, update: ReviewUpdateDTO): Promise<number> {
    const res = await fetchRest(`${config.apiUrl}/reviews/${tid}/${rid}/rating`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify(update),
    });
    return res.status;
  }

  async tutorReviewUpdateComment(tid: string, rid: string, update: ReviewUpdateDTO): Promise<number> {
    const res = await fetchRest(`${config.apiUrl}/reviews/${tid}/${rid}/comment`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify(update),
    });
    return res.status;
  }

  async tutorDeleteReview(tid: string, rid: string): Promise<number> {
    const res = await fetchRest(`${config.apiUrl}/reviews/${tid}/${rid}`, {
      headers: this.headers,
      method: 'DELETE',
    });
    return res.status;
  }
  async requestSubjectAdded(subjectRequest: SubjectRequestDTO): Promise<void> {
    await fetchRest(`${config.apiUrl}/subjects`, {
      headers: this.headers,
      method: 'POST',
      body: JSON.stringify(subjectRequest),
    });
  }
}
