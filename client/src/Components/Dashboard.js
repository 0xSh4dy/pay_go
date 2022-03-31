import axios from "axios";
import { useEffect } from "react";
import Navbar from "./Navbar";
import Cookies from "universal-cookie";
import {useState} from "react";

function getAuthToken(){
    const cookie = new Cookies();
    return cookie.get("auth");
}
function Dashboard(){
    const [apiData,setApiData] = useState([])
    function DashBox(props){
        if(props.key==0){
            return (<div className="grid grid-cols-2 gap-x-2 mt-4" key={props._id}>
                <div className="bg-[#8B1E98] p-2 text-[#1DF3E3]">New Entry(pay/receive)</div>
                <div className="bg-[#8B1E98] p-2 text-[#1DF3E3]">{props.description}</div>

            </div>)
        }else if(props.key==1 || props.key==2){
            return (
                <div className="grid grid-cols-2 gap-x-2 mb-4"  key={props._id}>
                    <div className="bg-[#1FF1CE] text-[#F3341D] p-2">Squared Off!</div>
                    <div className="bg-[#1FF1CE]  text-[#F3341D] p-2">{props.description}</div>
                </div>
            )
        }else if(props.key==3){
            return (
                <div className="grid grid-cols-2 gap-x-2 mb-4"  key={props._id}>
                    <div className="bg-[#F13F1F] text-[#56F51F] p-2">Paid something</div>
                    <div className="bg-[#F13F1F] text-[#56F51F] p-2">{props.description}</div>
                </div>
            )
        }else if(props.key==4){
            return (
                <div className="grid grid-cols-2 gap-x-2 mb-4"  key={props._id}>
                    <div className="bg-[#B5F49A] p-2">Received something</div>
                    <div className="bg-[#B5F49A] p-2">{props.description}</div>
                </div>
            )
        }else if(props.key==5){
        return (
            <div className="grid grid-cols-2 gap-x-2 mb-4"  key={props._id}>
                <div className="bg-[#F32D1D] text-[#56F51F] p-2">New Expense</div>
                <div className="bg-[#F32D1D] text-[#56F51F] p-2">{props.description}</div>
            </div>
        )
    }
    return(
        <div className="mb-4 pl-2 text-left text-2xl pt-14" key="-1">
                <div className="bg-[#F32D1D] text-[#56F51F] p-2">Welcome to pay_go</div>
                <div className="bg-[#F32D1D] text-[#56F51F] p-2">{props.description}</div>

        </div>
    )
    }
    useEffect(()=>{
        axios({
            method:"GET",
            url:"http://127.0.0.1:7000/api/dashboard",
            params:{
                token:getAuthToken()
            }
        }).then((res)=>{
            if(res.data==null){
                setApiData([{_id:-1,description:"Easily track your expenses, payments to receive, payments to do!"}]);
           }
            else{
            setApiData(res.data);
        }
        }).catch((err)=>{
            console.log(err)
        });
    },[])
    return <div>
        <Navbar/>
        <div  className="dashboard">
            {apiData.map(DashBox)}
        </div>
    </div>
}
export default Dashboard;