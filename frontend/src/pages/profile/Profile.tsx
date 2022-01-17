import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import Widget from '../../components/Widget';
import s from './Profile.module.scss';

export default function Profile() {
  return (
    <div className={s.root}>
      <h1 className="mb-lg">Edith</h1>
        <Row>
          <Col sm={6}>
            <Widget className={s.widget}>
              <h4 className="mt-0">Profile</h4>
            </Widget>
          </Col>
        </Row>
    </div>
  );
}
