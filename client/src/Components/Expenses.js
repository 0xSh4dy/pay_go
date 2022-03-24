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
    useEffect(()=>{
        const cookie = new Cookies();
        const authCookie = cookie.get("auth");
        if(authCookie===undefined){
            window.location.href = "";
        }
        let data = {authCookie:authCookie};
        console.log(year,month);
        axios({
            method:"get",
            url:"http://127.0.0.1:7000/api/expenses/",
            params:{
                token:authCookie,
                year:year,
                month:month
            }
        }).then((res)=>{
            if(res.data==="Unauthorized access"){
                window.location.href = "/unauthorized"
            }
            else if(res.data==="Bad request"){
                window.location.href = "/badrequest"
            }
            else{
                console.log(res.data)
            }
        })
    },[])

    function SubmitForm(){
        const cookie = new Cookies();
        let authCookie = cookie.get("auth");
        let data = {amount:amount,title:title,description:description,authToken:authCookie};
        console.log(data)
        axios({
            method:"post",
            url:"http://127.0.0.1:7000/api/expenses/",
            data:data,
        }).then((res)=>{
            if(res.data==="Added to database"){
                window.alert("Saved");
                setNewEntry(false);
            }
        });
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
    return (
        <div>
            <Navbar/>
            <div className="box">
            <FilterBox/>
            {newEntry===true?NewEntry():<RenderExpenses/>}
            </div>
        </div>
    )
}
export default Expenses;