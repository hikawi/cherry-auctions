const api = import.meta.env.VITE_API;

export const endpoints = {
  auth: {
    refresh: `${api}/v1/auth/refresh`,
    logout: `${api}/v1/auth/logout`,
  },
  self: `${api}/v1/users/me`,
};
