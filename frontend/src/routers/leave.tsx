import 'react-date-range/dist/styles.css';
import 'react-date-range/dist/theme/default.css';
import { DateRangePicker } from 'react-date-range';


export default function Leave() {

    //  function handleSelect(date){
    // console.log(date); // native Date object
    // }

    const selectionRange = {
      startDate: new Date(),
      endDate: new Date(),
      key: 'selection',
    }

    const handleSelect = (date: any) => {
        console.log('Hit')
        console.log(date);        
    //   selection: {
    //     startDate: [date.selection.startDate],
    //     endDate: [date.selection.endDate],
    //   }    

    }

    return (
      <DateRangePicker
        ranges={[selectionRange]}
        onChange={handleSelect}
      />
    )

}