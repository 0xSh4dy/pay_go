import { useEffect } from "react";
import Cookies from "universal-cookie";
import Navbar from "./Navbar";
import {useState} from "react";
import axios from "axios";
function Expenses(){
    const [newEntry,setNewEntry] = useState(false);
    const [description,setDescription] = useState("");
    const [amount,setAmount] = useState(0);
    const [title,setTitle] = useState("");
    const [year,setYear] = useState(0);
    const [month,setMonth] = useState(0);
    const [btnText,setbtnText] = useState("Add new");
    const [apiData,setApiData] = useState([]);
    const [total,setTotal] = useState(0);
    useEffect(ApplyFilter,[]);

    function SubmitForm(){
        const cookie = new Cookies();
        let authCookie = cookie.get("auth");
        let data = {amount:amount,title:title,description:description,authToken:authCookie};
        axios({
            method:"post",
            url:"http://127.0.0.1:7000/api/expenses/",
            data:data,
        }).then((res)=>{
            if(res.data==="Added to database"){
                window.alert("Saved");
                ApplyFilter();
                setNewEntry(false);
            }
        });
    }

   
    function ApplyFilter(){
        const cookie = new Cookies();
        const authCookie = cookie.get("auth");
        if(authCookie===undefined){
            window.location.href = "";
        }
        let yr;
        let mth;
        if(year==""){
            yr = 0;
        }
        else{
            yr= year;
        }
        if(month==""){
            mth=0;
        }
        else{
            mth = month;
        }
        let data = {authCookie:authCookie};
        axios({
            method:"get",
            url:"http://127.0.0.1:7000/api/expenses/",
            params:{
                token:authCookie,
                year:yr,
                month:mth
            }
        }).then((res)=>{
            if(res.data==="Unauthorized access"){
                window.location.href = "/"
                setTotal(0);
            }
            else if(res.data==="Bad request"){
                window.location.href = "/"
                setTotal(0);
            }
            else{
                if(res.data===null){
                    setApiData([{_id:-1,day:0,month:month,year:year,title:"No data found",amount:0}]);
                    setTotal(0);
                }
                else{
                setApiData(res.data);
                let sum = 0;
                for(let d of res.data){
                    sum += parseInt(d.amount);
                }
                setTotal(sum);
            }
            }
        })
    }

    function ExpenseBox(props){
        return (<div className="expensesGrid justify-center mt-2" key={props._id}>
            <div className="text-[#2929F0] bg-[#57F2E2] expData mr-2">{props.day+"/"+props.month+"/"+props.year}</div>
            <div className="text-[#2929F0] bg-[#57F2E2] expTitle mr-2">{props.title}</div>
            <div className="text-[#2929F0] bg-[#57F2E2] expAmount">{props.amount}</div> 
        </div>)
    }
    function FilterBox(){
        return (
            <div className="flex justify-evenly">
                <div className="mt-5">
                    <input className="bg-[#9AD0FA] p-2" type="number" min="0" placeholder="Year" onChange={(e)=>{setYear(e.target.value)}}></input>
                </div>
                <div className="mt-5">
                    <input className="bg-[#9AD0FA] p-2" type="number" min="0" placeholder="Month" onChange={(e)=>{setMonth(e.target.value)}}></input>
                </div>
                <div className="mt-5">
                    <button className="bg-[#BDFA9A] text-blue-500 p-2 rounded-3xl" onClick={ApplyFilter}>Apply Filter</button>
                </div>
                <div className="mt-5">
                    <button onClick={()=>{
                        if(newEntry===false){
                            setbtnText("Cancel");
                        setNewEntry(true);
                    }
                    else {
                        setNewEntry(false);
                        setbtnText("Add new");
                    }
                    }} className="bg-[#BDFA9A] text-blue-500 p-2 rounded-3xl">{btnText}</button>
                </div>
            </div>
        )
    }
    function RenderExpenses(){
        return (
            <div>
                <h1>Expenses</h1>
            </div>
        )
    }

    function NewEntry(){
        return (
            <div className="newEntry mt-10">
                <div>
                    <input type="number" placeholder="Enter amount" className="p-2 expensesForm" onChange={(e)=>{setAmount(e.target.value)}}/>
                </div>
                <div>
                    <input type="text" placeholder="Enter title" className="p-2" className="expensesForm p-2" onChange={(e)=>{setTitle(e.target.value)}}/>
                </div>
                <div>
                    <input type="text" placeholder="Enter description" className="p-2 expensesForm" id="expDesc" onChange={(e)=>{setDescription(e.target.value)}}/>
                </div>
                
                <div>
                    <button onClick={SubmitForm} className="bg-[#EDFFAB] p-2 text-blue-500 expensesForm">Submit</button>
                </div>
            </div>
            
        )
    }
    if(apiData.length==0){
    return (
        <div>
            <Navbar/>
            <div className="box">
            {FilterBox()}
            {newEntry===true?NewEntry():<RenderExpenses/>}
            </div>
            
        </div>
    )
}
else {
    const adata = [...apiData];
    return (<div>
        <Navbar/>
        <div className="box">
            {FilterBox()}
            {newEntry===true?NewEntry():<RenderExpenses/>}
            <div className="expensesGrid justify-center mt-2">
            <div className="text-[#2929F0] bg-[#6EF029] expData mr-2">Date</div>
            <div className="text-[#2929F0] bg-[#6EF029] expTitle mr-2">Title</div>
            <div className="text-[#2929F0] bg-[#6EF029] expAmount">Amount</div> 
        </div>
        {adata.map(ExpenseBox)}
        <div className="expensesGrid justify-center mt-2">
            <div className="text-[#2929F0] bg-[#6EF029] expData mr-2">Hmm</div>
            <div className="text-[#2929F0] bg-[#6EF029] expTitle mr-2">Total Expenses</div>
            <div className="text-[#2929F0] bg-[#6EF029] expAmount">{total}</div> 
        </div>
        </div>
        
    </div>)
}
}
export default Expenses;