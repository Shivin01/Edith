import Button  from 'react-bootstrap/Button';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Stack from 'react-bootstrap/Stack';
import Row from 'react-bootstrap/Row';
import { Link } from "react-router-dom";
import './App.css';

function App() {
  return (
    <Container fluid>
    {/* <Stack direction="horizontal" > */}
    <Row>
      <Col>
        <p> Welcome to Edith </p>
        </Col>            
      <Col>
        <Button as={Link as any} to='/login' variant="primary">Login</Button> {' '} <Button as={Link as any} to='/signup' variant="secondary">Sign Up</Button>
      </Col>
    </Row>
    {/* </Stack> */}
    </Container>
  );
}

export default App;
