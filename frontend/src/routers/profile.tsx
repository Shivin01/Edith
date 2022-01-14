import { useForm } from "react-hook-form";
import Image from 'react-bootstrap/Image'
import profileImage from '../images/p.jpg'


export default function Profile() {

    // UseEffect to fetch value for the profile information.

    return(        
        <main style={{ padding: "1rem 0" }}>
            <h1>Profile</h1>
            <Image src={profileImage} fluid rounded roundedCircle></Image>
            <p>username</p>
            <p>email address</p>
            <p>1234</p>
        </main>
    )
}

