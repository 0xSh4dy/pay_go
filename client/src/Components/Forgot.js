import axios from "axios";
import {useState} from "react";
export default function Forgot(){
    const [email,setEmail] = useState("");
    function ChangePass(){
        if(email==""){
            alert("Email cannot be empty");
        }
        else{
            const data = {email:email};
            axios({
                method:"post",
                url:"http://127.0.0.1:7000/api/change/",
                data:data
            }).then((res)=>console.log(res));
        }
    }
    return (<div className="block divResize bg-[#AEFED6] p-10 ">    
    <input type="email" placeholder="Enter email"className="inp block mt-5 p-2 bg-[#61F555] text-[#F57755] placeholder-[#F78F88] text-center" onChange={(e)=>{setEmail(e.target.value)}}/>
    <button  id="loginBtn" className="p-2 mt-5 bg-[#61F555] block rounded-md text-red-500" onClick={ChangePass}>Submit</button>
    </div>
    )
}