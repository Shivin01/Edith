// import logo from './logo.svg';
import Button  from 'react-bootstrap/Button';
import { Link } from "react-router-dom";

import './App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <p> Welcome to Edith </p>
        <Button as={Link as any} to='/login' variant="primary">Login</Button> {' '} <Button as={Link as any} to='/signup' variant="secondary">Sign Up</Button>
      </header>
    </div>
  );
}

export default App;
