export interface UserDetails {
  username: string
  password: string
}

export type UserData = {
  token: string
}

export type AuthContextData = {
  user: UserData | null,
  login: (data: UserDetails) => void,
  logout: () => void
  register: (data: UserDetails) => Promise<void>
}

export interface APIClientData {
  data: object,
  token: string,
  headers: object
}
