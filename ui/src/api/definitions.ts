import { PaymentMethod } from '@stripe/stripe-js';

export interface LoginRequestDTO {
  email: string;
  password: string;
}

export enum AccountType {
  Student = 'student',
  Tutor = 'tutor',
}

export interface AccountRequestDTO {
  email: string;
  password: string;
  type: AccountType;
  parents_email?: string;
}

export interface LoginRequestDTO {
  email: string;
  password: string;
}

export interface LoginResponseDTO {
  jwt: string;
}

export interface AccountResponseDTO {
  id: string;
  email: string;
  type: AccountType;
  parents_email?: string;
}
export interface ProfileResponseDTO {
  account_id: string;
  avatar: string;
  slug: string;
  first_name: string;
  last_name: string;
  city: string;
  country: string;
  subtitle: string;
  description: string;
  qualifications?: QualificationResponseDTO[];
  work_experience?: WorkExperienceResponseDTO[];
  availability?: boolean[];
  color: string;
}

export interface ProfileRequestDTO {
  avatar: string;
  first_name: string;
  last_name: string;
  city: string;
  country: string;
}

export interface QualificationRequestDTO {
  field: string;
  degree: string;
  school: string;
}

export interface QualificationResponseDTO {
  id: string;
  field: string;
  degree: string;
  school: string;
  verified: boolean;
}

export interface WorkExperienceRequestDTO {
  role: string;
  years_exp: string;
  description: string;
}

export interface WorkExperienceResponseDTO {
  id: string;
  role: string;
  years_exp: string;
  description: string;
  verified: boolean;
}

export interface LessonRequestDTO {
  start_time: string; // RFC3339 timestamp
  tutor_id: string;
  student_id: string;
  lesson_detail: string;
}

export enum LessonRequestStage {
  Requested = 'requested',
  Scheduled = 'scheduled',
  Denied = 'denied',
  Cancelled = 'cancelled',
  Rescheduled = 'rescheduled',
  PaymentRequired = 'payment-required',
  Completed = 'completed',
  NoShowStudent = 'no-show-student',
  NoShowTutor = 'no-show-tutor',
  Expired = 'expired',
}
export interface LessonResponseDTO {
  id: string;
  start_time: string; // RFC3339 timestap

  requester_id: string;
  student_id: string;
  tutor_id: string;
  lesson_detail: string;
  request_stage: LessonRequestStage;
  request_stage_detail: string;
  request_stage_changer_id: string;
}

// i.e confirmed, expired, etc
export interface LessonStageChangeDTO {
  stage_detail: string;
}

export interface LessonDenyRequestDTO {
  reason: string;
}

export interface LessonCancelRequestDTO {
  reason: string;
}

export interface LessonRescheduleRequestDTO {
  new_time: string;
  reason: string;
}

export interface SubjectDTO {
  id: string;
  name: string;
  image: string;
  slug: string;
}

export interface SubjectTaughtDTO {
  id: string;
  subject_id: string;
  name: string;
  slug: string;
  description: string;
  price: string;
}

export interface TutorSubjectsDTO {
  id: string;
  first_name: string;
  last_name: string;
  avatar: string;
  slug: string;
  description: string;
  color: string;
  city: string;
  country: string;
  subjects: SubjectTaughtDTO[];
}

export interface BillingTutorOnboardURLResponseDTO {
  url: string;
}

export interface BillingTutorPanelURLResponseDTO {
  url: string;
}

export interface BillingCardSetupSessionRequestDTO {
  success_path: string;
  cancel_path: string;
}

export interface BillingCardSetupSessionResponseDTO {
  id: string;
}

export interface BillingLessonPaymentIntentSecretResponseDTO {
  id: string;
}

export interface BillingCardsResponseDTO {
  cards: PaymentMethod[];
}

export interface BillingPayeePayment {
  description: string;
  date: string; // RFC3339
  amount: number; // in cents
  remarks: string;
}

export interface BillingPayerPayment {
  description: string;
  date: string; // RFC3339
  amount: number; // in cents
  remarks: string;
  available_for_payout: boolean;
  paid_out: boolean;
}

export interface BillingPayeesPaymentsResponseDTO {
  payments: BillingPayeePayment[];
}

export interface BillingPayersPaymentsResponseDTO {
  payments: BillingPayerPayment[];
}

export interface BillingPayoutInfoResponseDTO {
  payout_balance: number; // in cents
}

export interface SubjectTaughtRequestDTO {
  subject_id: string;
  price: string;
  description: string;
}

export interface SubjectTaughtDescriptionUpdateRequestDTO {
  description: string;
}

export interface SubjectTaughtPriceUpdateRequestDTO {
  price: string;
}

export interface PaginatedResponseDTO<T> {
  total_pages: number;
  items: T;
}

export interface TurnCredentials {
  username: string;
  password: string;
}

export interface ReviewCreateDTO {
  rating: number;
  comment: string;
}

export interface ReviewDTO {
  id: string;
  created_at: string; // RFC3339 timestap
  rating: number;
  comment: string;
  student: ProfileMin;
}

export interface ProfileMin {
  account_id: string;
  avatar: string;
  slug: string;
  first_name: string;
  last_name: string;
}

export interface ReviewUpdateDTO {
  rating?: number;
  comment?: string;
}

export interface ReviewAverageDTO {
  average: number;
}

export interface SubjectRequestDTO {
  name: string;
}
