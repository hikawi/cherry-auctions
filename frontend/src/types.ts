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
  email: string;
}

export interface ProductListing {
  id: number;
  name: string;
  description: string;
  thumbnail_url: string;
  bin_price: number;
  starting_bid: number;
  allows_unrated_buyers: boolean;
  auto_extends_time: boolean;
  step_bid_type: "percentage" | "fixed"; // Using a union type for better DX
  step_bid_value: number;
  seller: Seller;
  created_at: string; // ISO Date String
  expired_at: string; // ISO Date String
}
