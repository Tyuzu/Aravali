export interface ProfileSocialLinks {
  [key: string]: string;
}

export interface ProfileUser {
  userid?: string | number;
  userId?: string | number;
  id?: string | number;
  username?: string;
  name?: string;
  email?: string;
  bio?: string;
  avatar?: string;
  banner?: string;
  phone_number?: string;
  phoneNumber?: string;
  is_following?: boolean;
  isFollowing?: boolean;
  is_verified?: boolean;
  isVerified?: boolean;
  wallet_balance?: number;
  walletBalance?: number;
  followerscount?: number;
  followersCount?: number;
  followscount?: number;
  followingCount?: number;
  social_links?: ProfileSocialLinks;
  socialLinks?: ProfileSocialLinks;
  online?: boolean;
  last_login?: string | Date | null;
  lastLogin?: string | Date | null;
  [key: string]: unknown;
}

export type UserProfile = ProfileUser;

export interface SuggestedUser {
  userid: string | number;
  userId?: string | number;
  username?: string;
  bio?: string;
  avatar?: string;
  [key: string]: unknown;
}
