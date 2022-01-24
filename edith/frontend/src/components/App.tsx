import Container from 'react-bootstrap/Container';
import Login from "../pages/Login/Login";
import SignUp from "../routers/signup";
import React, {Fragment} from "react";
import Layout from "./layout/Layout"
import '../styles/theme.scss'
import {useAuth} from "../context/auth-context";
import { Route, Routes, Navigate, Outlet } from 'react-router-dom';

// @ts-ignore
const PrivateRoute = () => {
  const {user} = useAuth()
  return user?.token ? <Layout/> : <Navigate to="/login" />;
};

function App() {
  return (
    <Container>
      <Routes>
      <Fragment>
        <Route path="/" element={<PrivateRoute />} />
        <Route path="/login" element={<Login/>}/>
        <Route path="/signup" element={<SignUp/>}/>
      </Fragment>
      </Routes>
    </Container>
  );
}

export default App;