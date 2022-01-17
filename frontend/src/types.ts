
export interface UserDetails {
  username: string
  password: string
}

export interface ContextData {
  user: {
    token: string
  },
  login: (data: UserDetails) => Promise<any>,
  logout: () => Promise<any>
  register: (data: UserDetails) => Promise<any>
}

export interface APIClientData {
  data: object,
  token: string,
  headers: object
}
