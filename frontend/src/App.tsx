import Container from 'react-bootstrap/Container';
import Login from "./pages/Login/Login";
import SignUp from "./routers/signup";
import React, {Fragment} from "react";
import Layout from "./components/layout/Layout"
import '../styles/theme.scss'
import {useAuth} from "./context/auth-context";
import { Route, Routes, Navigate } from 'react-router-dom';

// @ts-ignore
const PrivateRoute = ({children}) => {
  const {user} = useAuth()
  return user?.token ? children : <Navigate to="/login" />;
};

function App() {
  return (
    <Container>
    
      <Routes>
      <Fragment>
        <Route path="/" element={
          <PrivateRoute>
            <Layout/>
          </PrivateRoute>
        } />
        <Route path="/login" element={<Login/>}/>
        <Route path="/signup" element={<SignUp/>}/>
        </Fragment>
      </Routes>      
    </Container>
  );
}

export default App;