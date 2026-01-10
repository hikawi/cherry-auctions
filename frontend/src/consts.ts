const api = import.meta.env.VITE_API;

export const endpoints = {
  auth: {
    login: `${api}/v1/auth/login`,
    register: `${api}/v1/auth/register`,
    forgot: `${api}/v1/auth/forgot`,
    refresh: `${api}/v1/auth/refresh`,
    logout: `${api}/v1/auth/logout`,
    verify: `${api}/v1/auth/verify`,
    verifyCheck: `${api}/v1/auth/verify/check`,
  },
  products: {
    get: `${api}/v1/products`,
    details: (id: unknown) => `${api}/v1/products/${id}`,
    favorite: `${api}/v1/products/favorite`,
    top: `${api}/v1/products/top`,
    bids: (id: unknown) => `${api}/v1/products/${id}/bids`,
    description: (id: unknown) => `${api}/v1/products/${id}/description`,
    denials: (id: unknown) => `${api}/v1/products/${id}/denials`,
  },
  categories: {
    get: `${api}/v1/categories`,
    post: `${api}/v1/categories`,
    edit: (id: unknown) => `${api}/v1/categories/${id}`,
    delete: (id: unknown) => `${api}/v1/categories/${id}`,
  },
  chat: {
    stream: `${api}/v1/chat/stream`,
    index: `${api}/v1/chat`,
    id: (id: unknown) => `${api}/v1/chat/${id}`,
  },
  users: {
    all: `${api}/v1/users`,
    request: `${api}/v1/users/request`,
    approve: `${api}/v1/users/approve`,
    avatar: `${api}/v1/users/avatar`,
    me: {
      index: `${api}/v1/users/me`,
      products: `${api}/v1/users/me/products`,
      bids: `${api}/v1/users/me/bids`,
      password: `${api}/v1/users/me/password`,
      ratings: `${api}/v1/users/me/ratings`,
      rated: `${api}/v1/users/me/rated`,
    },
  },
  questions: {
    index: `${api}/v1/questions`,
    id: (id: unknown) => `${api}/v1/questions/${id}`,
  },
};
