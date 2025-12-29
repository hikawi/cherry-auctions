const api = import.meta.env.VITE_API;

export const endpoints = {
  auth: {
    refresh: `${api}/v1/auth/refresh`,
    logout: `${api}/v1/auth/logout`,
  },
  categories: {
    get: `${api}/v1/categories`,
    post: `${api}/v1/categories`,
  },
  self: `${api}/v1/users/me`,
};
