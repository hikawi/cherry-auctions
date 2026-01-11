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
  average_rating: number;
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

export type ProductState = "active" | "ended" | "expired" | "cancelled";

export interface Product {
  id: number;
  name: string;
  description: string;
  thumbnail_url: string;
  starting_bid: number;
  bin_price?: number;
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
  product_state: ProductState;
  finalized_at?: string;
}

export interface Rating {
  rating: number;
  feedback: string;
  reviewer: Profile;
  reviewee: Profile;
  created_at: string;
  updated_at: string;
}

export interface Transaction {
  id: number;
  seller_id: number;
  buyer_id: number;
  product_id: number;
  final_price: number;
  transaction_status: string;
  created_at: string;
}

export interface ChatSession {
  id: number;
  seller: Profile;
  buyer: Profile;
  product: Product & { transaction?: Transaction };
}

export interface ChatMessage {
  chat_session_id: number;
  content: string;
  id: number;
  image_url: string;
  sender: Profile;
}
