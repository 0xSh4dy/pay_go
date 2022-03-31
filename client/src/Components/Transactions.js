import { useEffect } from 'react'
import Cookies from 'universal-cookie'
import Navbar from './Navbar'
import { useState } from 'react'
import axios from 'axios'

function Transactions() {
  const [month, setMonth] = useState(0)
  const [year, setYear] = useState(0)
  const [addNew, setAddNew] = useState(false)
  const [newText, setnewText] = useState('Add new')
  const [mode, setMode] = useState('debt')
  const [saveMode, setSaveMode] = useState('none')
  const [amount, setAmount] = useState(0)
  const [description, setDescription] = useState('')
  const [title, setTitle] = useState('')
  const [fetchedData, setFetchedData] = useState([])
  const [totalPay, setTotalPay] = useState(0)
  const [totalReceive, setTotalReceive] = useState(0)

  useEffect(() => {
    const cks = new Cookies()
    const token = cks.get('auth')
    axios({
      method: 'get',
      url: 'http://127.0.0.1:7000/api/transaction/',
      params: {
        month: month,
        year: year,
        mode: saveMode,
        token: token,
      },
    }).then((res) => {
      if (res.data === null) {
        setFetchedData([
          {
            _id: -1,
            day: 0,
            month: month,
            year: year,
            title: 'No data found',
            amount: 0,
          },
        ])
      } else {
        setFetchedData(res.data)
        let dbt = 0
        let crdt = 0
        for (let d of res.data) {
          if (d.mode === 'credit') {
            crdt += parseInt(d.amount)
          } else if (d.mode === 'debt') {
            dbt += parseInt(d.amount)
          }
        }
        setTotalPay(dbt)
        setTotalReceive(crdt)
      }
    }).catch((err)=>{
      window.location.href = "/";
    })
  }, [])

  function SubmitForm() {
    const cookie = new Cookies()
    const authToken = cookie.get('auth')
    let md = mode.toLowerCase()
    const data = {
      amount: amount,
      title: title,
      mode: md,
      description: description,
      authToken: authToken,
    }
    if (md === 'debt' || md === 'credit') {
      axios({
        url: 'http://127.0.0.1:7000/api/transaction/',
        method: 'POST',
        data: data,
      }).then((res) => console.log(res.data))
    } else {
      alert('Invalid mode! Mode can either be credit or debt')
    }
  }

  function ApplyFilter() {
    const cks = new Cookies()
    const token = cks.get('auth')
    console.log(saveMode)
    let ck = saveMode.toLowerCase()
    let mth = month
    let yr = year
    if (month == '') {
      mth = 0
    }
    if (year == '') {
      yr = 0
    }

    if (ck === 'debt' || ck == 'credit' || ck == 'none') {
      axios({
        method: 'get',
        url: 'http://127.0.0.1:7000/api/transaction/',
        params: {
          month: mth,
          year: yr,
          mode: ck,
          token: token,
        },
      }).then((res) => {
        if (res.data === null) {
          setFetchedData([
            {
              _id: -1,
              day: 0,
              month: month,
              year: year,
              title: 'No data found',
              amount: 0,
            },
          ])
        } else {
          setFetchedData(res.data)
          let dbt = 0
          let crdt = 0
          for (let d of res.data) {
            if (d.mode === 'credit') {
              crdt += parseInt(d.amount)
            } else if (d.mode === 'debt') {
              dbt += parseInt(d.amount)
            }
          }
          setTotalPay(dbt)
          setTotalReceive(crdt)
        }
      })
    } else {
      window.alert('Invalid mode filter! Set it to one of none/debt/credit')
    }
  }

  function TransactionBox(props) {
    return (
      <div className="grid grid-cols-4 gap-x-2 mt-2" key={props._id}>
        <div className="text-[#2929F0] bg-[#57F2E2] expData ">
          {props.day + '/' + props.month + '/' + props.year}
        </div>
        <div className="text-[#2929F0] bg-[#57F2E2] expTitle hover:text-[#F30D0A]">
          {' '}
          <a href={'/transactions/user?_id=' + props._id}>{props.title}</a>
        </div>
        <div className="text-[#2929F0] bg-[#57F2E2] expAmount">
          {props.amount}
        </div>
        <div className="text-[#2929F0] bg-[#57F2E2] expMode">{props.mode=="credit"?"To receive":"To pay"}</div>
      </div>
    )
  }

  function AddNewBox() {
    return (
      <div className="newEntry mt-10">
        <div>
          <input
            type="text"
            placeholder="Debt/credit"
            className="p-2"
            onChange={(e) => {
              setMode(e.target.value)
            }}
          />
        </div>
        <div>
          <input
            type="number"
            placeholder="Enter amount"
            className="p-2 expensesForm"
            required
            onChange={(e) => {
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
            onChange={(e) => {
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
            onChange={(e) => {
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
  function TransactionsFilter() {
    return (
      <div className="flex justify-around pt-2">
        <input
          className="p-2 bg-[#9AD0FA] trfilter"
          type="number"
          placeholder="Enter year"
          min="0"
          onChange={(e) => {
            setYear(e.target.value)
          }}
        />
        <input
          className="p-2 bg-[#9AD0FA] trfilter"
          type="number"
          placeholder="Enter month"
          min="0"
          onChange={(e) => {
            setMonth(e.target.value)
          }}
        />
        <input
          type="text"
          placeholder="Debt/credit"
          className="p-2"
          onChange={(e) => {
            setSaveMode(e.target.value)
          }}
        />

        <button
          className="trfilter bg-[#BDFA9A] rounded-2xl p-2"
          onClick={ApplyFilter}
        >
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
        {TransactionsFilter()}
        {addNew === true ? AddNewBox() : null}
        <div className="grid grid-cols-4 mt-5 gap-x-2">
          <div className="text-[#2929F0] bg-[#6EF029] expData ">Date</div>
          <div className="text-[#2929F0] bg-[#6EF029] expTitle ">Title</div>
          <div className="text-[#2929F0] bg-[#6EF029] expAmount">Amount</div>
          <div className="text-[#2929F0] bg-[#6EF029] expMode">Mode</div>
        </div>
        {fetchedData.map(TransactionBox)}
        <div className="flex justify-around mt-5">
          <div className="bg-[#9AB6F4] p-3 text-lg">
            <span>Total amount you have to pay: </span>
            {totalPay}
          </div>
          <div className="bg-[#9AB6F4] p-3 text-lg">
            <span>Total amount you have to receive: </span>
            {totalReceive}
          </div>
        </div>
      </div>
    </div>
  )
}
export default Transactions
