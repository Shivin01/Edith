import axios from "axios";
import {UserDetails} from "./types"

export const loginUser = ({username, password}: UserDetails) => {
    return axios.post(`http://localhost:8000/rest-auth/login/`, {
      username: username,
      password: password
    })
}