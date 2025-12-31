export type Profile = {
  id: number;
  name: string;
  email: string;
  roles: string[];
  oauth_type: string;
  verified: boolean;
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

export interface Seller {
  name: string;
  email?: string;
}

export interface Bid {
  id: number;
  price: number;
  automated: boolean;
  bidder: Seller;
  created_at: string;
  updated_at: string;
}

export type StepBidType = "percentage" | "fixed";

export interface Product {
  id: number;
  name: string;
  description: string;
  thumbnail_url: string;
  starting_bid: number;
  bin_price: number;
  step_bid_type: StepBidType;
  step_bid_value: number;
  bids_count: number;
  current_highest_bid?: Bid;
  allows_unrated_buyers: boolean;
  auto_extends_time: boolean;
  created_at: string;
  expired_at: string;
  seller: Seller;
}
