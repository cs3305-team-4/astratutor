export interface LoginDTO {
  email: string;
  password: string;
}

export interface LoginResponseDTO {
  jwt: string;
}

export interface ReadProfileDTO {
  avatar: string;
  slug: string;
  first_name: string;
  last_name: string;
  city: string;
  country: string;
  subtitle: string;
  description: string;
  qualifications?: QualificationDTO[];
  work_experience?: WorkExperienceDTO[];
  availability?: boolean[];
}

export interface CreateProfileDTO {
  avatar: string;
  first_name: string;
  last_name: string;
  city: string;
  country: string;
}

export interface QualificationDTO {
  field: string;
  degree: string;
  school: string;
  verified: boolean;
}

export interface WorkExperienceDTO {
  role: string;
  years_exp: string;
  description: string;
  verified: boolean;
}

export enum AccountType {
  Student = 'student',
  Tutor = 'tutor',
}
export interface AccountDTO {
  email: string;
  type: AccountType;
  parents_email?: string;
}

export interface SubjectDTO {
  name: string;
  image: string;
  slug: string;
}

export interface ReadSubjectsDTO {
  subjects: SubjectDTO[];
}

export interface SubjectTaughtDTO {
  subject_taught_id: string;

  subject_id: string;
  subject_name: string;

  tutor_first_name: string;
  tutor_last_name: string;
  tutor_avatar: string;

  price: number;
  description: string;
}
