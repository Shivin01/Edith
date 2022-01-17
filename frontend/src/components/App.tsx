import Container from 'react-bootstrap/Container';
import Login from "../pages/Login/Login";
import SignUp from "../routers/signup";
import React from "react";
import Layout from "./layout/Layout"
import '../styles/theme.scss'
import {useAuth} from "../context/auth-context";
import { Route, Routes, Navigate } from 'react-router-dom';

// @ts-ignore
const PrivateRoute = ({element: Element, ...rest}) => {
  const {user} = useAuth()
  return (
    // Show the component only when the user is logged in
    // Otherwise, redirect the user to /signin page
    // @ts-ignore
    <Route {...rest} render={props => (
      user?.token ?
        <Element {...props} />
        : <Navigate to="/signin" />
    )} />
  );
};

function App() {
  return (
    <Container>
      <Routes>
        <PrivateRoute path="/" element={<Layout/>}/>
        <Route path="login" element={<Login/>}/>
        <Route path="signup" element={<SignUp/>}/>
      </Routes>
    </Container>
  );
}

export default App;
