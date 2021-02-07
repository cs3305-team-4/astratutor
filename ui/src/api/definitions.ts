
export interface LoginDTO {
    email: string
    password: string
}

export interface LoginResponseDTO {
    jwt: string
}

export interface ProfileDTO {
    avatar: string
    slug: string
    first_name: string
    last_name: string
    city: string
    country: string
    subtitle: string
    description: string
}

export interface QualificationDTO {
    field: string
    degree: string
    school: string
    verified: boolean
}

export interface WorkExperienceDTO {
    role: string
    years_exp: string
    description: string
    verified: boolean
}

export interface TutorDTO extends ProfileDTO {
    qualifications: QualificationDTO[]
    work_experience: WorkExperienceDTO[]
    availability: boolean[]
}