const api = import.meta.env.VITE_API;

export const endpoints = {
  auth: {
    refresh: `${api}/v1/auth/refresh`,
    logout: `${api}/v1/auth/logout`,
  },
  products: {
    get: `${api}/v1/products`,
    details: (id: unknown) => `${api}/v1/products/${id}`,
    favorite: `${api}/v1/products/favorite`,
    top: `${api}/v1/products/top`,
    me: `${api}/v1/products/me`,
    description: `${api}/v1/products/description`,
  },
  categories: {
    get: `${api}/v1/categories`,
    post: `${api}/v1/categories`,
    edit: (id: unknown) => `${api}/v1/categories/${id}`,
    delete: (id: unknown) => `${api}/v1/categories/${id}`,
  },
  self: `${api}/v1/users/me`,
  users: {
    all: `${api}/v1/users`,
    request: `${api}/v1/users/request`,
    approve: `${api}/v1/users/approve`,
    avatar: `${api}/v1/users/avatar`,
    profile: `${api}/v1/users/profile`,
  },
  questions: {
    index: `${api}/v1/questions`,
    id: (id: unknown) => `${api}/v1/questions/${id}`,
  },
};
