import { useForm } from "react-hook-form";

export default function Login() {

    const {register, handleSubmit, formState: {errors}} = useForm({
        defaultValues: {
            username: "",
            password: ""
        }
    })

    console.log(`errors`);
    console.log(errors);

    return (
      <main style={{ padding: "1rem 0" }}>
          <form onSubmit={ handleSubmit((data) => {
            console.log(data);
          })}>
          <input {...register('username', { required: "Please eneter your first name" } )} placeholder="Username"></input>
          <p>{errors.username?.message}</p>
          <input {...register('password', { required: "Please eneter your Password " })} placeholder="Password"></input>
          <p>{errors.password?.message}</p>
          <input type="submit"></input>
          </form>
        <h2>Your Login Page</h2>
      </main>
    );
  }
