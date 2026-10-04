export interface BaitoRecord {
  baitoid?: string | number;
  baitoId?: string | number;
  entityType?: string;
  entityTypeName?: string;
  entityId?: string | number;
  entityID?: string | number;
  title?: string;
  description?: string;
  category?: string;
  subcategory?: string;
  subCategory?: string;
  location?: string;
  wage?: string | number;
  phone?: string;
  requirements?: string | string[];
  banner?: string;
  bannerURL?: string;
  images?: string[];
  workHours?: string;
  workhours?: string;
  benefits?: string;
  email?: string;
  tags?: string[];
  duration?: string;
  lastdate?: string | number | Date;
  lastDateToApply?: string | number | Date;
  createdAt?: string | number | Date;
  updatedAt?: string | number | Date;
  ownerId?: string | number;
  ownerid?: string | number;
  applicationcount?: number;
  applicationsCount?: number;
  employer?: {
    avatar?: string;
    name?: string;
    verified?: boolean;
    [key: string]: unknown;
  };
  coords?: {
    lat?: number;
    lng?: number;
    latitude?: number;
    longitude?: number;
    [key: string]: unknown;
  };
  [key: string]: unknown;
}

export interface Baito extends BaitoRecord {
  baitoid: string | number;
  title: string;
  category?: string;
  location?: string;
  description?: string;
  wage?: string | number;
  createdAt?: string | number | Date;
}

export interface BaitoApplicationRecord {
  baitoappid?: string | number;
  baitoAppId?: string | number;
  baitoid?: string | number;
  baitoId?: string | number;
  userid?: string | number;
  userId?: string | number;
  username?: string;
  pitch?: string;
  submittedAt?: string | number | Date;
  createdAt?: string | number | Date;
  status?: string;
  [key: string]: unknown;
}

export interface BaitoWorker {
  userid?: string | number;
  userId?: string | number;
  baitoWorkerId?: string | number;
  baitoWorkerID?: string | number;
  workerId?: string | number;
  id?: string | number;
  name?: string;
  age?: number | string;
  phone?: string;
  location?: string;
  preferredRoles?: string[] | string;
  preferredroles?: string[] | string;
  bio?: string;
  avatar?: string;
  profilePic?: string;
  email?: string;
  experience?: string;
  skills?: string[] | string;
  profession?: string;
  availability?: string;
  expectedWage?: string | number;
  expectedwage?: string | number;
  languages?: string;
  documents?: string[];
  createdAt?: string | number | Date;
  updatedAt?: string | number | Date;
  [key: string]: unknown;
}

export interface WorkerProfileData extends BaitoWorker {
  preferredRoles?: string[] | string;
  category?: string;
  expectedWage?: number | string;
}

export interface PageResult<T> {
  items?: T[];
  results?: T[];
  data?: T[];
  total?: number;
  page?: number;
  limit?: number;
  [key: string]: unknown;
}
