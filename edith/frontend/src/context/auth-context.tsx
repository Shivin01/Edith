import React, {useState} from 'react'
import * as auth from '../auth-provider'
import {client} from '../utils/api-client'
import {useAsync} from '../utils/hooks'
import {AuthContextData, UserData, UserDetails} from "../types";
import axios from "axios";
import {loginUser} from "../api"
import {QueryClient, QueryClientProvider, useQuery, useMutation} from 'react-query'


const defaultValue: AuthContextData = {
  user: {
    token: ''
  },
  login: (data: UserDetails) => {},
  logout: () => {},
  register: (data: any) => new Promise((resolve) => resolve())
}

const AuthContext = React.createContext<AuthContextData>(defaultValue)
AuthContext.displayName = 'AuthContext'

type AuthProviderProps = {
  children: JSX.Element
}

function AuthProvider(props: AuthProviderProps): JSX.Element {

  const {mutate: loginMutate} = useMutation(loginUser)

  const [user, setUser] = useState<UserData | null>({} as UserData);


  const login = ({username, password}: UserDetails) => {

    const l = loginMutate({username, password}, {
      onError: (error) => {
        console.log(error);
      },
      onSuccess: (data) => {
        console.log(data);
        // useHistory()
      },
    })

    console.log(l);
  }

  // save in auth provoder and windows localstorgae
  const register = React.useCallback(
    form => auth.register(form).then(user => setUser(user)),
    [],
  )

  // delete from auth provoder and windows localstorgae
  const logout = React.useCallback(() => {
    auth.logout()
    // queryCache.clear()
    setUser(null)
  }, [])

  const value = React.useMemo(
    () => ({user, login, logout, register}),
    [login, logout, register, user],
  )
  return (
    <AuthContext.Provider value={value} {...props} />
  );
}

function useAuth(): AuthContextData {
  const context = React.useContext(AuthContext)
  if (context === undefined) {
    throw new Error(`useAuth must be used within a AuthProvider`)
  }
  return context
}

// function useClient(): any {
//   const { user } = useAuth()
//   const token = user?.token
//   return React.useCallback(
//     (endpoint: string, config: any) => client(endpoint, {...config, token}),
//     [token],
//   )
// }

export {AuthProvider, useAuth}
