import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import {useForm} from "react-hook-form";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import Col from "react-bootstrap/Col";
import Widget from '../../components/Widget';
import s from './Login.module.scss';
import {Link} from "react-router-dom";

export default function Login() {

  const {register, handleSubmit, formState: {errors}} = useForm({
    defaultValues: {
      username: "",
      password: ""
    }
  })

  return (
    <Row className="d-flex justify-content-center">
      <Col xs={10} sm={6} lg={4}>
        <p className="text-center">Edith</p>
        <Widget className={s.widget}>
          <h4 className="mt-0">Login</h4>
          <p className="fs-sm text-muted">
            User your username and password to sign in<br/>
            Don&#39;t have an account? <Link to="/signup">Sign up now!</Link>
          </p>
          <Form onSubmit={handleSubmit((data) => {
            console.log({data});
          })}>
            <Form.Group className="mb-3" controlId="formBasicEmail">
              <Form.Label>Email address</Form.Label>
              <Form.Control
                placeholder="Enter email" {...register('username', {required: "Please enter your first name"})} />
              <Form.Text className="text-muted">
                We'll never share your email with anyone else.
              </Form.Text>
            </Form.Group>
            <Form.Group className="mb-3" controlId="formBasicPassword">
              <Form.Label>Password</Form.Label>
              <Form.Control
                placeholder="Password" {...register('password', {required: "Please eneter your Password "})} />
            </Form.Group>
            <Button variant="primary" type="submit">
              Submit
            </Button>
          </Form>
        </Widget>
      </Col>
    </Row>
  );
}
