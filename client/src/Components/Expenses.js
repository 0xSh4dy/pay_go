import { useEffect } from "react";
import Cookies from "universal-cookie";
import Navbar from "./Navbar";
import {useState} from "react";
function Expenses(){
    const [newEntry,setNewEntry] = useState(false);
    const [description,setDescription] = useState("");
    const [amount,setAmount] = useState(0);
    const [title,setTitle] = useState("");

    useEffect(()=>{
        const cookie = new Cookies();
        const authCookie = cookie.get("auth");
        if(authCookie===undefined){
            window.location.href = "";
        }
    },[])

    function FilterBox(){
        return (
            <div className="flex justify-evenly">
                <div className="mt-5">
                    <input className="bg-[#9AD0FA] p-2" type="number" min="0" placeholder="Year" ></input>
                </div>
                <div className="mt-5">
                    <input className="bg-[#9AD0FA] p-2" type="number" min="0" placeholder="Month" ></input>
                </div>
                <div className="mt-5">
                    <button onClick={()=>{
                        setNewEntry(true);
                    }} className="bg-[#BDFA9A] text-blue-500 p-2 rounded-3xl">Add New</button>
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
                    <input type="number" placeholder="Enter amount" className="p-2 expensesForm"/>
                </div>
                <div>
                    <input type="text" placeholder="Enter title" className="p-2" className="expensesForm p-2"/>
                </div>
                <div>
                    <input type="text" placeholder="Enter description" className="p-2 expensesForm" id="expDesc"/>
                </div>
                
                <div>
                    <button className="bg-[#EDFFAB] p-2 text-blue-500 expensesForm">Submit</button>
                </div>
            </div>
            
        )
    }
    return (
        <div>
            <Navbar/>
            <div className="box">
            <FilterBox/>
            {newEntry===true?<NewEntry/>:<RenderExpenses/>}
            </div>
        </div>
    )
}
export default Expenses;