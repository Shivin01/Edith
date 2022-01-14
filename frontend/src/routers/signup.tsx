import { Form } from "react-bootstrap";
import { useForm } from "react-hook-form";

export default function SignUp() {

    const {register, handleSubmit, formState: {errors}} = useForm({
        defaultValues: {
            username: "",
            username2: "",
            password: ""
        }
    })

return(
    <main style={{ padding: "1rem 0" }}>
        <Form onSubmit={ handleSubmit((data) => {
            console.log(data);
        })}>
        <input {...register('username', { required: "Please eneter your first name" } )} placeholder="Username"></input>
          <p>{errors.username?.message}</p>
          <input {...register('username2', { required: "Please eneter your first name" } )} placeholder="Username"></input>
          <p>{errors.username2?.message}</p>
          <input {...register('password', { required: "Please eneter your Password " })} placeholder="Password"></input>
          <p>{errors.password?.message}</p>
          <input type="submit"></input>
        <h2>Your Login Page</h2>
        </Form>
    </main>

)
}
