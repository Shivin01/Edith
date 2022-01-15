import 'react-date-range/dist/styles.css';
import 'react-date-range/dist/theme/default.css';
import { DateRangePicker } from 'react-date-range';
import { useForm } from "react-hook-form";
import { useState } from 'react';

// TODO: Remove yesterday and month options.

export default function Leave() {

    //  function handleSelect(date){
    // console.log(date); // native Date object
    // }

    const [selectionRange, setSelectionRange] = useState({
      startDate: new Date(),
      endDate: new Date(),
      key: 'selection',
    })

    const {register, handleSubmit, formState: {errors}} = useForm({
        defaultValues: {
            typeofLeave: "",
            comment: ""
        }
    })

    const handleSelect = (ranges: any) => {
        console.log('Hit')
        console.log(ranges);

        setSelectionRange(ranges.selection);

    //   selection: {
    //     startDate: [date.selection.startDate],
    //     endDate: [date.selection.endDate],
    //   }



    }

    return (
        <>
      <DateRangePicker
        ranges={[selectionRange]}
        onChange={handleSelect}
      />

      <main style={{ padding: "1rem 0" }}>
          <form onSubmit={ handleSubmit((data) => {
            console.log(data);
          })}>
          <input {...register('typeofLeave', { required: "Please eneter typeofLeave" } )} placeholder="Username"></input>
          <p>{errors.typeofLeave?.message}</p>
          <input {...register('comment', { required: "Please eneter your comment " })} placeholder="comment"></input>
          <p>{errors.comment?.message}</p>
          <input type="comment"></input>
          </form>
        <h2>Your Login Page</h2>
      </main>
      </>
    )

}