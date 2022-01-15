import CardGroup from 'react-bootstrap/CardGroup'
import Card from 'react-bootstrap/Card'
import Accordion from 'react-bootstrap/Accordion'
import profileImage from '../images/p.jpg'

export default function ClientDetails() {
    
    return (
    <>
    <Accordion defaultActiveKey="0">
  <Accordion.Item eventKey="0">
    <Accordion.Header>Client</Accordion.Header>
    <Accordion.Body>
      ClientInfo
    </Accordion.Body>
  </Accordion.Item>
  </ Accordion>

    <CardGroup>
    <Card>
    <Card.Img variant="top" src={profileImage} />
    <Card.Body>
      <Card.Title>Card title</Card.Title>
      <Card.Text>
        This is a wider card with supporting text below as a natural lead-in to
        additional content. This content is a little bit longer.
      </Card.Text>
    </Card.Body>
    <Card.Footer>
      <small className="text-muted">Last updated 3 mins ago</small>
    </Card.Footer>
  </Card>  
    </CardGroup>
    </>
    )


}