export interface BandMember {
  memberid: string | number;
  memberId?: string | number;
  ref_artistid?: string | number;
  refArtistId?: string | number;
  name?: string;
  role?: string;
  dob?: string;
  image?: string;
  [key: string]: unknown;
}

export interface Artist {
  artistid: string | number;
  artistId?: string | number;
  category?: string;
  name?: string;
  place?: string;
  country?: string;
  bio?: string;
  dob?: string;
  photo?: string;
  banner?: string;
  genres?: string[];
  socials?: Record<string, string>;
  events?: string[];
  members?: BandMember[];
  createdAt?: string | Date;
  creatorid?: string | number;
  creatorId?: string | number;
  issubscribed?: boolean;
  subscribed?: boolean;
  [key: string]: unknown;
}

export interface ExistingArtist {
  artistid?: string | number;
  artistId?: string | number;
  category?: string;
  name?: string;
  bio?: string;
  dob?: string;
  place?: string;
  country?: string;
  genres?: string[];
  socials?: Record<string, string>;
  [key: string]: unknown;
}

export type ArtistProfile = Artist;
