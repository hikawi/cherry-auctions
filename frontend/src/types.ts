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
