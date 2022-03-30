import {useEffect} from "react";
import {useState} from "react";
import axios from "axios";
import Cookies from "universal-cookie";
import Navbar from "./Navbar";
function getId(){
    let search = window.location.search;
    let params = new Proxy(new URLSearchParams(window.location.search),{
        get: function(searchParams, prop) {return searchParams.get(prop)},
    })
    return params._id;
}

function getUsername(){
    return window.location.pathname.substring(14);
}

function getAuthCookie(){
    let cookie = new Cookies();
    let authCookie = cookie.get("auth");
    return authCookie;
}
function PatchAmount(amount){
    let data = {amount:amount,token:getAuthCookie(),id:getId()};
    console.log(data);
    axios({
        method:"PATCH",
        url:"http://127.0.0.1:7000/api/transaction/",
        data:data
    }).then(res=>console.log(res)).catch(err=>console.log(err));
}
export default function SpecificTranscn(){
    const [apiData,setApiData] = useState([]);
    const [mode,setMode] = useState("Edit");
    useEffect(()=>{
        axios({
            method:"GET",
            url:"http://127.0.0.1:7000/api/transaction/"+getUsername(),
            params:{
                _id:getId(),
                token:getAuthCookie()
            }
        }).then((res)=>{
            setApiData(res.data);
            console.log(res.data);
        }).catch((error)=>{
            alert(error.response.data);
        });
    },[])
    return (
        <div className="box">
            <Navbar/>
            <div className="trnBox bg-[#7B0AF3] grid content-around">
                <div className="text-center text-2xl text-[#F3E50A]">{apiData.title}</div>
                <div className="pl-4 text-[#0AF3B7] text-xl">{apiData.description}</div>
                <div className="flex justify-around text-[#46F30A] text-xl"><div>{`${apiData.day}/${apiData.month}/${apiData.year}`}</div><div>{apiData.mode}</div><div>Amount: <span id="amount" className="p-2">{apiData.amount}</span> <button className="bg-[#E4FAF1] text-[#195DEF] rounded p-1 ml-2" onClick={()=>{
                    let amnt = document.getElementById("amount");
                    if(mode=="Edit"){
                    setMode("Save");
                    amnt.contentEditable = true;
                    amnt.focus();
                }
                else{
                    setMode("Edit");
                    amnt.contentEditable = false;
                    PatchAmount(amnt.innerHTML);
                }
                }}>{mode}</button></div></div>
            </div>
        </div>
    )
}