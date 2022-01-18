import React, { useState } from 'react'
// import {queryCache} from 'react-query'
import * as auth from '../auth-provider'
import {client} from '../utils/api-client'
import {useAsync} from '../utils/hooks'
import { ContextData, UserData } from "../types";

const defaultValue: ContextData = {
  user: {
    token: ''
  },
  login: auth.login,
  logout: auth.logout,
  register: auth.register
}

const AuthContext = React.createContext<ContextData>(defaultValue)
AuthContext.displayName = 'AuthContext'

function AuthProvider(props: any): any {
  // const {
  //   data: user,
  //   status,
  //   error,
  //   isLoading,
  //   isIdle,
  //   isError,
  //   isSuccess,
  //   run,
  //   setData,
  // } = useAsync()

  const [user, setData] = useState<UserData | null>({} as UserData);
  const r = () => {

  }

  // save in auth provoder and windows localstorgae
  const login = React.useCallback(
    form => {

      // api call 
      auth.login(form).then(user => setData(user))
    },
    [],
  )

  // save in auth provoder and windows localstorgae
  const register = React.useCallback(
    form => auth.register(form).then(user => setData(user)),
    [],
  )

  // delete from auth provoder and windows localstorgae
  const logout = React.useCallback(() => {
    auth.logout()
    // queryCache.clear()
    setData(null)
  }, [])

  const value = React.useMemo(
    () => ({user, login, logout, register}),
    [login, logout, register, user],
  )

  // if (isLoading || isIdle) {
  //   // TODO: return spinner
  //   return (
  //     <span>spinner</span>
  //   )
  // }
  //
  // if (isError) {
  //   return <FullPageErrorFallback error={error} />
  // }

  // if (isSuccess) {
  return  ( 
  < AuthContext.Provider value={value} {...props} />      
);
  // }

  // throw new Error(`Unhandled status: ${status}`)
}

function useAuth(): ContextData {
  const context = React.useContext(AuthContext)
  if (context === undefined) {
    throw new Error(`useAuth must be used within a AuthProvider`)
  }
  return context
}

function useClient(): any {
  const { user } = useAuth()
  const token = user?.token
  return React.useCallback(
    (endpoint: string, config: any) => client(endpoint, {...config, token}),
    [token],
  )
}

export {AuthProvider, useAuth, useClient}
