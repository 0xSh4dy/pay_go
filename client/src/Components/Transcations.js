import { useEffect } from 'react'
import Cookies from 'universal-cookie'
import Navbar from './Navbar'
import { useState } from 'react'
import axios from 'axios'

function Transcations() {
  const [month, setMonth] = useState(0)
  const [year, setYear] = useState(0)
  const [addNew, setAddNew] = useState(false)
  const [newText, setnewText] = useState('Add new')
  const [mode, setMode] = useState('debt')
  const [saveMode, setSaveMode] = useState('debt')
  const [amount, setAmount] = useState(0)
  const [description, setDescription] = useState('')
  const [title, setTitle] = useState('')

  function SubmitForm() {
      const cookie = new Cookies();
      const authToken = cookie.get("auth");
      const data = {amount:amount,title:title,mode:mode,description:description,authToken:authToken}
      axios({
          url:"http://127.0.0.1:7000/api/transactions/",
          method:"POST",
          data:data
      }).then(res=>console.log(res.data))
  }
  function AddNewBox() {
    return (
      <div className="newEntry mt-10">
        <div>
          <select
            onChange={(e)=>{
                console.log(e.target.value)
                // if(e.target.value==="Debt"){
                //     setMode("credit");
                // }
                // else if (e.target.value==="Credit"){
                //     setMode("debt");
                // }
                setMode(e.target.value)
                console.log(mode)
                // console.log(mode);
            }}
            className="h-11 pl-2 rounded-2xl"
          >
            <option>Debt</option>
            <option>Credit</option>
          </select>
        </div>
        <div>
          <input
            type="number"
            placeholder="Enter amount"
            className="p-2 expensesForm"
            required
            onChange={(e)=>{
                setAmount(e.target.value)
            }}
          />
        </div>
        <div>
          <input
            type="text"
            placeholder="Enter title"
            className="p-2"
            className="expensesForm p-2"
            required
            onChange={(e)=>{
                setTitle(e.target.value)
            }}
          />
        </div>
        <div>
          <input
            type="text"
            placeholder="Enter description"
            className="p-2 expensesForm"
            id="expDesc"
            onChange={(e)=>{
                setDescription(e.target.value)
            }}
          />
        </div>

        <div>
          <button
            className="bg-[#EDFFAB] p-2 text-blue-500 expensesForm"
            onClick={SubmitForm}
          >
            Submit
          </button>
        </div>
      </div>
    )
  }
  function TranscationsFilter() {
    return (
      <div className="flex justify-around pt-2">
        <input
          className="p-2 bg-[#9AD0FA] trfilter"
          type="number"
          placeholder="Enter year"
          min="0"
          onChange={(e) => {
            setMonth(e.target.value)
          }}
        />
        <input
          className="p-2 bg-[#9AD0FA] trfilter"
          type="number"
          placeholder="Enter month"
          min="0"
          onChange={(e) => {
            setYear(e.target.value)
          }}
        />
        <select
          className="trfilter bg-[#BDFA9A]"
          onChange={(e) => {
            setMode(e.target.value)
          }}
        >
          <option>Debt</option>
          <option>Credit</option>
        </select>
        <button className="trfilter bg-[#BDFA9A] rounded-2xl p-2">
          Apply filter
        </button>
        <button
          className="trfilter bg-[#BDFA9A] rounded-2xl p-2"
          onClick={() => {
            if (addNew === true) {
              setAddNew(false)
              setnewText('Add new')
            } else {
              setAddNew(true)
              setnewText('Cancel')
            }
          }}
        >
          {newText}
        </button>
      </div>
    )
  }
  useEffect(() => {
    const cookie = new Cookies()
    const authCookie = cookie.get('auth')
    if (authCookie === undefined) {
      window.location.href = ''
    }
  }, [])

  return (
    <div>
      <Navbar />
      <div className="box">
        <TranscationsFilter />
        {addNew === true ? AddNewBox() : null}
      </div>
    </div>
  )
}
export default Transcations
