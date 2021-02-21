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
  Accepted = 'accepted',
  Denied = 'denied',
  Cancelled = 'cancelled',
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

export interface SubjectDTO {
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
  subjects: SubjectTaughtDTO[];
}
