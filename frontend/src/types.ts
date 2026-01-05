export type Subscription = {
  expired_at: string;
};

export type Profile = {
  id: number;
  name?: string;
  email?: string;
  avatar_url?: string;
  address?: string;
  verified: boolean;
  created_at: string;
  average_rating: number;
  waiting_approval: boolean;
  roles: string[];
  subscription?: Subscription;
};

export type ProductImage = {
  url: string;
  alt: string;
};

export type Category = {
  id: number;
  name: string;
  slug: string;
  parent_id?: number;
  subcategories: Category[];
  created_at: string;
  updated_at: string;
  deleted_at?: string;
};

export interface SmallUser {
  id: number;
  name?: string;
  email?: string;
  avatar_url?: string;
}

export type Question = {
  id: number;
  content: string;
  answer?: string;
  user: SmallUser;
  created_at: string;
  updated_at: string;
};

export interface Bid {
  id: number;
  price: number;
  automated: boolean;
  bidder: SmallUser;
  created_at: string;
  updated_at: string;
}

export interface DescriptionChanges {
  id: number;
  changes: string;
  created_at: string;
}

export type StepBidType = "percentage" | "fixed";

export interface Product {
  id: number;
  name: string;
  description: string;
  thumbnail_url: string;
  starting_bid: number;
  bin_price?: number;
  step_bid_type: StepBidType;
  step_bid_value: number;
  bids_count: number;
  current_highest_bid?: Bid;
  allows_unrated_buyers: boolean;
  auto_extends_time: boolean;
  created_at: string;
  expired_at: string;
  seller: SmallUser;
  questions: Question[];
  bids: Bid[];
  product_images: ProductImage[];
  denied_bidders: SmallUser[];
  description_changes: DescriptionChanges[];
  is_favorite?: boolean;
}
