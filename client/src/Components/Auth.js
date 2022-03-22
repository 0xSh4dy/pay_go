import React from 'react';
import ReactDOM from 'react-dom';
import {useState} from 'react';
import { useEffect } from 'react';
import axios from 'axios';

function Home(){
    const [mode,setMode] = useState("login");
    const [modeText,setModeText] = useState("Do not have an account? Register");
    const [username,setUsername] = useState("");
    const [password,setPassword] = useState("");
    const [email,setEmail] = useState("");
    function handleLogin(){
        let data = {
            username:username,
            password:password
        }
        axios({
            method:"post",
            url:"http://127.0.0.1:7000/api/login/",
            data:data
        }).then((res)=>{console.log(res.data)});
    }
    function handleRegister(){
        let data = {
            username:username,
            password:password,
            email:email
        }
        axios({
            method:"post",
            url:"http://127.0.0.1:7000/api/register/",
            data:data
        }).then((res)=>{console.log(res.data)});
    }
    function Login(){
        return <div className="p-5 loginForm  text-blue-600">
            <p className="text-2xl">Login</p>
            <input type="text" placeholder="Enter username" id="loginName" className="inp block p-2 mt-2" onChange={(e)=>setUsername(e.target.value)}/>
            <input type="password" placeholder="Enter password" id="loginPassword" className="inp block mt-5 p-2 " onChange={(e)=>setPassword(e.target.value)}/>
            <button onClick={handleLogin} id="loginBtn" className="p-2 mt-5 bg-[#a21caf] block rounded-md text-[#EDFC0F]" >Login!</button>
            <p className='mt-5 txtHover' onClick={()=>{
                setMode("register");
                setModeText("Already have an account? Login");
            }}>{modeText}</p>
        </div>    
    }
    function Signup(){
        return <div className="p-5 signupForm  text-blue-600">
            <p className="text-2xl">Register</p>
            <input type="text" placeholder="Enter username" id="signupName" className="inp block p-2 mt-2" onChange={(e)=>{setUsername(e.target.value)}}/>
            <input type="password" placeholder="Enter password" id="signupPassword" className="inp block mt-5 p-2 " onChange={(e)=>{setPassword(e.target.value)}}/>
            <input type="email" placeholder="Enter email" id='signupEmail' className='inp block p-2 mt-5' onChange={(e)=>{setEmail(e.target.value)}}/>
            <button onClick={handleRegister} id="signupBtn" className="p-2 mt-5 bg-[#a21caf] block rounded-md text-[#EDFC0F]" >Register!</button>
            <p className='mt-5 txtHover' onClick={
                ()=>{
                    setMode("login");
                setModeText("Do not have an account? Register");
                }
            } >{modeText}</p>
        </div>  
    }
    return <div>
        {mode==="login"?Login():Signup()}
    </div>
}

export default Home;